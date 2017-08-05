package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/camptocamp/terraboard/api"
	"github.com/camptocamp/terraboard/auth"
	"github.com/camptocamp/terraboard/config"
	"github.com/camptocamp/terraboard/db"
	"github.com/camptocamp/terraboard/s3"
	"github.com/camptocamp/terraboard/util"
)

// idx serves index.html, always,
// so as to let AngularJS manage the app routing.
// The <base> HTML tag is edited on the fly
// to reflect the proper base URL
func idx(w http.ResponseWriter, r *http.Request) {
	idx, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		log.Errorf("Failed to open index.html: %v", err)
		// TODO: Return error page
	}
	idxStr := string(idx)
	idxStr = util.ReplaceBase(idxStr, "base href=\"/\"", "base href=\"%s\"")
	io.WriteString(w, idxStr)
}

// Pass the DB to API handlers
// This takes a callback and returns a HandlerFunc
// which calls the callback with the DB
func handleWithDB(apiF func(w http.ResponseWriter, r *http.Request, d *db.Database), d *db.Database) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiF(w, r, d)
	})
}

func isKnownStateVersion(statesVersions map[string][]string, versionID, path string) bool {
	if v, ok := statesVersions[versionID]; ok {
		for _, s := range v {
			if s == path {
				return true
			}
		}
	}
	return false
}

// Refresh the DB from S3
// This should be the only direct bridge between S3 and the DB
func refreshDB(d *db.Database) {
	for {
		log.Infof("Refreshing DB from S3")
		states, err := s3.GetStates()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to retrieve states from S3. Retrying in 1 minute.")
			time.Sleep(1 * time.Minute)
			continue
		}

		statesVersions := d.ListStatesVersions()
		for _, st := range states {
			versions, _ := s3.GetVersions(st)
			for _, v := range versions {
				if _, ok := statesVersions[*v.VersionId]; ok {
					log.WithFields(log.Fields{
						"version_id": *v.VersionId,
					}).Debug("Version is already in the database, skipping")
				} else {
					d.InsertVersion(v)
				}

				if isKnownStateVersion(statesVersions, *v.VersionId, st) {
					log.WithFields(log.Fields{
						"path":       st,
						"version_id": *v.VersionId,
					}).Debug("State is already in the database, skipping")
					continue
				}
				state, _ := s3.GetState(st, *v.VersionId)
				d.InsertState(st, *v.VersionId, state)
				if err != nil {
					log.WithFields(log.Fields{
						"path":       st,
						"version_id": *v.VersionId,
						"error":      err,
					}).Error("Failed to insert state in the database")
				}
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

var version = "undefined"

func getVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	j, err := json.Marshal(map[string]string{
		"version":   version,
		"copyright": "Copyright © 2017 Camptocamp",
	})
	if err != nil {
		api.JSONError(w, "Failed to marshal version", err)
		return
	}
	io.WriteString(w, string(j))
}

// Main
func main() {
	c := config.LoadConfig(version)

	log.Infof("Terraboard v%s is starting...", version)

	err := c.SetupLogging()
	if err != nil {
		log.Fatal(err)
	}

	// Set up S3
	s3.Setup(c)

	// Set up auth
	auth.Setup(c)

	// Set up the DB and start S3->DB sync
	database := db.Init(
		c.DB.Host, c.DB.User,
		c.DB.Name, c.DB.Password,
		c.Log.Level)
	if c.DB.NoSync {
		log.Infof("Not syncing database, as requested.")
	} else {
		go refreshDB(database)
	}
	defer database.Close()

	// Index is a wildcard for all paths
	http.HandleFunc(util.AddBase(""), idx)

	// Serve static files (CSS, JS, images) from dir
	staticFs := http.FileServer(http.Dir("static"))
	http.Handle(util.AddBase("static/"), http.StripPrefix(util.AddBase("static"), staticFs))

	// Handle API points
	http.HandleFunc(util.AddBase("api/version"), getVersion)
	http.HandleFunc(util.AddBase("api/user"), api.GetUser)
	http.HandleFunc(util.AddBase("api/states"), handleWithDB(api.ListStates, database))
	http.HandleFunc(util.AddBase("api/states/stats"), handleWithDB(api.ListStateStats, database))
	http.HandleFunc(util.AddBase("api/states/tfversion/count"), handleWithDB(api.ListTerraformVersionsWithCount, database))
	http.HandleFunc(util.AddBase("api/state/"), handleWithDB(api.GetState, database))
	http.HandleFunc(util.AddBase("api/state/activity/"), handleWithDB(api.GetStateActivity, database))
	http.HandleFunc(util.AddBase("api/state/compare/"), handleWithDB(api.StateCompare, database))
	http.HandleFunc(util.AddBase("api/locks"), api.GetLocks)
	http.HandleFunc(util.AddBase("api/search/attribute"), handleWithDB(api.SearchAttribute, database))
	http.HandleFunc(util.AddBase("api/resource/types"), handleWithDB(api.ListResourceTypes, database))
	http.HandleFunc(util.AddBase("api/resource/types/count"), handleWithDB(api.ListResourceTypesWithCount, database))
	http.HandleFunc(util.AddBase("api/resource/names"), handleWithDB(api.ListResourceNames, database))
	http.HandleFunc(util.AddBase("api/attribute/keys"), handleWithDB(api.ListAttributeKeys, database))
	http.HandleFunc(util.AddBase("api/tf_versions"), handleWithDB(api.ListTfVersions, database))

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.Port), nil))
}
