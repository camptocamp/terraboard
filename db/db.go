package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/terraform/terraform"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var db *gorm.DB

type State struct {
	gorm.Model `json:"-"`
	Path       string   `gorm:"index" json:"path"`
	VersionId  string   `gorm:"index" json:"version_id"`
	TFVersion  string   `json:"terraform_version"`
	Serial     int64    `json:"serial"`
	Modules    []Module `json:"modules"`
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
	Type       string        `json:"type"`
	Name       string        `json:"name"`
	Attributes []Attribute   `json:"attributes"`
}

type Attribute struct {
	ID         uint          `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	ResourceID sql.NullInt64 `gorm:"index" json:"-"`
	Key        string        `json:"key"`
	Value      string        `json:"value"`
}

func Init() {
	var err error
	db, err = gorm.Open("sqlite3", "./db/terraboard.db")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	log.Infof("Automigrate")

	db.AutoMigrate(&State{}, &Module{}, &Resource{}, &Attribute{})

	db.LogMode(true)

	log.Infof("New db is %v", db)
}

func stateS3toDB(state *terraform.State, path string, versionId string) (st State) {
	st = State{
		Path:      path,
		VersionId: versionId,
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

func InsertState(path string, versionId string, state *terraform.State) error {
	var testState State
	db.Find(&testState, "path = ? AND version_id = ?", path, versionId)
	if testState.Path == path {
		log.Infof("State %s/%s is already in the DB", path, versionId)
		return nil
	}

	st := stateS3toDB(state, path, versionId)
	db.Create(&st)
	return nil
}

func GetState(path, versionId string) (state State) {
	db.Preload("Modules").Preload("Modules.Resources").Preload("Modules.Resources.Attributes").Find(&state, "path = ? AND version_id = ?", path, versionId)
	return
}

func KnownVersions() (versions []string) {
	// TODO: err
	rows, _ := db.Table("states").Select("DISTINCT version_id").Rows()
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

func SearchResource(query url.Values) (results []SearchResult) {
	log.Infof("Searching for resource with query=%v", query)

	selectQuery := make(map[string]interface{})

	statesSelect := "states.serial"
	targetVersion := string(query.Get("versionid"))
	switch targetVersion {
	case "":
		// Select most recent version of each path
		statesSelect = "max(states.serial)"
	case "*":
		targetVersion = ""
	default:
	}

	if v := query.Get("type"); string(v) != "" {
		selectQuery["type"] = string(v)
	}

	if v := query.Get("name"); string(v) != "" {
		selectQuery["name"] = string(v)
	}

	baseSelect := db.Table("resources").
		Select(fmt.Sprintf("states.path, states.version_id, states.tf_version, %s as serial, modules.path as module_path, resources.type, resources.name", statesSelect)).
		Group("states.path, modules.path, resources.type, resources.name").
		Joins("LEFT JOIN modules ON resources.module_id = modules.id LEFT JOIN states ON modules.state_id = states.id")

	if targetVersion != "" {
		// filter by version unless we want all (*) or most recent ("")
		baseSelect = baseSelect.Where("states.version_id = ?", targetVersion)
	}

	baseSelect.Where(selectQuery).
		Order("states.path, states.serial, modules.path, resources.type, resources.name").
		Find(&results)

	return
}

func SearchAttribute(query url.Values) (results []SearchResult) {
	log.Infof("Searching for attribute with query=%v", query)

	selectQuery := make(map[string]interface{})

	statesSelect := "states.serial"
	targetVersion := string(query.Get("versionid"))
	switch targetVersion {
	case "":
		// Select most recent version of each path
		statesSelect = "max(states.serial)"
	case "*":
		targetVersion = ""
	default:
	}

	baseSelect := db.Table("attributes").
		Select(fmt.Sprintf("states.path, states.version_id, states.tf_version, %s as serial, modules.path as module_path, resources.type, resources.name, attributes.key, attributes.value", statesSelect)).
		Group("states.path, modules.path, resources.type, resources.name, attributes.key").
		Joins("LEFT JOIN resources ON attributes.resource_id = resources.id LEFT JOIN modules ON resources.module_id = modules.id LEFT JOIN states ON modules.state_id = states.id")

	if targetVersion != "" {
		// filter by version unless we want all (*) or most recent ("")
		baseSelect = baseSelect.Where("states.version_id = ?", targetVersion)
	}

	if v := query.Get("type"); string(v) != "" {
		baseSelect = baseSelect.Where("resources.type LIKE ?", fmt.Sprintf("%%%s%%", v))
	}

	if v := query.Get("name"); string(v) != "" {
		baseSelect = baseSelect.Where("resources.name LIKE ?", fmt.Sprintf("%%%s%%", v))
	}

	if v := query.Get("key"); string(v) != "" {
		baseSelect = baseSelect.Where("attributes.key LIKE ?", fmt.Sprintf("%%%s%%", v))
	}

	if v := query.Get("value"); string(v) != "" {
		baseSelect = baseSelect.Where("attributes.value LIKE ?", fmt.Sprintf("%%%s%%", v))
	}

	baseSelect.Where(selectQuery).
		Order("states.path, states.serial, modules.path, resources.type, resources.name, attributes.key").
		Find(&results)

	return
}

func listField(table, field string) (results []string, err error) {
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

func ListResourceTypes() ([]string, error) {
	return listField("resources", "type")
}

func ListResourceNames() ([]string, error) {
	return listField("resources", "name")
}

func ListAttributeKeys(resourceType string) (results []string, err error) {
	query := db.Table("attributes").
		Select(fmt.Sprintf("DISTINCT %s", "key")).
		Joins("LEFT JOIN resources ON attributes.resource_id = resources.id")

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
