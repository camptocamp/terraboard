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

// GCP is a state provider type, leveraging GCS
type GCP struct {
	svc     *storage.Client
	buckets []string
}

// NewGCP creates an GCP object
func NewGCP(c *config.Config) ([]*GCP, error) {
	ctx := context.Background()

	var client *storage.Client
	var gcpInstances []*GCP
	var err error
	for _, gcp := range c.GCP {
		if gcp.GCPSAKey != "" {
			log.WithFields(log.Fields{
				"path": gcp.GCPSAKey,
			}).Info("Authenticating using service account key")
			opt := option.WithCredentialsFile(gcp.GCPSAKey)
			client, err = storage.NewClient(ctx, opt) // Use service account key
		} else {
			client, err = storage.NewClient(ctx) // Use base credentials
		}

		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
			return nil, err
		}

		instance := &GCP{
			svc:     client,
			buckets: gcp.GCSBuckets,
		}
		gcpInstances = append(gcpInstances, instance)

		log.WithFields(log.Fields{
			"buckets": gcp.GCSBuckets,
		}).Info("Client successfully created")
	}

	return gcpInstances, nil
}

// GetLocks returns a map of locks by State path
func (a *GCP) GetLocks() (locks map[string]LockInfo, err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	var lockFiles []string
	for _, bucketName := range a.buckets {
		it := a.svc.Bucket(bucketName).Objects(ctx, nil)
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
			rc, err := a.svc.Bucket(bucketName).Object(lockFile).NewReader(ctx)
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

			locks[strings.Join([]string{bucketName, lockFile}, "/")] = info
		}

	}

	return locks, nil
}

// GetStates returns a slice of State files in the GCS bucket
func (a *GCP) GetStates() (states []string, err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	var stateFiles []string
	for _, bucketName := range a.buckets {
		it := a.svc.Bucket(bucketName).Objects(ctx, nil)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}

			if strings.HasSuffix(attrs.Name, ".tfstate") {
				stateFiles = append(stateFiles, strings.Join([]string{bucketName, attrs.Name}, "/"))
			}
		}
	}

	return stateFiles, nil
}

// GetState retrieves a single State from the GCS bucket
func (a *GCP) GetState(st, versionID string) (sf *statefile.File, err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	bucketSplit := strings.Index(st, "/")
	bucketName := st[0:bucketSplit]
	fileName := st[bucketSplit+1:]

	obj := a.svc.Bucket(bucketName).Object(fileName)
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

	log.WithFields(log.Fields{
		"path":       st,
		"version_id": versionID,
	}).Info("State read from GCS")

	return
}

// GetVersions returns a slice of Version objects
func (a *GCP) GetVersions(state string) (versions []Version, err error) {
	versions = []Version{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	bucketSplit := strings.Index(state, "/")
	bucketName := state[0:bucketSplit]
	fileName := state[bucketSplit+1:]

	q := storage.Query{
		Versions: true,
		Prefix:   fileName,
	}

	it := a.svc.Bucket(bucketName).Objects(ctx, &q)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		if attrs.Name == fileName {
			tm := attrs.Updated
			versions = append(versions, Version{
				ID:           strconv.FormatInt(attrs.Generation, 10),
				LastModified: tm,
			})
		}
	}

	return
}
