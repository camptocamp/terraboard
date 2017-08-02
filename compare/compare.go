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
		InBoth       []string          `json:"in_both"`
		ResourceDiff map[string]string `json:"resource_diff"`
	} `json:"differences"`
}

// Return all resources of a state
func stateResources(state db.State) (res []string) {
	for _, m := range state.Modules {
		for _, r := range m.Resources {
			res = append(res, fmt.Sprintf("%s.%s.%s", m.Path, r.Type, r.Name))
		}
	}
	return
}

// Returns elements only in s1
func sliceDiff(s1, s2 []string) (diff []string) {
	for _, e1 := range s1 {
		found := false
		for _, e2 := range s2 {
			if e1 == e2 {
				found = true
				break
			}
		}

		if !found {
			diff = append(diff, e1)
		}
	}
	return
}

// Returns elements in both s1 and s2
func sliceInter(s1, s2 []string) (inter []string) {
	for _, e1 := range s1 {
		for _, e2 := range s2 {
			if e1 == e2 {
				inter = append(inter, e1)
				break
			}
		}
	}
	return
}

func Compare(from, to db.State) (comp StateCompare, err error) {
	if from.Path == "" {
		err = fmt.Errorf("from version is unknown")
		return
	}
	fromResources := stateResources(from)
	comp.Stats.From = StateInfo{
		VersionID:     from.Version.VersionID,
		ResourceCount: len(fromResources),
	}

	if to.Path == "" {
		err = fmt.Errorf("to version is unknown")
		return
	}
	toResources := stateResources(to)
	comp.Stats.To = StateInfo{
		VersionID:     to.Version.VersionID,
		ResourceCount: len(toResources),
	}

	comp.Differences.OnlyInOld = sliceDiff(fromResources, toResources)
	comp.Differences.OnlyInNew = sliceDiff(toResources, fromResources)
	comp.Differences.InBoth = sliceInter(toResources, fromResources)

	log.WithFields(log.Fields{
		"path": from.Path,
		"from": from.Version.VersionID,
		"to":   to.Version.VersionID,
	}).Info("Comparing state versions")

	return
}
