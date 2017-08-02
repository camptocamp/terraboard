package db

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform/terraform"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	*gorm.DB
}

type Version struct {
	ID           uint      `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	VersionID    string    `gorm:"index" json:"version_id"`
	LastModified time.Time `json:"last_modified"`
}

type State struct {
	gorm.Model `json:"-"`
	Path       string        `gorm:"index" json:"path"`
	Version    Version       `json:"version"`
	VersionID  sql.NullInt64 `gorm:"index" json:"-"`
	TFVersion  string        `json:"terraform_version"`
	Serial     int64         `json:"serial"`
	Modules    []Module      `json:"modules"`
}

type Module struct {
	ID        uint          `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	StateID   sql.NullInt64 `gorm:"index" json:"-"`
	Path      string        `json:"path"`
	Resources []Resource    `json:"resources"`
}

type Resource struct {
	ID         uint          `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	ModuleID   sql.NullInt64 `gorm:"index" json:"-"`
	Type       string        `gorm:"index" json:"type"`
	Name       string        `gorm:"index" json:"name"`
	Attributes []Attribute   `json:"attributes"`
}

type Attribute struct {
	ID         uint          `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	ResourceID sql.NullInt64 `gorm:"index" json:"-"`
	Key        string        `gorm:"index" json:"key"`
	Value      string        `gorm:"index" json:"value"`
}

var pageSize = 20

func Init(host, user, dbname, password, logLevel string) *Database {
	var err error
	connString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, dbname, password)
	db, err := gorm.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Automigrate")
	db.AutoMigrate(&Version{}, &State{}, &Module{}, &Resource{}, &Attribute{})

	if logLevel == "debug" {
		db.LogMode(true)
	}
	return &Database{db}
}

func (db *Database) stateS3toDB(state *terraform.State, path string, versionId string) (st State) {
	var version Version
	db.First(&version, Version{VersionID: versionId})
	st = State{
		Path:      path,
		Version:   version,
		TFVersion: state.TFVersion,
		Serial:    state.Serial,
	}

	for _, m := range state.Modules {
		mod := Module{
			Path: strings.Join(m.Path, "/"),
		}
		for n, r := range m.Resources {
			res := Resource{
				Type: r.Type,
				Name: n,
			}

			for k, v := range r.Primary.Attributes {
				if !isASCII(v) {
					log.WithFields(log.Fields{
						"key":          k,
						"value_base64": base64.StdEncoding.EncodeToString([]byte(v)),
					}).Info("Attribute has non-ASCII value, skipping")
					continue
				}
				res.Attributes = append(res.Attributes, Attribute{
					Key:   k,
					Value: v,
				})
			}

			mod.Resources = append(mod.Resources, res)
		}
		st.Modules = append(st.Modules, mod)
	}
	return
}

func isASCII(s string) bool {
	for _, c := range s {
		if c > 127 {
			return false
		}
	}
	return true
}

func (db *Database) InsertState(path string, versionId string, state *terraform.State) error {
	st := db.stateS3toDB(state, path, versionId)
	db.Create(&st)
	return nil
}

func (db *Database) InsertVersion(version *s3.ObjectVersion) error {
	var v Version
	db.FirstOrCreate(&v, Version{
		VersionID:    *version.VersionId,
		LastModified: *version.LastModified,
	})
	return nil
}

func (db *Database) GetState(path, versionId string) (state State) {
	db.Joins("JOIN versions on states.version_id=versions.id").
		Preload("Version").Preload("Modules").Preload("Modules.Resources").Preload("Modules.Resources.Attributes").
		Find(&state, "states.path = ? AND versions.version_id = ?", path, versionId)
	return
}

func (db *Database) GetStateActivity(path string) (states []StateStat) {
	sql := "SELECT t.path, t.serial, t.tf_version, t.version_id, t.last_modified, count(resources.*) as resource_count" +
		fmt.Sprintf(" FROM (SELECT states.id, states.path, states.serial, states.tf_version, versions.version_id, versions.last_modified FROM states JOIN versions ON versions.id = states.version_id WHERE states.path = '%s' ORDER BY states.path, versions.last_modified ASC) t", path) +
		" JOIN modules ON modules.state_id = t.id" +
		" JOIN resources ON resources.module_id = modules.id" +
		" GROUP BY t.path, t.serial, t.tf_version, t.version_id, t.last_modified" +
		" ORDER BY last_modified ASC"

	db.Raw(sql).Find(&states)
	return
}

func (db *Database) KnownVersions() (versions []string) {
	// TODO: err
	rows, _ := db.Table("versions").Select("DISTINCT version_id").Rows()
	defer rows.Close()
	for rows.Next() {
		var version string
		rows.Scan(&version) // TODO: err
		versions = append(versions, version)
	}
	return
}

type SearchResult struct {
	Path           string `gorm:"column:path" json:"path"`
	VersionId      string `gorm:"column:version_id" json:"version_id"`
	TFVersion      string `gorm:"column:tf_version" json:"tf_version"`
	Serial         int64  `gorm:"column:serial" json:"serial"`
	ModulePath     string `gorm:"column:module_path" json:"module_path"`
	ResourceType   string `gorm:"column:type" json:"resource_type"`
	ResourceName   string `gorm:"column:name" json:"resource_name"`
	AttributeKey   string `gorm:"column:key" json:"attribute_key"`
	AttributeValue string `gorm:"column:value" json:"attribute_value"`
}

func (db *Database) SearchAttribute(query url.Values) (results []SearchResult, page int, total int) {
	log.WithFields(log.Fields{
		"query": query,
	}).Info("Searching for attribute with query")

	targetVersion := string(query.Get("versionid"))

	sqlQuery := ""
	if targetVersion == "" {
		sqlQuery += " FROM (SELECT states.path, max(states.serial) as mx FROM states GROUP BY states.path) t" +
			" JOIN states ON t.path = states.path AND t.mx = states.serial"
	} else {
		sqlQuery += " FROM states"
	}

	sqlQuery += " JOIN modules ON states.id = modules.state_id" +
		" JOIN resources ON modules.id = resources.module_id" +
		" JOIN attributes ON resources.id = attributes.resource_id"

	var where []string
	if targetVersion != "" && targetVersion != "*" {
		// filter by version unless we want all (*) or most recent ("")
		where = append(where, fmt.Sprintf("states.version_id = '%s'", targetVersion))
	}

	if v := query.Get("type"); string(v) != "" {
		where = append(where, fmt.Sprintf("resources.type LIKE '%s'", fmt.Sprintf("%%%s%%", string(v))))
	}

	if v := query.Get("name"); string(v) != "" {
		where = append(where, fmt.Sprintf("resources.name LIKE '%s'", fmt.Sprintf("%%%s%%", v)))
	}

	if v := query.Get("key"); string(v) != "" {
		where = append(where, fmt.Sprintf("attributes.key LIKE '%s'", fmt.Sprintf("%%%s%%", v)))
	}

	if v := query.Get("value"); string(v) != "" {
		where = append(where, fmt.Sprintf("attributes.value LIKE '%s'", fmt.Sprintf("%%%s%%", v)))
	}

	if len(where) > 0 {
		sqlQuery += fmt.Sprintf(" WHERE %s", strings.Join(where, " AND "))
	}

	// Count everything
	row := db.Raw("SELECT count(*)" + sqlQuery).Row()
	row.Scan(&total)

	// Now get results
	// gorm doesn't support subqueries...
	sql := "SELECT states.path, states.version_id, states.tf_version, states.serial, modules.path as module_path, resources.type, resources.name, attributes.key, attributes.value" +
		sqlQuery +
		" ORDER BY states.path, states.serial, modules.path, resources.type, resources.name, attributes.key" +
		fmt.Sprintf(" LIMIT %v", pageSize)

	if v := string(query.Get("page")); v != "" {
		page, _ = strconv.Atoi(v) // TODO: err
		o := (page - 1) * pageSize
		sql += fmt.Sprintf(" OFFSET %v", o)
	} else {
		page = 1
	}

	db.Raw(sql).Find(&results)

	return
}

func (db *Database) ListStatesVersions() (statesVersions map[string][]string) {
	rows, _ := db.Table("states").
		Joins("JOIN versions ON versions.id = states.version_id").
		Select("states.path, versions.version_id").Rows()
	defer rows.Close()
	statesVersions = make(map[string][]string)
	for rows.Next() {
		var path string
		var versionId string
		rows.Scan(&path, &versionId)
		statesVersions[versionId] = append(statesVersions[versionId], path)
	}
	return
}

func (db *Database) ListStates() (states []string) {
	rows, _ := db.Table("states").Select("DISTINCT path").Rows()
	defer rows.Close()
	for rows.Next() {
		var state string
		rows.Scan(&state)
		states = append(states, state)
	}
	return
}

func (db *Database) ListTerraformVersionsWithCount(query url.Values) (results []map[string]string, err error) {
	orderBy := string(query.Get("orderBy"))
	sql := "SELECT t.tf_version, COUNT(*) FROM (SELECT DISTINCT ON(states.path) states.id, states.path, states.serial, states.tf_version, versions.version_id, versions.last_modified FROM states JOIN versions ON versions.id = states.version_id ORDER BY states.path, versions.last_modified DESC) t GROUP BY t.tf_version ORDER BY "

	if orderBy == "version" {
		sql += "string_to_array(t.tf_version, '.')::int[] DESC"
	} else {
		sql += "count DESC"
	}

	rows, err := db.Raw(sql).Rows()
	defer rows.Close()
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var name string
		var count string
		r := make(map[string]string)
		rows.Scan(&name, &count)
		r["name"] = name
		r["count"] = count
		results = append(results, r)
	}
	return
}

type StateStat struct {
	Path          string    `json:"path"`
	TFVersion     string    `json:"terraform_version"`
	Serial        int64     `json:"serial"`
	VersionID     string    `json:"version_id"`
	LastModified  time.Time `json:"last_modified"`
	ResourceCount int       `json:"resource_count"`
}

func (db *Database) ListStateStats(query url.Values) (states []StateStat, page int, total int) {
	row := db.Table("states").Select("count(DISTINCT path)").Row()
	row.Scan(&total)

	offset := 0
	page = 1
	if v := string(query.Get("page")); v != "" {
		page, _ = strconv.Atoi(v) // TODO: err
		offset = (page - 1) * pageSize
	}

	sql := "SELECT t.path, t.serial, t.tf_version, t.version_id, t.last_modified, count(resources.*) as resource_count" +
		" FROM (SELECT DISTINCT ON(states.path) states.id, states.path, states.serial, states.tf_version, versions.version_id, versions.last_modified FROM states JOIN versions ON versions.id = states.version_id ORDER BY states.path, versions.last_modified DESC) t" +
		" JOIN modules ON modules.state_id = t.id" +
		" JOIN resources ON resources.module_id = modules.id" +
		" GROUP BY t.path, t.serial, t.tf_version, t.version_id, t.last_modified" +
		" ORDER BY last_modified DESC" +
		" LIMIT 20" +
		fmt.Sprintf(" OFFSET %v", offset)

	db.Raw(sql).Find(&states)
	return
}

func (db *Database) listField(table, field string) (results []string, err error) {
	rows, err := db.Table(table).Select(fmt.Sprintf("DISTINCT %s", field)).Rows()
	defer rows.Close()
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var t string
		rows.Scan(&t)
		results = append(results, t)
	}

	return
}

func (db *Database) listFieldWithCount(table, field string) (results []map[string]string, err error) {
	rows, err := db.Table(table).Select(fmt.Sprintf("%s, COUNT(*)", field)).Group(field).Order("count DESC").Rows()
	defer rows.Close()
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var name string
		var count string
		r := make(map[string]string)
		rows.Scan(&name, &count)
		r["name"] = name
		r["count"] = count
		results = append(results, r)
	}

	return
}

func (db *Database) ListResourceTypes() ([]string, error) {
	return db.listField("resources", "type")
}

func (db *Database) ListResourceTypesWithCount() (results []map[string]string, err error) {
	sql := "SELECT resources.type, COUNT(*) FROM (SELECT DISTINCT ON(states.path) states.id, states.path, states.serial, states.tf_version, versions.version_id, versions.last_modified FROM states JOIN versions ON versions.id = states.version_id ORDER BY states.path, versions.last_modified DESC) t JOIN modules ON modules.state_id = t.id JOIN resources ON resources.module_id = modules.id GROUP BY resources.type ORDER BY count DESC"

	rows, err := db.Raw(sql).Rows()
	defer rows.Close()
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var name string
		var count string
		r := make(map[string]string)
		rows.Scan(&name, &count)
		r["name"] = name
		r["count"] = count
		results = append(results, r)
	}
	return
}

func (db *Database) ListResourceNames() ([]string, error) {
	return db.listField("resources", "name")
}

func (db *Database) ListAttributeKeys(resourceType string) (results []string, err error) {
	query := db.Table("attributes").
		Select(fmt.Sprintf("DISTINCT %s", "key")).
		Joins("JOIN resources ON attributes.resource_id = resources.id")

	if resourceType != "" {
		query = query.Where("resources.type = ?", resourceType)
	}

	rows, err := query.Rows()
	defer rows.Close()
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var t string
		rows.Scan(&t)
		results = append(results, t)
	}

	return
}

func (db *Database) DefaultVersion(path string) (version string, err error) {
	sqlQuery := "SELECT versions.version_id FROM" +
		" (SELECT states.path, max(states.serial) as mx FROM states GROUP BY states.path) t" +
		" JOIN states ON t.path = states.path AND t.mx = states.serial" +
		" JOIN versions on states.version_id=versions.id" +
		fmt.Sprintf(" WHERE states.path = '%s'", path)

	row := db.Raw(sqlQuery).Row()
	row.Scan(&version)
	return
}
