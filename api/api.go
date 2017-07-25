package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/camptocamp/terraboard/db"
	"github.com/camptocamp/terraboard/util"
)

var svc *s3.S3
var bucket string
var states []string

func init() {
	sess := session.Must(session.NewSession())
	svc = s3.New(sess, &aws.Config{})
	bucket = os.Getenv("AWS_BUCKET")

	db.Init()
	go refreshDB()
}

func refreshDB() {
	for {
		log.Infof("Refreshing DB from S3")
		err := refreshStates()
		if err != nil {
			log.Errorf("Failed to build cache: %s", err)
		}

		// Refresh knownVersions
		knownVersions := db.KnownVersions()

		for _, st := range states {
			// FIXME: UpdateState duplicates entries!
			// We should not use it, instead point to the right version in the UI
			// s3State, _ := GetS3State(st, "")
			//db.UpdateState(st, "", s3State)

			versions, _ := getVersions(st)
			for _, v := range versions {
				if isKnownVersion(knownVersions, *v.VersionId) {
					log.Infof("Version %s for %s is already known, skipping", *v.VersionId, st)
					continue
				}
				state, _ := GetS3State(st, *v.VersionId)
				db.InsertState(st, *v.VersionId, state)
				if err != nil {
					log.Errorf("Failed to insert state %s/%s: %v", st, *v.VersionId, err)
				}
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func isKnownVersion(knownVersions []string, version string) bool {
	for _, v := range knownVersions {
		if v == version {
			return true
		}
	}
	return false
}

func refreshStates() error {
	result, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	var keys []string
	for _, obj := range result.Contents {
		if strings.HasSuffix(*obj.Key, ".tfstate") {
			keys = append(keys, *obj.Key)
		}
	}
	states = keys
	return nil
}

func ApiStates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err := refreshStates()
	if err != nil {
		errObj := make(map[string]string)
		errObj["error"] = "Failed to list states"
		errObj["details"] = fmt.Sprintf("%v", err)
		j, _ := json.Marshal(errObj)
		io.WriteString(w, string(j))
		return
	}

	j, _ := json.Marshal(states)
	io.WriteString(w, string(j))
}

func ApiState(w http.ResponseWriter, r *http.Request) {
	st := util.TrimBase(r, "api/state/")
	versionId := r.URL.Query().Get("versionid")
	if versionId == "" {
		versionId, _ = defaultVersion(st) // TODO: err
	}
	state := db.GetState(st, versionId)

	jState, _ := json.Marshal(state)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, string(jState))
}

func defaultVersion(path string) (string, error) {
	versions, err := getVersions(path) // TODO: err
	if err != nil {
		return "", fmt.Errorf("Failed to list versions for %s: %v", path, err)
	}

	for _, v := range versions {
		if *v.IsLatest {
			return *v.VersionId, nil
		}
	}

	return "", fmt.Errorf("Failed to find default version for %s", path)
}

func getVersions(prefix string) (versions []*s3.ObjectVersion, err error) {
	result, err := svc.ListObjectVersions(&s3.ListObjectVersionsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return versions, err
	}
	return result.Versions, nil
}

func ApiHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := util.TrimBase(r, "api/history/")
	result, err := getVersions(st)
	if err != nil {
		errObj := make(map[string]string)
		errObj["error"] = fmt.Sprintf("State file history not found: %v", st)
		errObj["details"] = fmt.Sprintf("%v", err)
		j, err := json.Marshal(errObj)
		if err != nil {
			log.Errorf("Failed to marshal json: %v", err)
		}
		io.WriteString(w, string(j))
		return
	}

	j, err := json.Marshal(result)
	if err != nil {
		log.Errorf("Failed to marshal json: %v", err)
	}
	io.WriteString(w, string(j))
}

func ApiSearchResource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.URL.Query()
	result := db.SearchResource(query)
	j, err := json.Marshal(result)
	if err != nil {
		log.Errorf("Failed to marshal json: %v", err)
	}
	io.WriteString(w, string(j))
}

func ApiSearchAttribute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.URL.Query()
	defaultV, _ := defaultVersion("")
	result := db.SearchAttribute(query, defaultV)
	j, err := json.Marshal(result)
	if err != nil {
		log.Errorf("Failed to marshal json: %v", err)
	}
	io.WriteString(w, string(j))
}
