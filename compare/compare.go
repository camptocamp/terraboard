package compare

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/camptocamp/terraboard/db"
)

type StateInfo struct {
	VersionID     string `json:"version_id"`
	ResourceCount int    `json:"resource_count"`
}

type StateCompare struct {
	Stats struct {
		From StateInfo `json:"from"`
		To   StateInfo `json:"to"`
	} `json:"stats"`
	Differences struct {
		OnlyInOld    []string          `json:"only_in_old"`
		OnlyInNew    []string          `json:"only_in_new"`
		ResourceDiff map[string]string `json:"resource_diff"`
	} `json:"differences"`
}

func countResources(state db.State) (count int) {
	for _, m := range state.Modules {
		count += len(m.Resources)
	}
	return
}

func Compare(from, to db.State) (comp StateCompare, err error) {
	if from.Path == "" {
		err = fmt.Errorf("from version is unknown")
		return
	}
	comp.Stats.From = StateInfo{
		VersionID:     from.Version.VersionID,
		ResourceCount: countResources(from),
	}

	if to.Path == "" {
		err = fmt.Errorf("to version is unknown")
		return
	}
	comp.Stats.To = StateInfo{
		VersionID:     to.Version.VersionID,
		ResourceCount: countResources(to),
	}

	log.WithFields(log.Fields{
		"path": from.Path,
		"from": from.Version.VersionID,
		"to":   to.Version.VersionID,
	}).Info("Comparing state versions")

	return
}
