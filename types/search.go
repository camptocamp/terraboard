package types

import "time"

/**********************************************
 * Search types
 *
 * Used to cast DB requests as search results
 **********************************************/

// SearchResult returns a single search result
type SearchResult struct {
	Path           string `gorm:"column:path" json:"path"`
	VersionID      string `gorm:"column:version_id" json:"version_id"`
	TFVersion      string `gorm:"column:tf_version" json:"tf_version"`
	Serial         int64  `gorm:"column:serial" json:"serial"`
	LineageValue   string `json:"lineage_value"`
	ModulePath     string `gorm:"column:module_path" json:"module_path"`
	ResourceType   string `gorm:"column:type" json:"resource_type"`
	ResourceName   string `gorm:"column:name" json:"resource_name"`
	ResourceIndex  string `gorm:"column:index" json:"resource_index"`
	AttributeKey   string `gorm:"column:key" json:"attribute_key"`
	AttributeValue string `gorm:"column:value" json:"attribute_value"`
}

// StateStat stores State stats
// NOTE: do we want to merge this with StateInfo?
type StateStat struct {
	Path          string    `json:"path"`
	LineageValue  string    `json:"lineage_value"`
	TFVersion     string    `json:"terraform_version"`
	Serial        int64     `json:"serial"`
	VersionID     string    `json:"version_id"`
	LastModified  time.Time `json:"last_modified"`
	ResourceCount int       `json:"resource_count"`
}

// LineageStat stores Lineage stats
type LineageStat struct {
	LineageValue string `json:"lineage_value"`
	StateCount   int    `json:"state_count"`
}
