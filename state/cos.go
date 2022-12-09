package state

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/camptocamp/terraboard/config"
	"github.com/camptocamp/terraboard/internal/terraform/states/statefile"
	log "github.com/sirupsen/logrus"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

// COS is a state provider type, leveraging Tencent Cloud COS
type COS struct {
	svc    *cos.Client
	bucket string // one svc one bucket
	Ext    cosExts
}

type cosExts struct {
	keyPrefix    string
	noLocks      bool
	noVersioning bool
}

type CosExt interface {
	apply(*cosExts)
}

type funcCosExt struct {
	f func(*cosExts)
}

func (fdo *funcCosExt) apply(do *cosExts) {
	fdo.f(do)
}

func newFuncCosExt(f func(*cosExts)) *funcCosExt {
	return &funcCosExt{f: f}
}

func WithKeyPrefix(kp string) CosExt {
	return newFuncCosExt(func(exts *cosExts) {
		exts.keyPrefix = kp
	})
}

func WithNoLocks(noLocks bool) CosExt {
	return newFuncCosExt(func(exts *cosExts) {
		exts.noLocks = noLocks
	})
}

func WithNoVersioning(noVersioning bool) CosExt {
	return newFuncCosExt(func(exts *cosExts) {
		exts.noVersioning = noVersioning
	})
}

var defaultExt = cosExts{
	keyPrefix:    "",
	noLocks:      false,
	noVersioning: false,
}

// NewCOS creates an COS object
func NewCOS(cosCfg config.COSConfig, exts ...CosExt) (cosInstance *COS, err error) {
	if len(cosCfg.Bucket) == 0 {
		return nil, nil
	}

	log.WithFields(log.Fields{
		"SecretId":  cosCfg.SecretId,
		"SecretKey": cosCfg.SecretKey,
		"Bucket":    cosCfg.Bucket,
		"Region":    cosCfg.Region,
		"KeyPrefix": cosCfg.KeyPrefix,
	}).Info("Begin to NewCOS:")

	client, err := UseTencentCosClient(&cosCfg)
	if err != nil {
		return
	}

	cosInstance = &COS{
		svc:    client,
		bucket: cosCfg.Bucket,
		// default extension
		Ext:    defaultExt,
	}

	// the specified extension
	for _, ext := range exts {
		ext.apply(&cosInstance.Ext)
	}

	log.WithFields(log.Fields{
		"bucket": cosCfg.Bucket,
		"exts":   cosInstance.Ext,
	}).Info("Client successfully created")

	return
}

func UseTencentCosClient(cosCfg *config.COSConfig) (client *cos.Client, err error) {
	u, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cosCfg.Bucket, cosCfg.Region))
	if err != nil {
		return
	}

	baseUrl := &cos.BaseURL{
		BucketURL: u,
	}

	client = cos.NewClient(baseUrl, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:     cosCfg.SecretId,
			SecretKey:    cosCfg.SecretKey,
			SessionToken: cosCfg.SecretToken,
		},
	})

	return
}

// NewCOSCollection instantiate all needed COS objects configurated by the user and return a slice
func NewCOSCollection(cfg *config.Config) ([]*COS, error) {
	var cosInstances []*COS
	for _, cos := range cfg.COS {
		var exts []CosExt
		if len(cos.KeyPrefix) > 0 {
			exts = append(exts, WithKeyPrefix(cos.KeyPrefix))
		}
		if cfg.Provider.NoLocks {
			exts = append(exts, WithNoLocks(cfg.Provider.NoLocks))
		}
		if cfg.Provider.NoVersioning {
			exts = append(exts, WithNoVersioning(cfg.Provider.NoVersioning))
		}

		cosIns, err := NewCOS(cos, exts...)
		if err != nil || cosIns == nil {
			return nil, err
		}

		cosInstances = append(cosInstances, cosIns)
	}

	return cosInstances, nil
}

// GetLocks returns a map of locks by State path
func (a *COS) GetLocks() (locks map[string]LockInfo, err error) {
	if a.Ext.noLocks {
		locks = make(map[string]LockInfo)
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	var lockFiles []string
	opt := &cos.BucketGetOptions{
		Prefix:  a.Ext.keyPrefix,
		MaxKeys: 100,
	}

	ret, _, err := a.svc.Bucket.Get(context.Background(), opt)
	if err != nil {
		log.WithFields(log.Fields{
			"bucket": a.bucket,
		}).Error("Error retrieving lockfile key from COS bucket!")
		return nil, err
	}

	for _, c := range ret.Contents {
		if strings.HasSuffix(c.Key, ".tflock") {
			lockFiles = append(lockFiles, c.Key)
			log.WithFields(log.Fields{
				"key":  c.Key,
				"size": c.Size,
			}).Debug("Got one lockfile key from COS.")
		}
	}

	locks = make(map[string]LockInfo)
	for _, lockFile := range lockFiles {
		ctx, cancel := context.WithTimeout(ctx, time.Second*50)
		defer cancel()

		resp, err := a.svc.Object.Get(ctx, lockFile, nil)
		if err != nil {
			if err != nil {
				log.WithFields(log.Fields{
					"key": lockFile,
				}).Error("Error retrieving lockfile from COS!")
				return nil, err
			}
		}
		defer resp.Body.Close()

		log.WithFields(log.Fields{
			"key":  lockFile,
			"body": resp.Body,
		}).Debug("Got one lockfile from COS.")

		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var info LockInfo
		err = json.Unmarshal(bs, &info)
		if err != nil {
			return nil, err
		}

		// key:lockFileName
		locks[lockFile] = info
	}
	return locks, nil
}

// GetStates returns a slice of State files in the COS bucket
func (a *COS) GetStates() (states []string, err error) {
	ctx := context.Background()
	_, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	var stateFiles []string
	truncatedListing := true
	opt := &cos.BucketGetOptions{
		Prefix:  a.Ext.keyPrefix,
		MaxKeys: 100,
	}

	for truncatedListing {
		ret, _, err := a.svc.Bucket.Get(context.Background(), opt)

		log.WithFields(log.Fields{
			"bucket":  ret.Name,
			"content": ret.Contents,
			"ALL":     ret,
		}).Info("begin to GetStates,")

		if err != nil {
			log.WithFields(log.Fields{
				"bucket": a.bucket,
			}).Error("Error retrieving statefiles from COS bucket!")
			return nil, err
		}

		for _, c := range ret.Contents {
			if strings.HasSuffix(c.Key, ".tfstate") {
				// item:stateFileName
				stateFiles = append(stateFiles, c.Key)
				log.WithFields(log.Fields{
					"bucket": a.bucket,
					"key":    c.Key,
				}).Debug("Got one statefile from COS.")
			}
		}
		truncatedListing = ret.IsTruncated
	}

	log.WithFields(log.Fields{
		"count": len(stateFiles),
	}).Info("Found statefiles from COS.")
	return stateFiles, nil
}

// GetState retrieves a single State from the COS bucket
func (a *COS) GetState(st, versionID string) (sf *statefile.File, err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	var verId string
	if versionID != "" && !a.Ext.noVersioning {
		log.WithField("versionID", versionID).Debug("Set the versionID in GetState.")
		verId = versionID
	}

	ret, err := a.svc.Object.Get(ctx, st, nil, verId)
	if err != nil {
		log.WithFields(log.Fields{
			"path":       st,
			"version_id": verId,
			"error":      err,
		}).Error("Error retrieving state from COS")

		errObj := make(map[string]string)
		errObj["error"] = fmt.Sprintf("State file not found: %v", st)
		errObj["details"] = fmt.Sprintf("%v", err)
		j, _ := json.Marshal(errObj)
		return sf, fmt.Errorf("%s", string(j))
	}
	defer ret.Body.Close()

	sf, err = statefile.Read(ret.Body)
	if sf == nil || err != nil {
		return sf, fmt.Errorf("failed to find state: %v", err)
	}

	log.WithFields(log.Fields{
		"path":       st,
		"version_id": versionID,
	}).Info("Read state from COS.")
	return
}

// GetVersions returns a slice of Version objects from COS bucket
func (a *COS) GetVersions(state string) (versions []Version, err error) {
	log.WithFields(log.Fields{
		"noVersioning": a.Ext.noVersioning,
		"noLocks":      a.Ext.noLocks,
		"state":        state,
	}).Info("Begin to GetVersions,")

	versions = []Version{}
	if a.Ext.noVersioning {
		versions = append(versions, Version{
			ID:           state,
			LastModified: time.Now(),
		})
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	opt := &cos.BucketGetObjectVersionsOptions{
		Prefix: state,
	}

	ret, _, err := a.svc.Bucket.GetObjectVersions(ctx, opt)
	if err != nil {
		return
	}

	for _, v := range ret.Version {
		modified, _ := time.Parse(time.RFC3339, v.LastModified)
		versions = append(versions, Version{
			ID:           v.VersionId,
			LastModified: modified,
		})
	}

	log.WithFields(log.Fields{
		"key":   state,
		"count": len(versions),
	}).Info("Read versions from COS.")
	return
}
