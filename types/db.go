package types

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
)

/*********************************************
 * Database object types
 *
 * Each type corresponds to a table in the DB
 *********************************************/

// Version is an S3 bucket version
type Version struct {
	ID           uint      `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	VersionID    string    `gorm:"index" json:"version_id"`
	LastModified time.Time `json:"last_modified"`
}

// State is a Terraform State
type State struct {
	gorm.Model `json:"-"`
	Path       string        `gorm:"index" json:"path"`
	Version    Version       `json:"version"`
	VersionID  sql.NullInt64 `gorm:"index" json:"-"`
	TFVersion  string        `json:"terraform_version"`
	Serial     int64         `json:"serial"`
	Modules    []Module      `json:"modules"`
}

// Module is a Terraform module in a State
type Module struct {
	ID        uint          `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	StateID   sql.NullInt64 `gorm:"index" json:"-"`
	Path      string        `json:"path"`
	Resources []Resource    `json:"resources"`
}

// Resource is a Terraform resource in a Module
type Resource struct {
	ID         uint          `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	ModuleID   sql.NullInt64 `gorm:"index" json:"-"`
	Type       string        `gorm:"index" json:"type"`
	Name       string        `gorm:"index" json:"name"`
	Attributes []Attribute   `json:"attributes"`
}

// Attribute is a Terraform attribute in a Resource
type Attribute struct {
	ID         uint          `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	ResourceID sql.NullInt64 `gorm:"index" json:"-"`
	Key        string        `gorm:"index" json:"key"`
	Value      string        `json:"value"`
}
