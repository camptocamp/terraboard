package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/camptocamp/terraboard/db"
	"github.com/camptocamp/terraboard/s3"
	"github.com/camptocamp/terraboard/util"
)

var states []string

func JSONError(w http.ResponseWriter, message string, err error) {
	errObj := make(map[string]string)
	errObj["error"] = message
	errObj["details"] = fmt.Sprintf("%v", err)
	j, _ := json.Marshal(errObj)
	io.WriteString(w, string(j))
}

func ListStates(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	states := d.ListStates()

	j, err := json.Marshal(states)
	if err != nil {
		JSONError(w, "Failed to marshal states", err)
	}
	io.WriteString(w, string(j))
}

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
	}
	io.WriteString(w, string(j))
}

func GetState(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := util.TrimBase(r, "api/state/")
	versionId := r.URL.Query().Get("versionid")
	var err error
	if versionId == "" {
		versionId, err = d.DefaultVersion(st)
		if err != nil {
			JSONError(w, "Failed to retrieve default version", err)
		}
	}
	state := d.GetState(st, versionId)

	jState, err := json.Marshal(state)
	if err != nil {
		JSONError(w, "Failed to marshal state", err)
	}
	io.WriteString(w, string(jState))
}

func GetStateActivity(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := util.TrimBase(r, "api/state/activity/")
	activity := d.GetStateActivity(st)

	jActivity, err := json.Marshal(activity)
	if err != nil {
		JSONError(w, "Failed to marshal state activity", err)
	}
	io.WriteString(w, string(jActivity))
}

func GetLocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	locks, err := s3.GetLocks()
	if err != nil {
		JSONError(w, "Failed to get locks", err)
	}

	j, err := json.Marshal(locks)
	if err != nil {
		JSONError(w, "Failed to marshal locks", err)
	}
	io.WriteString(w, string(j))
}

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
	}
	io.WriteString(w, string(j))
}

func ListResourceTypes(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, _ := d.ListResourceTypes()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
	}
	io.WriteString(w, string(j))
}

func ListResourceNames(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, _ := d.ListResourceNames()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
	}
	io.WriteString(w, string(j))
}

func ListAttributeKeys(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resourceType := r.URL.Query().Get("resource_type")
	result, _ := d.ListAttributeKeys(resourceType)
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
	}
	io.WriteString(w, string(j))
}
