package types

/*******************************************************
 * Compare types
 *
 * Used to compute the diff between two state versions
 *******************************************************/

type StateInfo struct {
	VersionID     string `json:"version_id"`
	ResourceCount int    `json:"resource_count"`
}

type ResourceDiff struct {
	OnlyInOld   map[string]string `json:"only_in_old"`
	OnlyInNew   map[string]string `json:"only_in_new"`
	UnifiedDiff string            `json:"unified_diff"`
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
