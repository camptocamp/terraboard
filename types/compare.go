package types

/*******************************************************
 * Compare types
 *
 * Used to compute the diff between two state versions
 *******************************************************/

// StateInfo stores general information and stats for a State
type StateInfo struct {
	Path          string `json:"path"`
	VersionID     string `json:"version_id"`
	ResourceCount int    `json:"resource_count"`
	TFVersion     string `json:"terraform_version"`
	Serial        int64  `json:"serial"`
}

// ResourceDiff represents a diff between two versions of a Resource
type ResourceDiff struct {
	OnlyInOld   map[string]string `json:"only_in_old"`
	OnlyInNew   map[string]string `json:"only_in_new"`
	UnifiedDiff string            `json:"unified_diff"`
}

// StateCompare represents a diff between two versions of a State
type StateCompare struct {
	Stats struct {
		From StateInfo `json:"from"`
		To   StateInfo `json:"to"`
	} `json:"stats"`
	Differences struct {
		OnlyInOld    map[string]string       `json:"only_in_old"`
		OnlyInNew    map[string]string       `json:"only_in_new"`
		InBoth       []string                `json:"in_both"`
		ResourceDiff map[string]ResourceDiff `json:"resource_diff"`
	} `json:"differences"`
}
