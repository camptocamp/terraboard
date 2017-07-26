package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/camptocamp/terraboard/db"
	"github.com/camptocamp/terraboard/s3"
	"github.com/camptocamp/terraboard/util"
)

var states []string

func RefreshDB(d *db.Database) {
	for {
		log.Infof("Refreshing DB from S3")
		states, err := s3.GetStates()
		if err != nil {
			log.Errorf("Failed to build cache: %s", err)
		}

		for _, st := range states {
			versions, _ := s3.GetVersions(st)
			for _, v := range versions {
				d.InsertVersion(v)

				s := d.GetState(st, *v.VersionId)
				if s.Path == st {
					log.Infof("State %s/%s is already in the DB, skipping", st, *v.VersionId)
					continue
				}
				state, _ := s3.GetState(st, *v.VersionId)
				d.InsertState(st, *v.VersionId, state)
				if err != nil {
					log.Errorf("Failed to insert state %s/%s: %v", st, *v.VersionId, err)
				}
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func ApiStates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	states, err := s3.GetStates()
	if err != nil {
		errObj := make(map[string]string)
		errObj["error"] = "Failed to list states"
		errObj["details"] = fmt.Sprintf("%v", err)
		j, _ := json.Marshal(errObj)
		io.WriteString(w, string(j))
		return
	}

	j, _ := json.Marshal(states)
	io.WriteString(w, string(j))
}

func ApiState(w http.ResponseWriter, r *http.Request, d *db.Database) {
	st := util.TrimBase(r, "api/state/")
	versionId := r.URL.Query().Get("versionid")
	if versionId == "" {
		versionId, _ = d.DefaultVersion(st) // TODO: err
	}
	state := d.GetState(st, versionId)

	jState, _ := json.Marshal(state)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, string(jState))
}

func ApiHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := util.TrimBase(r, "api/history/")
	result, err := s3.GetVersions(st)
	if err != nil {
		errObj := make(map[string]string)
		errObj["error"] = fmt.Sprintf("State file history not found: %v", st)
		errObj["details"] = fmt.Sprintf("%v", err)
		j, err := json.Marshal(errObj)
		if err != nil {
			log.Errorf("Failed to marshal json: %v", err)
		}
		io.WriteString(w, string(j))
		return
	}

	j, err := json.Marshal(result)
	if err != nil {
		log.Errorf("Failed to marshal json: %v", err)
	}
	io.WriteString(w, string(j))
}

func ApiSearchAttribute(w http.ResponseWriter, r *http.Request, d *db.Database) {
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
		log.Errorf("Failed to marshal json: %v", err)
	}
	io.WriteString(w, string(j))
}

func ApiResourceTypes(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, _ := d.ListResourceTypes()
	j, err := json.Marshal(result)
	if err != nil {
		log.Errorf("Failed to marshal json: %v", err)
	}
	io.WriteString(w, string(j))
}

func ApiResourceNames(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, _ := d.ListResourceNames()
	j, err := json.Marshal(result)
	if err != nil {
		log.Errorf("Failed to marshal json: %v", err)
	}
	io.WriteString(w, string(j))
}

func ApiAttributeKeys(w http.ResponseWriter, r *http.Request, d *db.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resourceType := r.URL.Query().Get("resource_type")
	result, _ := d.ListAttributeKeys(resourceType)
	j, err := json.Marshal(result)
	if err != nil {
		log.Errorf("Failed to marshal json: %v", err)
	}
	io.WriteString(w, string(j))
}
