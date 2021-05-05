package state

import (
	"context"
	"fmt"
	"time"

	"github.com/camptocamp/terraboard/config"
	"github.com/hashicorp/terraform/states/statefile"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Minio is a state provider type
type Minio struct {
	minioClient *minio.Client
	ctx         *context.Context
	bucket      string
}

// NewMinio creates an Minio object
func NewMinio(c *config.Config) (*Minio, error) {

	endpoint := c.Minio.EndpointURL
	accessKeyID := c.Minio.AccessKey
	secretAccessKey := c.Minio.SecretKey
	disableSSL := c.Minio.DisableSSL
	bucket := c.Minio.BucketName

	useSSL := true
	if disableSSL {
		useSSL = false
	}

	ctx := context.Background()

	options := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	}

	client, _ := minio.New(endpoint, options)

	return &Minio{
		minioClient: client,
		ctx:         &ctx,
		bucket:      bucket,
	}, nil
}

// GetLocks returns a map of locks by State path
func (a *Minio) GetLocks() (locks map[string]LockInfo, err error) {
	locks = make(map[string]LockInfo)

	// do not read any lock states

	return
}

// GetStates returns a slice of State files in the minio bucket
func (a *Minio) GetStates() (states []string, err error) {

	var stateFiles []string

	result := a.minioClient.ListObjects(*a.ctx, a.bucket, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})
	if err != nil {
		return states, err
	}

	var keys []string
	for obj := range result {
		keys = append(keys, obj.Key)
	}
	stateFiles = keys

	return stateFiles, nil
}

// GetState retrieves a single State from the minio bucket
func (a *Minio) GetState(path, version string) (sf *statefile.File, err error) {

	ctx := context.Background()

	state, err := a.minioClient.GetObject(ctx, a.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return
	}

	// Parse the statefile
	sf, err = statefile.Read(state)
	if sf == nil {
		return nil, fmt.Errorf("Unable to parse the statefile for workspace %s version %s", path, version)
	}

	return
}

// GetVersions returns a slice of Version objects
func (a *Minio) GetVersions(_ string) (versions []Version, err error) {

	// currently only unversioned buckets supported
	versions = []Version{}
	versions = append(versions, Version{
		ID:           "1",
		LastModified: time.Now(),
	})

	return
}
