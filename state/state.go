package state

import (
	"time"

	"github.com/camptocamp/terraboard/config"
	"github.com/hashicorp/terraform/states/statefile"
	log "github.com/sirupsen/logrus"
)

// LockInfo stores information on a State Lock
type LockInfo struct {
	ID        string
	Operation string
	Info      string
	Who       string
	Version   string
	Created   *time.Time
	Path      string
}

// Lock is a single State Lock
type Lock struct {
	LockID string
	Info   string
}

// Version is a handler for state versions
type Version struct {
	ID string
	LastModified time.Time
}

// Provider is an interface for supported state providers
type Provider interface {
	GetLocks() (map[string]LockInfo, error)
	GetVersions(string) ([]Version, error)
	GetStates() ([]string, error)
	GetState(string, string) (*statefile.File, error)
}

// Provider is a handler for the configured one
var p Provider

// Configure returns the configured provider
func Configure(c *config.Config) {
	log.Info("Using AWS (S3+DynamoDB) as the state/locks provider")
	provider := NewAWS(c)
	p = &provider
}

// Functions wrappers
//
func GetLocks() (map[string]LockInfo, error) {
	return p.GetLocks()
}

func GetVersions(state string) ([]Version, error) {
	return p.GetVersions(state)
}

func GetStates() ([]string, error) {
	return p.GetStates()
}

func GetState(state, versionID string) (*statefile.File, error) {
	return p.GetState(state, versionID)
}
