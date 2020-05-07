package state

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/camptocamp/terraboard/config"
	"github.com/hashicorp/terraform/states/statefile"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GCP is a state provider type, leveraging S3 and DynamoDB
type GCP struct {
	svc    *storage.Client
	bucket string
}

// NewGCP creates an GCP object
func NewGCP(c *config.Config) (GCP, error) {
	ctx := context.Background()
	var client *storage.Client
	var err error
	if c.GCP.GCPSAKey != "" {
		log.WithFields(log.Fields{
			"path": c.GCP.GCPSAKey,
		}).Info("Authenticating using service account key")
		opt := option.WithCredentialsFile(c.GCP.GCPSAKey)
		client, err = storage.NewClient(ctx, opt) // Use service account key
	} else {
		client, err = storage.NewClient(ctx) // Use base credentials
	}

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return GCP{}, err
	}

	return GCP{
		svc:    client,
		bucket: c.GCP.GCSBucket,
	}, nil
}

// GetLocks returns a map of locks by State path
func (a *GCP) GetLocks() (locks map[string]LockInfo, err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	var lockFiles []string
	it := a.svc.Bucket(a.bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(attrs.Name, ".tflock") {
			lockFiles = append(lockFiles, attrs.Name)
		}
	}

	locks = make(map[string]LockInfo)
	for _, lockFile := range lockFiles {
		ctx, cancel := context.WithTimeout(ctx, time.Second*50)
		defer cancel()
		rc, err := a.svc.Bucket(a.bucket).Object(lockFile).NewReader(ctx)
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		data, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil, err
		}

		var info LockInfo
		err = json.Unmarshal([]byte(data), &info)
		if err != nil {
			return nil, err
		}

		locks[lockFile] = info
	}

	return locks, nil
}

// GetStates returns a slice of State files in the GCS bucket
func (a *GCP) GetStates() (states []string, err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	var stateFiles []string
	it := a.svc.Bucket(a.bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(attrs.Name, ".tfstate") {
			stateFiles = append(stateFiles, attrs.Name)
		}
	}

	return stateFiles, nil
}

// GetState retrieves a single State from the S3 bucket
func (a *GCP) GetState(st, versionID string) (sf *statefile.File, err error) {
	log.WithFields(log.Fields{
		"path":       st,
		"version_id": versionID,
	}).Info("Retrieving state from GCS")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	obj := a.svc.Bucket(a.bucket).Object(st)
	if versionID != "" {
		version, err := strconv.ParseInt(versionID, 10, 64)
		if err != nil {
			return nil, err
		}
		obj = obj.Generation(version)
	}
	rc, err := obj.NewReader(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"path":       st,
			"version_id": versionID,
			"error":      err,
		}).Error("Error retrieving state from GCS")
		errObj := make(map[string]string)
		errObj["error"] = fmt.Sprintf("State file not found: %v", st)
		errObj["details"] = fmt.Sprintf("%v", err)
		j, _ := json.Marshal(errObj)
		return sf, fmt.Errorf("%s", string(j))
	}
	defer rc.Close()

	sf, err = statefile.Read(rc)

	if sf == nil {
		return sf, fmt.Errorf("Failed to find state")
	}

	return
}

// GetVersions returns a slice of Version objects
func (a *GCP) GetVersions(state string) (versions []Version, err error) {
	versions = []Version{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	q := storage.Query{
		Versions: true,
		Prefix:   state,
	}

	it := a.svc.Bucket(a.bucket).Objects(ctx, &q)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		if attrs.Name == state {
			tm := attrs.Updated
			versions = append(versions, Version{
				ID:           strconv.FormatInt(attrs.Generation, 10),
				LastModified: tm,
			})
		}
	}

	return
}
