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
func Configure(c *config.Config) ([]Provider, error) {
	var providers []Provider
	if len(c.TFE) > 0 {
		objs, err := NewTFE(c)
		if err != nil {
			return []Provider{}, err
		}
		if len(objs) > 0 {
			log.Info("Using Terraform Enterprise as state/locks provider")
			for _, tfeObj := range objs {
				providers = append(providers, tfeObj)
			}
		}
	}

	if len(c.GCP) > 0 {
		objs, err := NewGCP(c)
		if err != nil {
			return []Provider{}, err
		}
		if len(objs) > 0 {
			log.Info("Using Google Cloud as state/locks provider")
			for _, gcpObj := range objs {
				providers = append(providers, gcpObj)
			}
		}
	}

	if len(c.Gitlab) > 0 {
		objs := NewGitlab(c)
		if len(objs) > 0 {
			log.Info("Using Gitab as state/locks provider")
			for _, glObj := range objs {
				providers = append(providers, glObj)
			}
		}
	}

	if len(c.AWS) > 0 {
		objs := NewAWS(c)
		if len(objs) > 0 {
			log.Info("Using AWS (S3+DynamoDB) as state/locks provider")
			for _, awsObj := range objs {
				providers = append(providers, awsObj)
			}
		}
	}

	return providers, nil
}
