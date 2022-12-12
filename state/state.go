package state

import (
	"time"

	"github.com/camptocamp/terraboard/config"
	"github.com/camptocamp/terraboard/internal/terraform/states/statefile"
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
	var err error
	var providers []Provider

	providers, err = handleTFE(c, providers)
	if err != nil {
		return []Provider{}, err
	}

	providers, err = handleGCP(c, providers)
	if err != nil {
		return []Provider{}, err
	}

	providers = handleGitLab(c, providers)

	providers = handleAWS(c, providers)

	providers, err = handleCOS(c, providers)
	if err != nil {
		return []Provider{}, err
	}

	return providers, nil
}

func handleTFE(c *config.Config, providers []Provider) ([]Provider, error) {
	if len(c.TFE) > 0 {
		objs, err := NewTFECollection(c)
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
	return providers, nil
}

func handleGCP(c *config.Config, providers []Provider) ([]Provider, error) {
	if len(c.GCP) > 0 {
		objs, err := NewGCPCollection(c)
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
	return providers, nil
}

func handleGitLab(c *config.Config, providers []Provider) []Provider {
	if len(c.Gitlab) > 0 {
		objs := NewGitlabCollection(c)
		if len(objs) > 0 {
			log.Info("Using Gitab as state/locks provider")
			for _, glObj := range objs {
				providers = append(providers, glObj)
			}
		}
	}
	return providers
}

func handleAWS(c *config.Config, providers []Provider) []Provider {
	if len(c.AWS) > 0 {
		objs := NewAWSCollection(c)
		if len(objs) > 0 {
			log.Info("Using AWS (S3+DynamoDB) as state/locks provider")
			for _, awsObj := range objs {
				providers = append(providers, awsObj)
			}
		}
	}
	return providers
}

func handleCOS(c *config.Config, providers []Provider) ([]Provider, error) {
	if len(c.COS) > 0 {
		objs, err := NewCOSCollection(c)
		if err != nil {
			return nil, err
		}
		if len(objs) > 0 {
			log.Info("Using Tencent Cloud Object Storage as state/locks provider")
			for _, cosObj := range objs {
				providers = append(providers, cosObj)
			}
		}
	}
	return providers, nil
}
