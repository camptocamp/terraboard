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
	ID           string
	LastModified time.Time
}

// Provider is an interface for supported state providers
type Provider interface {
	GetLocks() (map[string]LockInfo, error)
	GetVersions(string) ([]Version, error)
	GetStates() ([]string, error)
	GetState(string, string) (*statefile.File, error)
}

// Configure the state provider
func Configure(c *config.Config) (Provider, error) {
	if len(c.TFE.Token) > 0 {
		log.Info("Using Terraform Enterprise as the state/locks provider")
		provider, err := NewTFE(c)
		if err != nil {
			return nil, err
		}
		return &provider, nil
	}

	if c.GCP.GCSBuckets != nil {
		log.Info("Using Google Cloud as the state/locks provider")
		provider, err := NewGCP(c)
		if err != nil {
			return nil, err
		}
		return &provider, nil
	}

	log.Info("Using AWS (S3+DynamoDB) as the state/locks provider")
	provider := NewAWS(c)
	return &provider, nil
}
