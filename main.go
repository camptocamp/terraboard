package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/camptocamp/terraboard/api"
	"github.com/camptocamp/terraboard/auth"
	"github.com/camptocamp/terraboard/config"
	"github.com/camptocamp/terraboard/db"
	"github.com/camptocamp/terraboard/state"
	"github.com/camptocamp/terraboard/util"
	tfversion "github.com/hashicorp/terraform/version"
	log "github.com/sirupsen/logrus"
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
	idxStr = util.ReplaceBasePath(idxStr, "base href=\"/\"", "base href=\"%s\"")
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

func handleWithStateProvider(apiF func(w http.ResponseWriter, r *http.Request, sp state.Provider), sp state.Provider) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiF(w, r, sp)
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

// Refresh the DB
// This should be the only direct bridge between the state provider and the DB
func refreshDB(syncInterval uint16, d *db.Database, sp state.Provider) {
	interval := time.Duration(syncInterval) * time.Minute
	for {
		log.Infof("Refreshing DB")
		states, err := sp.GetStates()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to retrieve states. Retrying in 1 minute.")
			time.Sleep(interval)
			continue
		}

		statesVersions := d.ListStatesVersions()
		for _, st := range states {
			versions, _ := sp.GetVersions(st)
			for _, v := range versions {
				if _, ok := statesVersions[v.ID]; ok {
					log.WithFields(log.Fields{
						"version_id": v.ID,
					}).Debug("Version is already in the database, skipping")
				} else {
					d.InsertVersion(&v)
				}

				if isKnownStateVersion(statesVersions, v.ID, st) {
					log.WithFields(log.Fields{
						"path":       st,
						"version_id": v.ID,
					}).Debug("State is already in the database, skipping")
					continue
				}
				state, err := sp.GetState(st, v.ID)
				if err != nil {
					log.WithFields(log.Fields{
						"path":       st,
						"version_id": v.ID,
						"error":      err,
					}).Error("Failed to fetch state")
					continue
				}
				d.InsertState(st, v.ID, state)
				if err != nil {
					log.WithFields(log.Fields{
						"path":       st,
						"version_id": v.ID,
						"error":      err,
					}).Error("Failed to insert state in the database")
				}
			}
		}

		log.Debugf("Waiting %d minutes until next DB sync", syncInterval)
		time.Sleep(interval)
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

	util.SetBasePath(c.Web.BaseURL)

	log.Infof("Terraboard v%s (built for Terraform v%s) is starting...", version, tfversion.Version)

	err := c.SetupLogging()
	if err != nil {
		log.Fatal(err)
	}

	// Set up the state provider
	sp, err := state.Configure(c)
	if err != nil {
		log.Fatal(err)
	}

	// Set up auth
	auth.Setup(c)

	// Set up the DB and start S3->DB sync
	database := db.Init(c.DB, c.Log.Level == "debug")
	if c.DB.NoSync {
		log.Infof("Not syncing database, as requested.")
	} else {
		go refreshDB(c.DB.SyncInterval, database, sp)
	}
	defer database.Close()

	// Index is a wildcard for all paths
	http.HandleFunc(util.GetFullPath(""), idx)

	// Serve static files (CSS, JS, images) from dir
	staticFs := http.FileServer(http.Dir("static"))
	http.Handle(util.GetFullPath("static/"), http.StripPrefix(util.GetFullPath("static"), staticFs))

	// Handle API points
	http.HandleFunc(util.GetFullPath("api/version"), getVersion)
	http.HandleFunc(util.GetFullPath("api/user"), api.GetUser)
	http.HandleFunc(util.GetFullPath("api/states"), handleWithDB(api.ListStates, database))
	http.HandleFunc(util.GetFullPath("api/states/stats"), handleWithDB(api.ListStateStats, database))
	http.HandleFunc(util.GetFullPath("api/states/tfversion/count"), handleWithDB(api.ListTerraformVersionsWithCount, database))
	http.HandleFunc(util.GetFullPath("api/state/"), handleWithDB(api.GetState, database))
	http.HandleFunc(util.GetFullPath("api/state/activity/"), handleWithDB(api.GetStateActivity, database))
	http.HandleFunc(util.GetFullPath("api/state/compare/"), handleWithDB(api.StateCompare, database))
	http.HandleFunc(util.GetFullPath("api/locks"), handleWithStateProvider(api.GetLocks, sp))
	http.HandleFunc(util.GetFullPath("api/search/attribute"), handleWithDB(api.SearchAttribute, database))
	http.HandleFunc(util.GetFullPath("api/resource/types"), handleWithDB(api.ListResourceTypes, database))
	http.HandleFunc(util.GetFullPath("api/resource/types/count"), handleWithDB(api.ListResourceTypesWithCount, database))
	http.HandleFunc(util.GetFullPath("api/resource/names"), handleWithDB(api.ListResourceNames, database))
	http.HandleFunc(util.GetFullPath("api/attribute/keys"), handleWithDB(api.ListAttributeKeys, database))
	http.HandleFunc(util.GetFullPath("api/tf_versions"), handleWithDB(api.ListTfVersions, database))

	// Start server
	log.Debugf("Listening on port %d\n", c.Web.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.Web.Port), nil))
}
