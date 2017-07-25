package main

import (
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/camptocamp/terraboard/api"
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

func main() {
	// Index is a wildcard for all paths
	http.HandleFunc(util.AddBase(""), idx)

	// Serve static files (CSS, JS, images) from dir
	staticFs := http.FileServer(http.Dir("static"))
	http.Handle(util.AddBase("static/"), http.StripPrefix(util.AddBase("static"), staticFs))

	// Handle API points
	http.HandleFunc(util.AddBase("api/states"), api.ApiStates)
	http.HandleFunc(util.AddBase("api/state/"), api.ApiState)
	http.HandleFunc(util.AddBase("api/history/"), api.ApiHistory)
	http.HandleFunc(util.AddBase("api/search/resource"), api.ApiSearchResource)
	http.HandleFunc(util.AddBase("api/search/attribute"), api.ApiSearchAttribute)
	http.HandleFunc(util.AddBase("api/resource/types"), api.ApiResourceTypes)
	http.HandleFunc(util.AddBase("api/resource/names"), api.ApiResourceNames)
	http.HandleFunc(util.AddBase("api/attribute/keys"), api.ApiAttributeKeys)

	// Start server
	http.ListenAndServe(":80", nil)
}
