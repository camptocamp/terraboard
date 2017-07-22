package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/camptocamp/terraboard/util"
)

type StateVersions struct {
	Versions map[string]string
}

var states map[string]*StateVersions

func init() {
	states = make(map[string]*StateVersions)
}

func GetState(w http.ResponseWriter, r *http.Request) (state string, err error) {
	st := util.TrimBase(r, "api/state")
	if _, ok := states[st]; !ok {
		// Init
		states[st] = &StateVersions{}
		states[st].Versions = make(map[string]string)
	}

	versionId := r.URL.Query().Get("versionid")
	if s, ok := states[st].Versions[versionId]; ok {
		// Return cached version
		log.Infof("Getting %s/%s from cache", st, versionId)
		return s, nil
	}

	// Retrieve from S3
	log.Infof("Retrieving %s/%s from S3", st, versionId)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(st),
	}
	if versionId != "" {
		input.VersionId = &versionId
	}
	result, err := svc.GetObjectWithContext(context.Background(), input)
	if err != nil {
		log.Errorf("Error retrieving %s/%s from S3: %v", st, versionId, err)
		errObj := make(map[string]string)
		errObj["error"] = fmt.Sprintf("State file not found: %v", st)
		errObj["details"] = fmt.Sprintf("%v", err)
		j, _ := json.Marshal(errObj)
		return state, fmt.Errorf("%s", string(j))
	}
	defer result.Body.Close()

	content, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Errorf("Error reading %s/%s from S3: %v", st, versionId, err)
		errObj := make(map[string]string)
		errObj["error"] = fmt.Sprintf("Failed to read S3 response: %v", st)
		errObj["details"] = fmt.Sprintf("%v", err)
		j, _ := json.Marshal(errObj)
		return state, fmt.Errorf("%s", string(j))
	}
	state = string(content)

	if versionId != "" {
		// Store in cache
		log.Infof("Adding %s/%s to cache", st, versionId)
		states[st].Versions[versionId] = state
	}

	return
}
