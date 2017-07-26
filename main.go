package main

import (
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/camptocamp/terraboard/api"
	"github.com/camptocamp/terraboard/db"
	"github.com/camptocamp/terraboard/util"
)

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

func handleWithDB(apiF func(w http.ResponseWriter, r *http.Request, d *db.Database), d *db.Database) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiF(w, r, d)
	})
}

func main() {
	database := db.Init()
	go api.RefreshDB(database)
	defer database.Close()

	// Index is a wildcard for all paths
	http.HandleFunc(util.AddBase(""), idx)

	// Serve static files (CSS, JS, images) from dir
	staticFs := http.FileServer(http.Dir("static"))
	http.Handle(util.AddBase("static/"), http.StripPrefix(util.AddBase("static"), staticFs))

	// Handle API points
	http.HandleFunc(util.AddBase("api/states"), api.ListStates)
	http.HandleFunc(util.AddBase("api/state/"), handleWithDB(api.GetState, database))
	http.HandleFunc(util.AddBase("api/history/"), api.GetHistory)
	http.HandleFunc(util.AddBase("api/search/attribute"), handleWithDB(api.SearchAttribute, database))
	http.HandleFunc(util.AddBase("api/resource/types"), handleWithDB(api.ListResourceTypes, database))
	http.HandleFunc(util.AddBase("api/resource/names"), handleWithDB(api.ListResourceNames, database))
	http.HandleFunc(util.AddBase("api/attribute/keys"), handleWithDB(api.ListAttributeKeys, database))

	// Start server
	http.ListenAndServe(":80", nil)
}
