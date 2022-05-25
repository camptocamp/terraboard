package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/camptocamp/terraboard/config"
	"github.com/camptocamp/terraboard/internal/terraform/addrs"
	"github.com/camptocamp/terraboard/internal/terraform/states"
	"github.com/camptocamp/terraboard/internal/terraform/states/statefile"
	"github.com/camptocamp/terraboard/state"
	"github.com/camptocamp/terraboard/types"
	log "github.com/sirupsen/logrus"

	ctyJson "github.com/zclconf/go-cty/cty/json"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database is a wrapping structure to *gorm.DB
type Database struct {
	*gorm.DB
	lock sync.Mutex
}

var pageSize = 20

// Init setups up the Database and a pointer to it
func Init(config config.DBConfig, debug bool) *Database {
	var err error
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		config.Host,
		config.Port,
		config.User,
		config.Name,
		config.SSLMode,
		config.Password,
	)
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		Logger: &LogrusGormLogger,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Automigrate")
	err = db.AutoMigrate(
		&types.Lineage{},
		&types.Version{},
		&types.State{},
		&types.Module{},
		&types.Resource{},
		&types.Attribute{},
		&types.OutputValue{},
		&types.Plan{},
		&types.PlanModel{},
		&types.PlanModelVariable{},
		&types.PlanOutput{},
		&types.PlanResourceChange{},
		&types.PlanState{},
		&types.PlanStateModule{},
		&types.PlanStateOutput{},
		&types.PlanStateResource{},
		&types.PlanStateResourceAttribute{},
		&types.PlanStateValue{},
		&types.Change{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}

	if debug {
		db.Config.Logger.LogMode(logger.Info)
	}

	d := &Database{DB: db}
	if err = d.MigrateLineage(); err != nil {
		log.Fatalf("Lineage migration failed: %v\n", err)
	}

	return d
}

// MigrateLineage is a migration function to update db and its data to the
// new lineage db scheme. It will update State table data, delete "lineage" column
// and add corresponding Lineage entries
func (db *Database) MigrateLineage() error {
	if db.Migrator().HasColumn(&types.State{}, "lineage") {
		var states []types.State
		if err := db.Find(&states).Error; err != nil {
			return err
		}

		for _, st := range states {
			if err := db.UpdateState(st); err != nil {
				return fmt.Errorf("Failed to update %s state during lineage migration: %v", st.Path, err)
			}
		}

		// Custom migration rules
		if err := db.Migrator().DropColumn(&types.State{}, "lineage"); err != nil {
			return fmt.Errorf("Failed to drop lineage column during migration: %v", err)
		}
	}

	return nil
}

type attributeValues map[string]interface{}

func (db *Database) stateS3toDB(sf *statefile.File, path string, versionID string) (st types.State, err error) {
	var version types.Version
	db.First(&version, types.Version{VersionID: versionID})

	// Check if the associated lineage is already present in lineages table
	// If so, it recovers its ID otherwise it inserts it at the same time as the state
	var lineage types.Lineage
	db.lock.Lock()
	err = db.FirstOrCreate(&lineage, types.Lineage{Value: sf.Lineage}).Error
	if err != nil || lineage.ID == 0 {
		log.WithField("error", err).
			Error("Unknown error in stateS3toDB during lineage finding")
		return types.State{}, err
	}
	db.lock.Unlock()

	st = types.State{
		Path:      path,
		Version:   version,
		TFVersion: sf.TerraformVersion.String(),
		Serial:    int64(sf.Serial),
		LineageID: sql.NullInt64{Int64: int64(lineage.ID), Valid: true},
	}

	for _, m := range sf.State.Modules {
		mod := types.Module{
			Path: m.Addr.String(),
		}
		for _, r := range m.Resources {
			for index, i := range r.Instances {
				res := types.Resource{
					Type:       r.Addr.Resource.Type,
					Name:       r.Addr.Resource.Name,
					Index:      getResourceIndex(index),
					Attributes: marshalAttributeValues(i.Current),
				}
				mod.Resources = append(mod.Resources, res)
			}
		}

		for n, r := range m.OutputValues {
			jsonVal, err := ctyJson.Marshal(r.Value, r.Value.Type())
			if err != nil {
				log.WithError(err).Errorf("failed to load output for %s", r.Addr.String())
			}
			out := types.OutputValue{
				Sensitive: r.Sensitive,
				Name:      n,
				Value:     string(jsonVal),
			}

			mod.OutputValues = append(mod.OutputValues, out)
		}

		st.Modules = append(st.Modules, mod)
	}
	return
}

// getResourceIndex transforms an addrs.InstanceKey instance into a string representation
func getResourceIndex(index addrs.InstanceKey) string {
	switch index.(type) {
	case addrs.IntKey, addrs.StringKey:
		return index.String()
	}
	return ""
}

func marshalAttributeValues(src *states.ResourceInstanceObjectSrc) (attrs []types.Attribute) {
	vals := make(attributeValues)
	if src == nil {
		return
	}
	if src.AttrsFlat != nil {
		for k, v := range src.AttrsFlat {
			vals[k] = v
		}
	} else if err := json.Unmarshal(src.AttrsJSON, &vals); err != nil {
		log.Error(err.Error())
	}
	log.Debug(vals)

	for k, v := range vals {
		vJSON, _ := json.Marshal(v)
		attr := types.Attribute{
			Key:   k,
			Value: string(vJSON),
		}
		log.Debug(attrs)
		attrs = append(attrs, attr)
	}
	return attrs
}

// InsertState inserts a Terraform State in the Database
func (db *Database) InsertState(path string, versionID string, sf *statefile.File) error {
	st, err := db.stateS3toDB(sf, path, versionID)
	if err == nil {
		db.Create(&st)
	}
	return nil
}

// UpdateState update a Terraform State in the Database with Lineage foreign constraint
// It will also insert Lineage entry in the db if needed.
// This method is only use during the Lineage migration since States are immutable
func (db *Database) UpdateState(st types.State) error {
	// Get lineage from old column
	var lineageValue sql.NullString
	if err := db.Raw("SELECT lineage FROM states WHERE id = ?", st.ID).Scan(&lineageValue).Error; err != nil {
		return fmt.Errorf("Error on %s lineage recovering during migration: %v", st.Path, err)
	}
	if lineageValue.String == "" || !lineageValue.Valid {
		log.Warnf("Missing lineage for '%s' state, attempt to recover lineage from other states...", st.Path)
		var lineages []string
		db.Table("states").
			Distinct("lineage").
			Order("lineage desc").
			Where("path = ?", st.Path).
			Scan(&lineages)

		for _, l := range lineages {
			if l != "" {
				lineageValue.String = l
				lineageValue.Valid = true
				log.Infof("Missing lineage for '%s' state solved!", st.Path)
				break
			}
		}

		if lineageValue.String == "" || !lineageValue.Valid {
			log.Warnf("Failed to recover '%s' lineage from others states. Orphan state", st.Path)
			return nil
		}
	}

	// Create Lineage entry if not exist (value column is unique)
	lineage := types.Lineage{
		Value: lineageValue.String,
	}
	tx := db.FirstOrCreate(&lineage, lineage)
	if tx.Error != nil || lineage.ID == 0 {
		return tx.Error
	}

	// Get Lineage ID for foreign constraint
	st.LineageID = sql.NullInt64{Int64: int64(lineage.ID), Valid: true}

	return db.Save(&st).Error
}

// InsertVersion inserts an AWS S3 Version in the Database
func (db *Database) InsertVersion(version *state.Version) error {
	var v types.Version
	db.lock.Lock()
	db.FirstOrCreate(&v, types.Version{
		VersionID:    version.ID,
		LastModified: version.LastModified,
	})
	db.lock.Unlock()
	return nil
}

// GetState retrieves a State from the database by its path and versionID
func (db *Database) GetState(lineage, versionID string) (state types.State) {
	db.Joins("JOIN lineages on states.lineage_id=lineages.id").
		Joins("JOIN versions on states.version_id=versions.id").
		Preload("Version").Preload("Modules").Preload("Modules.Resources").Preload("Modules.Resources.Attributes").
		Preload("Modules.OutputValues").
		Find(&state, "lineages.value = ? AND versions.version_id = ?", lineage, versionID)
	return
}

// GetLineageActivity returns a slice of StateStat from the Database
// for a given lineage representing the State activity over time (Versions)
func (db *Database) GetLineageActivity(lineage string) (states []types.StateStat) {
	sql := "SELECT t.path, t.serial, t.tf_version, t.version_id, t.last_modified, count(resources.*) as resource_count" +
		" FROM (SELECT states.id, states.path, states.serial, states.tf_version, versions.version_id, versions.last_modified FROM states JOIN lineages ON lineages.id = states.lineage_id JOIN versions ON versions.id = states.version_id WHERE lineages.value = ? ORDER BY states.path, versions.last_modified ASC) t" +
		" JOIN modules ON modules.state_id = t.id" +
		" JOIN resources ON resources.module_id = modules.id" +
		" GROUP BY t.path, t.serial, t.tf_version, t.version_id, t.last_modified" +
		" ORDER BY last_modified ASC"

	db.Raw(sql, lineage).Find(&states)
	return
}

// KnownVersions returns a slice of all known Versions in the Database
func (db *Database) KnownVersions() (versions []string) {
	// TODO: err
	rows, _ := db.Table("versions").Select("DISTINCT version_id").Rows()
	defer rows.Close()
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			log.Error(err.Error())
		}
		versions = append(versions, version)
	}
	return
}

// SearchAttribute returns a slice of SearchResult given a query
// The query might contain parameters 'type', 'name', 'key', 'value' and 'tf_version'
// SearchAttribute also returns paging information: the page number and the total results
func (db *Database) SearchAttribute(query url.Values) (results []types.SearchResult, page int, total int) {
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
		" JOIN attributes ON resources.id = attributes.resource_id" +
		" JOIN lineages ON lineages.id = states.lineage_id" +
		" JOIN versions ON states.version_id = versions.id"

	var where []string
	var params []interface{}
	if targetVersion != "" && targetVersion != "*" {
		// filter by version unless we want all (*) or most recent ("")
		where = append(where, "states.version_id = ?")
		params = append(params, targetVersion)
	}

	if v := string(query.Get("type")); v != "" {
		where = append(where, "resources.type LIKE ?")
		params = append(params, fmt.Sprintf("%%%s%%", v))
	}

	if v := string(query.Get("name")); v != "" {
		where = append(where, "resources.name LIKE ?")
		params = append(params, fmt.Sprintf("%%%s%%", v))
	}

	if v := string(query.Get("key")); v != "" {
		where = append(where, "attributes.key LIKE ?")
		params = append(params, fmt.Sprintf("%%%s%%", v))
	}

	if v := string(query.Get("value")); v != "" {
		where = append(where, "attributes.value LIKE ?")
		params = append(params, fmt.Sprintf("%%%s%%", v))
	}

	if v := query.Get("tf_version"); string(v) != "" {
		where = append(where, "states.tf_version LIKE ?")
		params = append(params, fmt.Sprintf("%%%s%%", v))
	}

	if v := query.Get("lineage_value"); string(v) != "" {
		where = append(where, "lineages.value LIKE ?")
		params = append(params, fmt.Sprintf("%%%s%%", v))
	}

	if len(where) > 0 {
		sqlQuery += " WHERE " + strings.Join(where, " AND ")
	}

	// Count everything
	row := db.Raw("SELECT count(*)"+sqlQuery, params...).Row()
	if err := row.Scan(&total); err != nil {
		log.Error(err.Error())
	}

	// Now get results
	// gorm doesn't support subqueries...
	sql := "SELECT states.path, versions.version_id, states.tf_version, states.serial, lineages.value as lineage_value, modules.path as module_path, resources.type, resources.name, resources.index, attributes.key, attributes.value" +
		sqlQuery +
		" ORDER BY states.path, states.serial, lineage_value, modules.path, resources.type, resources.name, resources.index, attributes.key" +
		" LIMIT ?"

	params = append(params, pageSize)

	if v := string(query.Get("page")); v != "" {
		page, _ = strconv.Atoi(v) // TODO: err
		o := (page - 1) * pageSize
		sql += " OFFSET ?"
		params = append(params, o)
	} else {
		page = 1
	}

	db.Raw(sql, params...).Find(&results)

	return
}

// ListStatesVersions returns a map of Version IDs to a slice of State paths
// from the Database
func (db *Database) ListStatesVersions() (statesVersions map[string][]string) {
	rows, _ := db.Table("states").
		Joins("JOIN versions ON versions.id = states.version_id").
		Select("states.path, versions.version_id").Rows()
	defer rows.Close()
	statesVersions = make(map[string][]string)
	for rows.Next() {
		var path string
		var versionID string
		if err := rows.Scan(&path, &versionID); err != nil {
			log.Error(err.Error())
		}
		statesVersions[versionID] = append(statesVersions[versionID], path)
	}
	return
}

// ListTerraformVersionsWithCount returns a slice of maps of Terraform versions
// mapped to the count of most recent State paths using them.
// ListTerraformVersionsWithCount also takes a query with possible parameter 'orderBy'
// to sort results. Default sorting is by descending version number.
func (db *Database) ListTerraformVersionsWithCount(query url.Values) (results []map[string]string, err error) {
	orderBy := string(query.Get("orderBy"))
	sql := "SELECT t.tf_version, COUNT(*)" +
		" FROM (SELECT DISTINCT ON(states.path) states.id, states.path, states.serial, states.tf_version, versions.version_id, versions.last_modified" +
		" FROM states JOIN versions ON versions.id = states.version_id ORDER BY states.path, versions.last_modified DESC) t" +
		" GROUP BY t.tf_version ORDER BY "

	if orderBy == "version" {
		sql += "string_to_array(t.tf_version, '.')::int[] DESC"
	} else {
		sql += "count DESC"
	}

	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var count string
		r := make(map[string]string)
		if err = rows.Scan(&name, &count); err != nil {
			return
		}
		r["name"] = name
		r["count"] = count
		results = append(results, r)
	}
	return
}

// ListStateStats returns a slice of StateStat, along with paging information
func (db *Database) ListStateStats(query url.Values) (states []types.StateStat, page int, total int) {
	row := db.Raw("SELECT count(*) FROM (SELECT DISTINCT lineage_id FROM states) AS t").Row()
	if err := row.Scan(&total); err != nil {
		log.Error(err.Error())
	}

	var paginationQuery string
	var params []interface{}
	page = 1
	if v := string(query.Get("page")); v != "" {
		page, _ = strconv.Atoi(v) // TODO: err
		offset := (page - 1) * pageSize
		params = append(params, offset)
		paginationQuery = " LIMIT 20 OFFSET ?"
	} else {
		page = -1
	}

	sql := "SELECT t.path, lineages.value as lineage_value, t.serial, t.tf_version, t.version_id, t.last_modified, count(resources.*) as resource_count" +
		" FROM (SELECT DISTINCT ON(states.lineage_id) states.id, states.lineage_id, states.path, states.serial, states.tf_version, versions.version_id, versions.last_modified FROM states JOIN versions ON versions.id = states.version_id ORDER BY states.lineage_id, versions.last_modified DESC) t" +
		" JOIN modules ON modules.state_id = t.id" +
		" JOIN resources ON resources.module_id = modules.id" +
		" JOIN lineages ON lineages.id = t.lineage_id" +
		" GROUP BY t.path, lineages.value, t.serial, t.tf_version, t.version_id, t.last_modified" +
		" ORDER BY last_modified DESC" +
		paginationQuery

	db.Raw(sql, params...).Find(&states)
	return
}

// listField is a wrapper utility method to list distinct values in Database tables.
func (db *Database) listField(table, field string) (results []string, err error) {
	rows, err := db.Table(table).Select(fmt.Sprintf("DISTINCT %s", field)).Rows()
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		var t string
		if err = rows.Scan(&t); err != nil {
			return
		}
		results = append(results, t)
	}

	return
}

// ListResourceTypes lists all Resource types from the Database
func (db *Database) ListResourceTypes() ([]string, error) {
	return db.listField("resources", "type")
}

// ListResourceTypesWithCount returns a list of Resource types with associated counts
// from the Database
func (db *Database) ListResourceTypesWithCount() (results []map[string]string, err error) {
	sql := "SELECT resources.type, COUNT(*)" +
		" FROM (SELECT DISTINCT ON(states.path) states.id, states.path, states.serial, states.tf_version, versions.version_id, versions.last_modified" +
		" FROM states" +
		" JOIN versions ON versions.id = states.version_id" +
		" ORDER BY states.path, versions.last_modified DESC) t" +
		" JOIN modules ON modules.state_id = t.id" +
		" JOIN resources ON resources.module_id = modules.id" +
		" GROUP BY resources.type" +
		" ORDER BY count DESC"

	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var count string
		r := make(map[string]string)
		if err = rows.Scan(&name, &count); err != nil {
			return
		}
		r["name"] = name
		r["count"] = count
		results = append(results, r)
	}
	return
}

// ListResourceNames lists all Resource names from the Database
func (db *Database) ListResourceNames() ([]string, error) {
	return db.listField("resources", "name")
}

// ListTfVersions lists all Terraform versions from the Database
func (db *Database) ListTfVersions() ([]string, error) {
	return db.listField("states", "tf_version")
}

// ListAttributeKeys lists all Resource Attribute keys for a given Resource type
// from the Database
func (db *Database) ListAttributeKeys(resourceType string) (results []string, err error) {
	query := db.Table("attributes").
		Select("DISTINCT key").
		Joins("JOIN resources ON attributes.resource_id = resources.id")

	if resourceType != "" {
		query = query.Where("resources.type = ?", resourceType)
	}

	rows, err := query.Rows()
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		var t string
		if err = rows.Scan(&t); err != nil {
			return
		}
		results = append(results, t)
	}

	return
}

// InsertPlan inserts a Terraform plan with associated information in the Database
func (db *Database) InsertPlan(plan []byte) error {
	var lineage types.Lineage
	if err := json.Unmarshal(plan, &lineage); err != nil {
		return err
	}

	// Recover lineage from db if it's already exists or insert it
	res := db.FirstOrCreate(&lineage, lineage)
	if res.Error != nil {
		return fmt.Errorf("Error on lineage retrival during plan insertion: %v", res.Error)
	}

	var p types.Plan
	if err := json.Unmarshal(plan, &p); err != nil {
		return err
	}
	if err := json.Unmarshal(p.PlanJSON, &p.ParsedPlan); err != nil {
		return err
	}

	p.LineageID = lineage.ID
	return db.Create(&p).Error
}

// GetPlansSummary retrieves a summary of all Plans of a lineage from the database
func (db *Database) GetPlansSummary(lineage, limitStr, pageStr string) (plans []types.Plan, page int, total int) {
	var whereClause []interface{}
	var whereClauseTotal string
	if lineage != "" {
		whereClause = append(whereClause, `"Lineage"."value" = ?`, lineage)
		whereClauseTotal = ` JOIN lineages on lineages.id=t.lineage_id WHERE lineages.value = ?`
	}

	row := db.Raw("SELECT count(*) FROM plans AS t"+whereClauseTotal, lineage).Row()
	if err := row.Scan(&total); err != nil {
		log.Error(err.Error())
	}

	var limit int
	if limitStr == "" {
		limit = -1
	} else {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Warnf("GetPlans limit ignored: %v", err)
			limit = -1
		}
	}

	var offset int
	if pageStr == "" {
		offset = -1
	} else {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			log.Warnf("GetPlans offset ignored: %v", err)
		} else {
			offset = (page - 1) * pageSize
		}
	}

	db.Select(`"plans"."id"`, `"plans"."created_at"`, `"plans"."updated_at"`, `"plans"."tf_version"`,
		`"plans"."git_remote"`, `"plans"."git_commit"`, `"plans"."ci_url"`, `"plans"."source"`, `"plans"."exit_code"`).
		Joins("Lineage").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&plans, whereClause...)

	return
}

// GetPlan retrieves a specific Plan by his ID from the database
func (db *Database) GetPlan(id string) (plans types.Plan) {
	db.Joins("Lineage").
		Preload("ParsedPlan").
		Preload("ParsedPlan.PlanStateValue").
		Preload("ParsedPlan.PlanStateValue.PlanStateOutputs").
		Preload("ParsedPlan.PlanStateValue.PlanStateModule").
		Preload("ParsedPlan.PlanStateValue.PlanStateModule.PlanStateResources").
		Preload("ParsedPlan.PlanStateValue.PlanStateModule.PlanStateResources.PlanStateResourceAttributes").
		Preload("ParsedPlan.PlanStateValue.PlanStateModule.PlanStateModules").
		Preload("ParsedPlan.Variables").
		Preload("ParsedPlan.PlanResourceChanges").
		Preload("ParsedPlan.PlanResourceChanges.Change").
		Preload("ParsedPlan.PlanOutputs").
		Preload("ParsedPlan.PlanOutputs.Change").
		Preload("ParsedPlan.PlanState").
		Preload("ParsedPlan.PlanState.PlanStateValue").
		Preload("ParsedPlan.PlanState.PlanStateValue.PlanStateOutputs").
		Preload("ParsedPlan.PlanState.PlanStateValue.PlanStateModule").
		Preload("ParsedPlan.PlanState.PlanStateValue.PlanStateModule.PlanStateResources").
		Preload("ParsedPlan.PlanState.PlanStateValue.PlanStateModule.PlanStateResources.PlanStateResourceAttributes").
		Preload("ParsedPlan.PlanState.PlanStateValue.PlanStateModule.PlanStateModules").
		Find(&plans, `"plans"."id" = ?`, id)

	return
}

// GetPlans retrieves all Plan of a lineage from the database
func (db *Database) GetPlans(lineage, limitStr, pageStr string) (plans []types.Plan, page int, total int) {
	var whereClause []interface{}
	var whereClauseTotal string
	if lineage != "" {
		whereClause = append(whereClause, `"Lineage"."value" = ?`, lineage)
		whereClauseTotal = ` JOIN lineages on lineages.id=t.lineage_id WHERE lineages.value = ?`
	}

	row := db.Raw("SELECT count(*) FROM plans AS t"+whereClauseTotal, lineage).Row()
	if err := row.Scan(&total); err != nil {
		log.Error(err.Error())
	}

	var limit int
	if limitStr == "" {
		limit = -1
	} else {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Warnf("GetPlans limit ignored: %v", err)
			limit = -1
		}
	}

	var offset int
	if pageStr == "" {
		offset = -1
	} else {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			log.Warnf("GetPlans offset ignored: %v", err)
		} else {
			offset = (page - 1) * pageSize
		}
	}

	db.Joins("Lineage").
		Preload("ParsedPlan").
		Preload("ParsedPlan.PlanStateValue").
		Preload("ParsedPlan.PlanStateValue.PlanStateOutputs").
		Preload("ParsedPlan.PlanStateValue.PlanStateModule").
		Preload("ParsedPlan.PlanStateValue.PlanStateModule.PlanStateResources").
		Preload("ParsedPlan.PlanStateValue.PlanStateModule.PlanStateResources.PlanStateResourceAttributes").
		Preload("ParsedPlan.PlanStateValue.PlanStateModule.PlanStateModules").
		Preload("ParsedPlan.Variables").
		Preload("ParsedPlan.PlanResourceChanges").
		Preload("ParsedPlan.PlanResourceChanges.Change").
		Preload("ParsedPlan.PlanOutputs").
		Preload("ParsedPlan.PlanOutputs.Change").
		Preload("ParsedPlan.PlanState").
		Preload("ParsedPlan.PlanState.PlanStateValue").
		Preload("ParsedPlan.PlanState.PlanStateValue.PlanStateOutputs").
		Preload("ParsedPlan.PlanState.PlanStateValue.PlanStateModule").
		Preload("ParsedPlan.PlanState.PlanStateValue.PlanStateModule.PlanStateResources").
		Preload("ParsedPlan.PlanState.PlanStateValue.PlanStateModule.PlanStateResources.PlanStateResourceAttributes").
		Preload("ParsedPlan.PlanState.PlanStateValue.PlanStateModule.PlanStateModules").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&plans, whereClause...)

	return
}

// GetLineages retrieves all Lineage from the database
func (db *Database) GetLineages(limitStr string) (lineages []types.Lineage) {
	var limit int
	if limitStr == "" {
		limit = -1
	} else {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Warnf("GetLineages limit ignored: %v", err)
			limit = -1
		}
	}

	db.Order("created_at desc").
		Limit(limit).
		Find(&lineages)
	return
}

// DefaultVersion returns the default VersionID for a given Lineage
// Copied and adapted from github.com/hashicorp/terraform/command/jsonstate/state.go
func (db *Database) DefaultVersion(lineage string) (version string, err error) {
	sqlQuery := "SELECT versions.version_id FROM" +
		" (SELECT states.path, max(states.serial) as mx FROM states GROUP BY states.path) t" +
		" JOIN states ON t.path = states.path AND t.mx = states.serial" +
		" JOIN versions on states.version_id=versions.id" +
		" JOIN lineages on lineages.id=states.lineage_id" +
		" WHERE lineages.value = ?" +
		" ORDER BY versions.last_modified DESC"

	row := db.Raw(sqlQuery, lineage).Row()
	err = row.Scan(&version)
	return
}

// Close get generic database interface *sql.DB from the current *gorm.DB
// and close it
func (db *Database) Close() {
	sqlDb, err := db.DB.DB()
	if err != nil {
		log.Fatalf("Unable to terminate db instance: %v\n", err)
	}
	sqlDb.Close()
}
