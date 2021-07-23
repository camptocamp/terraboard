package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/camptocamp/terraboard/api"
	"github.com/camptocamp/terraboard/auth"
	"github.com/camptocamp/terraboard/config"
	"github.com/camptocamp/terraboard/db"
	"github.com/camptocamp/terraboard/state"
	"github.com/camptocamp/terraboard/util"
	"github.com/gorilla/mux"
	tfversion "github.com/hashicorp/terraform/version"
	log "github.com/sirupsen/logrus"
)

// Pass the DB to API handlers
// This takes a callback and returns a HandlerFunc
// which calls the callback with the DB
func handleWithDB(apiF func(w http.ResponseWriter, r *http.Request,
	d *db.Database), d *db.Database) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiF(w, r, d)
	})
}

func handleWithStateProviders(apiF func(w http.ResponseWriter, r *http.Request,
	sps []state.Provider), sps []state.Provider) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiF(w, r, sps)
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
// This should be the only direct bridge between the state providers and the DB
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
			for k, v := range versions {
				if _, ok := statesVersions[v.ID]; ok {
					log.WithFields(log.Fields{
						"version_id": v.ID,
					}).Debug("Version is already in the database, skipping")
				} else {
					if err := d.InsertVersion(&versions[k]); err != nil {
						log.Error(err.Error())
					}
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
					}).Error("Failed to fetch state from bucket")
					continue
				}
				if err = d.InsertState(st, v.ID, state); err != nil {
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

func getVersion(w http.ResponseWriter, _ *http.Request) {
	j, err := json.Marshal(map[string]string{
		"version":   version,
		"copyright": "Copyright © 2017-2021 Camptocamp",
	})
	if err != nil {
		api.JSONError(w, "Failed to marshal version", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
		next.ServeHTTP(w, r)
	})
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
	sps, err := state.Configure(c)
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
		log.Debugf("Total providers: %d\n", len(sps))
		for _, sp := range sps {
			go refreshDB(c.DB.SyncInterval, database, sp)
		}
	}
	defer database.Close()

	// Instantiate gorilla/mux router instance
	r := mux.NewRouter()

	// Handle API endpoints
	apiRouter := r.PathPrefix("/api/").Subrouter()
	apiRouter.HandleFunc(util.GetFullPath("version"), getVersion)
	apiRouter.HandleFunc(util.GetFullPath("user"), api.GetUser)
	apiRouter.HandleFunc(util.GetFullPath("lineages"), handleWithDB(api.GetLineages, database))
	apiRouter.HandleFunc(util.GetFullPath("lineages/stats"), handleWithDB(api.ListStateStats, database))
	apiRouter.HandleFunc(util.GetFullPath("lineages/tfversion/count"),
		handleWithDB(api.ListTerraformVersionsWithCount, database))
	apiRouter.HandleFunc(util.GetFullPath("lineages/{lineage}"), handleWithDB(api.GetState, database))
	apiRouter.HandleFunc(util.GetFullPath("lineages/{lineage}/activity"), handleWithDB(api.GetLineageActivity, database))
	apiRouter.HandleFunc(util.GetFullPath("lineages/{lineage}/compare"), handleWithDB(api.StateCompare, database))
	apiRouter.HandleFunc(util.GetFullPath("locks"), handleWithStateProviders(api.GetLocks, sps))
	apiRouter.HandleFunc(util.GetFullPath("search/attribute"), handleWithDB(api.SearchAttribute, database))
	apiRouter.HandleFunc(util.GetFullPath("resource/types"), handleWithDB(api.ListResourceTypes, database))
	apiRouter.HandleFunc(util.GetFullPath("resource/types/count"), handleWithDB(api.ListResourceTypesWithCount, database))
	apiRouter.HandleFunc(util.GetFullPath("resource/names"), handleWithDB(api.ListResourceNames, database))
	apiRouter.HandleFunc(util.GetFullPath("attribute/keys"), handleWithDB(api.ListAttributeKeys, database))
	apiRouter.HandleFunc(util.GetFullPath("tf_versions"), handleWithDB(api.ListTfVersions, database))
	apiRouter.HandleFunc(util.GetFullPath("plans"), handleWithDB(api.ManagePlans, database))

	// Serve static files (CSS, JS, images) from dir
	staticFs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFs))

	// Serve index page on all unhandled routes
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// Add CORS Middleware to mux router
	r.Use(corsMiddleware)

	// Start server
	log.Debugf("Listening on port %d\n", c.Web.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.Web.Port), r))
}
