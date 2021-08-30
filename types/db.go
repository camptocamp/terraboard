package types

import (
	"database/sql"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
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
	TFVersion  string        `gorm:"varchar(10)" json:"terraform_version"`
	Serial     int64         `json:"serial"`
	LineageID  sql.NullInt64 `gorm:"index" json:"-"`
	Modules    []Module      `json:"modules"`
}

type Lineage struct {
	gorm.Model
	Value  string  `gorm:"index;unique" json:"lineage"`
	States []State `json:"states"`
	Plans  []Plan  `json:"plans"`
}

// Module is a Terraform module in a State
type Module struct {
	ID           uint          `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	StateID      sql.NullInt64 `gorm:"index" json:"-"`
	Path         string        `json:"path"`
	Resources    []Resource    `json:"resources"`
	OutputValues []OutputValue `json:"outputs"`
}

// Resource is a Terraform resource in a Module
type Resource struct {
	ID         uint          `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	ModuleID   sql.NullInt64 `gorm:"index" json:"-"`
	Type       string        `gorm:"index" json:"type"`
	Name       string        `gorm:"index" json:"name"`
	Index      string        `gorm:"index" json:"index"`
	Attributes []Attribute   `json:"attributes"`
}

// OutputValue is a Terraform output in a Module
type OutputValue struct {
	ID        uint          `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	ModuleID  sql.NullInt64 `gorm:"index" json:"-"`
	Sensitive bool          `gorm:"index" json:"sensitive"`
	Name      string        `gorm:"index" json:"name"`
	Value     string        `json:"value"`
}

// Attribute is a Terraform attribute in a Resource
type Attribute struct {
	ID         uint          `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	ResourceID sql.NullInt64 `gorm:"index" json:"-"`
	Key        string        `gorm:"index" json:"key"`
	Value      string        `json:"value"`
}

// Plan is a Terraform plan
type Plan struct {
	gorm.Model
	LineageID    uint           `gorm:"index" json:"-"`
	Lineage      Lineage        `json:"lineage_data"`
	TFVersion    string         `gorm:"varchar(10)" json:"terraform_version"`
	GitRemote    string         `json:"git_remote"`
	GitCommit    string         `gorm:"varchar(50)" json:"git_commit"`
	CiURL        string         `json:"ci_url"`
	Source       string         `json:"source"`
	ParsedPlan   PlanModel      `json:"parsed_plan"`
	ParsedPlanID sql.NullInt64  `gorm:"index" json:"-"`
	PlanJSON     datatypes.JSON `json:"plan_json"`
}

// PlanModel represents the entire contents of an output Terraform plan.
type PlanModel struct {
	gorm.Model
	// The version of the plan format. This should always match the
	// PlanFormatVersion constant in this package, or else an unmarshal
	// will be unstable.
	FormatVersion string `json:"format_version,omitempty"`

	// The version of Terraform used to make the plan.
	TerraformVersion string `json:"terraform_version,omitempty"`

	// The common state representation of resources within this plan.
	// This is a product of the existing state merged with the diff for
	// this plan.
	PlanStateValue   PlanStateValue `json:"planned_values,omitempty"`
	PlanStateValueID sql.NullInt64  `gorm:"index" json:"-"`

	// The variables set in the root module when creating the plan.
	Variables planVariableList `json:"variables,omitempty"`

	// The change operations for resources and data sources within this
	// plan.
	PlanResourceChanges []PlanResourceChange `json:"resource_changes,omitempty"`

	// The change operations for outputs within this plan.
	PlanOutputs planOutputList `json:"output_changes,omitempty"`

	// The Terraform state prior to the plan operation. This is the
	// same format as PlannedValues, without the current diff merged.
	PlanState   PlanState     `json:"prior_state,omitempty"`
	PlanStateID sql.NullInt64 `gorm:"index" json:"-"`

	// The Terraform configuration used to make the plan.
	// PlanConfig   PlanConfig    `json:"configuration,omitempty"`
	// PlanConfigID sql.NullInt64 `gorm:"index" json:"-"`
}

// PlanState is the top-level representation of a Terraform state.
type PlanState struct {
	gorm.Model

	// The version of the state format. This should always match the
	// StateFormatVersion constant in this package, or else am
	// unmarshal will be unstable.
	FormatVersion string `json:"format_version,omitempty"`

	// The Terraform version used to make the state.
	TerraformVersion string `json:"terraform_version,omitempty"`

	// The values that make up the state.
	PlanStateValue   PlanStateValue `json:"values,omitempty"`
	PlanStateValueID sql.NullInt64  `gorm:"index" json:"-"`
}

// PlanStateValue is the common representation of resolved values for both the
// prior state (which is always complete) and the planned new state.
type PlanStateValue struct {
	gorm.Model
	// The Outputs for this common state representation.
	PlanStateOutputs planStateOutputList `json:"outputs,omitempty"`

	// The root module in this state representation.
	PlanStateModule   PlanStateModule `json:"root_module,omitempty"`
	PlanStateModuleID sql.NullInt64   `gorm:"index" json:"-"`
}

// PlanStateOutput is the representation of a state output in a plan
type PlanStateOutput struct {
	gorm.Model
	PlanStateValueID sql.NullInt64 `gorm:"index" json:"-"`
	Name             string        `gorm:"index" json:"name"`
	Sensitive        bool          `json:"sensitive"`
	Value            string        `json:"value"`
}

// PlanStateModule is the representation of a module in the common state
// representation. This can be the root module or a child module.
type PlanStateModule struct {
	gorm.Model
	// All resources or data sources within this module.
	PlanStateResources []PlanStateResource `json:"resources,omitempty"`

	// The absolute module address, omitted for the root module.
	Address string `json:"address,omitempty"`

	// Any child modules within this module.
	PlanStateModules  []PlanStateModule `json:"child_modules,omitempty"`
	PlanStateModuleID sql.NullInt64     `gorm:"index" json:"-"`
}

// PlanStateResource is the representation of a resource in the common
// state representation.
type PlanStateResource struct {
	gorm.Model
	PlanStateModuleID sql.NullInt64 `gorm:"index" json:"-"`

	// The absolute resource address.
	Address string `json:"address,omitempty"`

	// The resource mode.
	Mode string `json:"mode,omitempty"`

	// The resource type, example: "aws_instance" for aws_instance.foo.
	Type string `json:"type,omitempty"`

	// The resource name, example: "foo" for aws_instance.foo.
	Name string `json:"name,omitempty"`

	// The instance key for any resources that have been created using
	// "count" or "for_each". If neither of these apply the key will be
	// empty.
	//
	// This value can be either an integer (int) or a string.
	Index rawJSON `json:"index,omitempty"`

	// The name of the provider this resource belongs to. This allows
	// the provider to be interpreted unambiguously in the unusual
	// situation where a provider offers a resource type whose name
	// does not start with its own name, such as the "googlebeta"
	// provider offering "google_compute_instance".
	ProviderName string `json:"provider_name,omitempty"`

	//  The version of the resource type schema the "values" property
	//  conforms to.
	SchemaVersion uint `json:"schema_version,omitempty"`

	// The JSON representation of the attribute values of the resource,
	// whose structure depends on the resource type schema. Any unknown
	// values are omitted or set to null, making them indistinguishable
	// from absent values.
	PlanStateResourceAttributes planStateResourceAttributeList `json:"values,omitempty"`

	// The addresses of the resources that this resource depends on.
	DependsOn rawJSON `json:"depends_on,omitempty"`

	// If true, the resource has been marked as tainted and will be
	// re-created on the next update.
	Tainted bool `json:"tainted,omitempty"`

	// DeposedKey is set if the resource instance has been marked Deposed and
	// will be destroyed on the next apply.
	DeposedKey string `json:"deposed_key,omitempty"`
}

type PlanStateResourceAttribute struct {
	gorm.Model
	PlanStateResourceID sql.NullInt64 `gorm:"index" json:"-"`
	Key                 string        `gorm:"index" json:"key"`
	Value               string        `json:"value"`
}

// PlanResourceChange is a description of an individual change action
// that Terraform plans to use to move from the prior state to a new
// state matching the configuration.
type PlanResourceChange struct {
	gorm.Model
	PlanModelID sql.NullInt64 `gorm:"index" json:"-"`

	// The absolute resource address.
	Address string `json:"address,omitempty"`

	// The module portion of the above address. Omitted if the instance
	// is in the root module.
	ModuleAddress string `json:"module_address,omitempty"`

	// The resource mode.
	Mode string `json:"mode,omitempty"`

	// The resource type, example: "aws_instance" for aws_instance.foo.
	Type string `json:"type,omitempty"`

	// The resource name, example: "foo" for aws_instance.foo.
	Name string `json:"name,omitempty"`

	// The instance key for any resources that have been created using
	// "count" or "for_each". If neither of these apply the key will be
	// empty.
	//
	// This value can be either an integer (int) or a string.
	Index rawJSON `json:"index,omitempty"`

	// The name of the provider this resource belongs to. This allows
	// the provider to be interpreted unambiguously in the unusual
	// situation where a provider offers a resource type whose name
	// does not start with its own name, such as the "googlebeta"
	// provider offering "google_compute_instance".
	ProviderName string `json:"provider_name,omitempty"`

	// An identifier used during replacement operations, and can be
	// used to identify the exact resource being replaced in state.
	DeposedKey string `json:"deposed,omitempty"`

	// The data describing the change that will be made to this object.
	Change   Change        `json:"change,omitempty"`
	ChangeID sql.NullInt64 `gorm:"index" json:"-"`
}

type PlanOutput struct {
	gorm.Model
	Name        string        `gorm:"index" json:"name"`
	PlanModelID sql.NullInt64 `gorm:"index" json:"-"`
	// The data describing the change that will be made to this object.
	Change   Change        `json:"change,omitempty"`
	ChangeID sql.NullInt64 `gorm:"index" json:"-"`
}

// Change is the representation of a proposed change for an object.
type Change struct {
	gorm.Model
	// The action to be carried out by this change.
	Actions rawJSON `json:"actions,omitempty"`

	// Before and After are representations of the object value both
	// before and after the action. For create and delete actions,
	// either Before or After is unset (respectively). For no-op
	// actions, both values will be identical. After will be incomplete
	// if there are values within it that won't be known until after
	// apply.
	Before rawJSON `json:"before,omitempty"`
	After  rawJSON `json:"after,omitempty"`

	// A deep object of booleans that denotes any values that are
	// unknown in a resource. These values were previously referred to
	// as "computed" values.
	//
	// If the value cannot be found in this map, then its value should
	// be available within After, so long as the operation supports it.
	AfterUnknown rawJSON `json:"after_unknown,omitempty"`

	// BeforeSensitive and AfterSensitive are object values with similar
	// structure to Before and After, but with all sensitive leaf values
	// replaced with true, and all non-sensitive leaf values omitted. These
	// objects should be combined with Before and After to prevent accidental
	// display of sensitive values in user interfaces.
	BeforeSensitive rawJSON `json:"before_sensitive,omitempty"`
	AfterSensitive  rawJSON `json:"after_sensitive,omitempty"`
}

// PlanModelVariable is a top-level variable in the Terraform plan.
type PlanModelVariable struct {
	gorm.Model
	PlanModelID sql.NullInt64 `gorm:"index" json:"-"`
	Key         string        `gorm:"index" json:"key"`
	Value       string        `json:"value,omitempty"`
}
