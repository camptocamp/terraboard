package compare

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/camptocamp/terraboard/db"
)

type StateCompare struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func Compare(from, to db.State) (comp StateCompare, err error) {
	if from.Path == "" {
		err = fmt.Errorf("from version is unknown")
		return
	}
	comp.From = from.Version.VersionID

	if to.Path == "" {
		err = fmt.Errorf("to version is unknown")
		return
	}
	comp.To = to.Version.VersionID

	log.WithFields(log.Fields{
		"path": from.Path,
		"from": from.Version.VersionID,
		"to":   to.Version.VersionID,
	}).Info("Comparing state versions")

	return
}
