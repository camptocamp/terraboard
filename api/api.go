package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/camptocamp/terraboard/auth"
	"github.com/camptocamp/terraboard/compare"
	"github.com/camptocamp/terraboard/db"
	"github.com/camptocamp/terraboard/state"
	"github.com/camptocamp/terraboard/util"
	log "github.com/sirupsen/logrus"
)

// JSONError is a wrapper function for errors
// which prints them to the http.ResponseWriter as a JSON response
func JSONError(w http.ResponseWriter, message string, err error) {
	errObj := make(map[string]string)
	errObj["error"] = message
	errObj["details"] = fmt.Sprintf("%v", err)
	j, _ := json.Marshal(errObj)
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListStates lists States
func ListStates(w http.ResponseWriter, _ *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	states := d.ListStates()

	j, err := json.Marshal(states)
	if err != nil {
		JSONError(w, "Failed to marshal states", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListStatesWithLineages lists distinct State paths with associated Lineages
func ListStatesWithLineages(w http.ResponseWriter, _ *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	states := d.ListStatesWithLineages()

	j, err := json.Marshal(states)
	if err != nil {
		JSONError(w, "Failed to marshal states", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListTerraformVersionsWithCount lists Terraform versions with their associated
// counts, sorted by the 'orderBy' parameter (version by default)
func ListTerraformVersionsWithCount(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.URL.Query()
	versions, _ := d.ListTerraformVersionsWithCount(query)

	j, err := json.Marshal(versions)
	if err != nil {
		JSONError(w, "Failed to marshal states", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListStateStats returns State information for a given path as parameter
func ListStateStats(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.URL.Query()
	states, page, total := d.ListStateStats(query)

	// Build response object
	response := make(map[string]interface{})
	response["states"] = states
	response["page"] = page
	response["total"] = total
	j, err := json.Marshal(response)
	if err != nil {
		JSONError(w, "Failed to marshal states", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// GetState provides information on a State
func GetState(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	lineage := util.TrimBasePath(r, "api/state/")
	versionID := r.URL.Query().Get("versionid")
	// var err error
	// if versionID == "" {
	// 	versionID, err = d.DefaultVersion(st)
	// 	if err != nil {
	// 		JSONError(w, "Failed to retrieve default version", err)
	// 		return
	// 	}
	// }
	state := d.GetState(lineage, versionID)

	j, err := json.Marshal(state)
	if err != nil {
		JSONError(w, "Failed to marshal state", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// GetStateActivity returns the activity (version history) of a State
func GetStateActivity(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	lineage := util.TrimBasePath(r, "api/state/activity/")
	activity := d.GetStateActivity(lineage)

	j, err := json.Marshal(activity)
	if err != nil {
		JSONError(w, "Failed to marshal state activity", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// StateCompare compares two versions ('from' and 'to') of a State
func StateCompare(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	lineage := util.TrimBasePath(r, "api/state/compare/")
	query := r.URL.Query()
	fromVersion := query.Get("from")
	toVersion := query.Get("to")

	from := d.GetState(lineage, fromVersion)
	to := d.GetState(lineage, toVersion)
	compare, err := compare.Compare(from, to)
	if err != nil {
		JSONError(w, "Failed to compare state versions", err)
		return
	}

	j, err := json.Marshal(compare)
	if err != nil {
		JSONError(w, "Failed to marshal state compare", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// GetLocks returns information on locked States
func GetLocks(w http.ResponseWriter, _ *http.Request, sps []state.Provider) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	allLocks := make(map[string]state.LockInfo)
	for _, sp := range sps {
		locks, err := sp.GetLocks()
		if err != nil {
			JSONError(w, "Failed to get locks on a provider", err)
			return
		}
		for k, v := range locks {
			allLocks[k] = v
		}
	}

	j, err := json.Marshal(allLocks)
	if err != nil {
		JSONError(w, "Failed to marshal locks", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// SearchAttribute performs a search on Resource Attributes
// by various parameters
func SearchAttribute(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.URL.Query()
	result, page, total := d.SearchAttribute(query)

	// Build response object
	response := make(map[string]interface{})
	response["results"] = result
	response["page"] = page
	response["total"] = total

	j, err := json.Marshal(response)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListResourceTypes lists all Resource types
func ListResourceTypes(w http.ResponseWriter, _ *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, _ := d.ListResourceTypes()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListResourceTypesWithCount lists all Resource types with their associated count
func ListResourceTypesWithCount(w http.ResponseWriter, _ *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, _ := d.ListResourceTypesWithCount()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListResourceNames lists all Resource names
func ListResourceNames(w http.ResponseWriter, _ *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, _ := d.ListResourceNames()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListAttributeKeys lists all Resource Attribute Keys,
// optionally filtered by resource_type
func ListAttributeKeys(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resourceType := r.URL.Query().Get("resource_type")
	result, _ := d.ListAttributeKeys(resourceType)
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListTfVersions lists all Terraform versions
func ListTfVersions(w http.ResponseWriter, _ *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, _ := d.ListTfVersions()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// GetUser returns information about the logged user
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	name := r.Header.Get("X-Forwarded-User")
	email := r.Header.Get("X-Forwarded-Email")

	user := auth.UserInfo(name, email)

	j, err := json.Marshal(user)
	if err != nil {
		JSONError(w, "Failed to marshal user information", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// SubmitPlan inserts a new Terraform plan in the database.
// /api/plans POST endpoint callback
func SubmitPlan(w http.ResponseWriter, r *http.Request, db *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to read body: %v", err)
		JSONError(w, "Failed to read body during plan submit", err)
		return
	}

	if err = db.InsertPlan(body); err != nil {
		log.Errorf("Failed to insert plan to db: %v", err)
		JSONError(w, "Failed to insert plan to db", err)
		return
	}
}

// GetPlans provides all Plan by lineage.
// Optional "&limit=X" parameter to limit requested quantity of plans.
// Sorted by most recent to oldest.
// /api/plans GET endpoint callback
func GetPlans(w http.ResponseWriter, r *http.Request, db *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	lineage := r.URL.Query().Get("lineage")
	limit := r.URL.Query().Get("limit")
	plans := db.GetPlans(lineage, limit)

	j, err := json.Marshal(plans)
	if err != nil {
		log.Errorf("Failed to marshal plans: %v", err)
		JSONError(w, "Failed to marshal plans", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ManagePlans is used to route the request to the appropriated handler function
// on /api/plans request
func ManagePlans(w http.ResponseWriter, r *http.Request, db *db.Database) {
	if r.Method == "GET" {
		GetPlans(w, r, db)
	} else if r.Method == "POST" {
		SubmitPlan(w, r, db)
	} else {
		http.Error(w, "Invalid request method.", 405)
	}
}

// GetLineages recover all Lineage from db.
// Optional "&limit=X" parameter to limit requested quantity of them.
// Sorted by most recent to oldest.
func GetLineages(w http.ResponseWriter, r *http.Request, db *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	limit := r.URL.Query().Get("limit")
	lineages := db.GetLineages(limit)

	j, err := json.Marshal(lineages)
	if err != nil {
		log.Errorf("Failed to marshal lineages: %v", err)
		JSONError(w, "Failed to marshal lineages", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}
