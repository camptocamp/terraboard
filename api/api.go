package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/camptocamp/terraboard/auth"
	"github.com/camptocamp/terraboard/compare"
	"github.com/camptocamp/terraboard/db"
	"github.com/camptocamp/terraboard/s3"
	"github.com/camptocamp/terraboard/util"
)

var states []string

// JSONError is a wrapper function for errors
// which prints them to the http.ResponseWriter as a JSON response
func JSONError(w http.ResponseWriter, message string, err error) {
	errObj := make(map[string]string)
	errObj["error"] = message
	errObj["details"] = fmt.Sprintf("%v", err)
	j, _ := json.Marshal(errObj)
	io.WriteString(w, string(j))
}

// ListStates lists States
func ListStates(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	states := d.ListStates()

	j, err := json.Marshal(states)
	if err != nil {
		JSONError(w, "Failed to marshal states", err)
		return
	}
	io.WriteString(w, string(j))
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
	io.WriteString(w, string(j))
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
	io.WriteString(w, string(j))
}

// GetState provides information on a State
func GetState(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := util.TrimBasePath(r, "api/state/")
	versionID := r.URL.Query().Get("versionid")
	var err error
	if versionID == "" {
		versionID, err = d.DefaultVersion(st)
		if err != nil {
			JSONError(w, "Failed to retrieve default version", err)
			return
		}
	}
	state := d.GetState(st, versionID)

	jState, err := json.Marshal(state)
	if err != nil {
		JSONError(w, "Failed to marshal state", err)
		return
	}
	io.WriteString(w, string(jState))
}

// GetStateActivity returns the activity (version history) of a State
func GetStateActivity(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := util.TrimBasePath(r, "api/state/activity/")
	activity := d.GetStateActivity(st)

	jActivity, err := json.Marshal(activity)
	if err != nil {
		JSONError(w, "Failed to marshal state activity", err)
		return
	}
	io.WriteString(w, string(jActivity))
}

// StateCompare compares two versions ('from' and 'to') of a State
func StateCompare(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := util.TrimBasePath(r, "api/state/compare/")
	query := r.URL.Query()
	fromVersion := query.Get("from")
	toVersion := query.Get("to")

	from := d.GetState(st, fromVersion)
	to := d.GetState(st, toVersion)
	compare, err := compare.Compare(from, to)
	if err != nil {
		JSONError(w, "Failed to compare state versions", err)
		return
	}

	jCompare, err := json.Marshal(compare)
	if err != nil {
		JSONError(w, "Failed to marshal state compare", err)
		return
	}
	io.WriteString(w, string(jCompare))
}

// GetLocks returns information on locked States
func GetLocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	locks, err := s3.GetLocks()
	if err != nil {
		JSONError(w, "Failed to get locks", err)
		return
	}

	j, err := json.Marshal(locks)
	if err != nil {
		JSONError(w, "Failed to marshal locks", err)
		return
	}
	io.WriteString(w, string(j))
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
	io.WriteString(w, string(j))
}

// ListResourceTypes lists all Resource types
func ListResourceTypes(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, _ := d.ListResourceTypes()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	io.WriteString(w, string(j))
}

// ListResourceTypesWithCount lists all Resource types with their associated count
func ListResourceTypesWithCount(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, _ := d.ListResourceTypesWithCount()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	io.WriteString(w, string(j))
}

// ListResourceNames lists all Resource names
func ListResourceNames(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, _ := d.ListResourceNames()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	io.WriteString(w, string(j))
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
	io.WriteString(w, string(j))
}

// ListTfVersions lists all Terraform versions
func ListTfVersions(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, _ := d.ListTfVersions()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	io.WriteString(w, string(j))
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
	io.WriteString(w, string(j))
}
