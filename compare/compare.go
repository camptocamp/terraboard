package compare

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/camptocamp/terraboard/db"
	"github.com/pmezard/go-difflib/difflib"
)

type StateInfo struct {
	VersionID     string `json:"version_id"`
	ResourceCount int    `json:"resource_count"`
}

type ResourceDiff struct {
	OnlyInOld   []string `json:"only_in_old"`
	OnlyInNew   []string `json:"only_in_new"`
	UnifiedDiff string   `json:"unified_diff"`
}

type StateCompare struct {
	Stats struct {
		From StateInfo `json:"from"`
		To   StateInfo `json:"to"`
	} `json:"stats"`
	Differences struct {
		OnlyInOld    []string                `json:"only_in_old"`
		OnlyInNew    []string                `json:"only_in_new"`
		InBoth       []string                `json:"in_both"`
		ResourceDiff map[string]ResourceDiff `json:"resource_diff"`
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

// Return all attributes of a resource
func resourceAttributes(res db.Resource) (attrs []string) {
	for _, a := range res.Attributes {
		attrs = append(attrs, a.Key)
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

func getResource(state db.State, key string) (res db.Resource) {
	for _, m := range state.Modules {
		if strings.HasPrefix(key, m.Path) {
			for _, r := range m.Resources {
				if key == fmt.Sprintf("%s.%s.%s", m.Path, r.Type, r.Name) {
					return r
				}
			}
		} else {
			continue
		}
	}
	return
}

func formatResource(res db.Resource) (out string) {
	out = fmt.Sprintf("resource \"%s\" \"%s\" {\n", res.Type, res.Name)
	for _, attr := range res.Attributes {
		out += fmt.Sprintf("  %s = \"%s\"\n", attr.Key, attr.Value)
	}
	out += "}\n"

	return
}

func stateInfo(state db.State) (info string) {
	return fmt.Sprintf("%s (%s)", state.Path, state.Version.LastModified)
}

// Compare a resource in two states
func compareResource(st1, st2 db.State, key string) (comp ResourceDiff) {
	res1 := getResource(st1, key)
	attrs1 := resourceAttributes(res1)
	res2 := getResource(st2, key)
	attrs2 := resourceAttributes(res2)

	comp.OnlyInOld = sliceDiff(attrs1, attrs2)
	comp.OnlyInNew = sliceDiff(attrs2, attrs1)

	// Compute unified diff
	diff := difflib.ContextDiff{
		A:        difflib.SplitLines(formatResource(res1)),
		B:        difflib.SplitLines(formatResource(res2)),
		FromFile: stateInfo(st1),
		ToFile:   stateInfo(st2),
		Context:  3,
		Eol:      "\n",
	}
	result, _ := difflib.GetContextDiffString(diff)
	comp.UnifiedDiff = result

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
	comp.Differences.ResourceDiff = make(map[string]ResourceDiff)

	for _, r := range comp.Differences.InBoth {
		comp.Differences.ResourceDiff[r] = compareResource(to, from, r)
	}

	log.WithFields(log.Fields{
		"path": from.Path,
		"from": from.Version.VersionID,
		"to":   to.Version.VersionID,
	}).Info("Comparing state versions")

	return
}
