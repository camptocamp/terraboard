package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/camptocamp/terraboard/api"
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
	idx, err := ioutil.ReadFile("index.html")
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

// Refresh the DB from S3
// This should be the only direct bridge between S3 and the DB
func refreshDB(d *db.Database) {
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

var version = "undefined"

// Main
func main() {
	log.Infof("Terraboard v%s is starting...", version)

	c := config.LoadConfig(version)

	// Set up S3
	s3.Setup(c)

	// Set up the DB and start S3->DB sync
	database := db.Init(
		c.DB.Host, c.DB.User,
		c.DB.Name, c.DB.Password)
	go refreshDB(database)
	defer database.Close()

	// Index is a wildcard for all paths
	http.HandleFunc(util.AddBase(""), idx)

	// Serve static files (CSS, JS, images) from dir
	staticFs := http.FileServer(http.Dir("static"))
	http.Handle(util.AddBase("static/"), http.StripPrefix(util.AddBase("static"), staticFs))

	// Handle API points
	http.HandleFunc(util.AddBase("api/states"), handleWithDB(api.ListStates, database))
	http.HandleFunc(util.AddBase("api/state/"), handleWithDB(api.GetState, database))
	http.HandleFunc(util.AddBase("api/history/"), api.GetHistory)
	http.HandleFunc(util.AddBase("api/search/attribute"), handleWithDB(api.SearchAttribute, database))
	http.HandleFunc(util.AddBase("api/resource/types"), handleWithDB(api.ListResourceTypes, database))
	http.HandleFunc(util.AddBase("api/resource/names"), handleWithDB(api.ListResourceNames, database))
	http.HandleFunc(util.AddBase("api/attribute/keys"), handleWithDB(api.ListAttributeKeys, database))

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.Port), nil))
}
