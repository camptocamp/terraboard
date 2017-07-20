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
	http.HandleFunc(util.AddBase(""), idx)
	staticFs := http.FileServer(http.Dir("static"))
	http.Handle(util.AddBase("static/"), http.StripPrefix(util.AddBase("static"), staticFs))
	http.HandleFunc(util.AddBase("api/states"), api.States)
	http.HandleFunc(util.AddBase("api/state/"), api.State)
	http.HandleFunc(util.AddBase("api/history/"), api.History)
	http.ListenAndServe(":80", nil)
}
