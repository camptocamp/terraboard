package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/camptocamp/terraboard/util"
)

var svc *s3.S3
var bucket string

func init() {
	sess := session.Must(session.NewSession())
	svc = s3.New(sess, &aws.Config{})
	bucket = os.Getenv("AWS_BUCKET")
}

func States(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		errObj := make(map[string]string)
		errObj["error"] = "Failed to list states"
		errObj["details"] = fmt.Sprintf("%v", err)
		j, _ := json.Marshal(errObj)
		io.WriteString(w, string(j))
		return
	}

	var keys []string

	for _, obj := range result.Contents {
		if strings.HasSuffix(*obj.Key, ".tfstate") {
			keys = append(keys, *obj.Key)
		}
	}

	j, _ := json.Marshal(keys)
	io.WriteString(w, string(j))
}

func State(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := util.TrimBase(r, "api/state")
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(st),
	}
	if versionId := r.URL.Query().Get("versionid"); versionId != "" {
		input.VersionId = &versionId
	}
	result, err := svc.GetObjectWithContext(context.Background(), input)
	if err != nil {
		errObj := make(map[string]string)
		errObj["error"] = fmt.Sprintf("State file not found: %v", st)
		errObj["details"] = fmt.Sprintf("%v", err)
		j, _ := json.Marshal(errObj)
		io.WriteString(w, string(j))
		return
	}
	defer result.Body.Close()

	content, _ := ioutil.ReadAll(result.Body)
	io.WriteString(w, string(content))
}

func History(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := util.TrimBase(r, "api/history/")
	result, err := svc.ListObjectVersions(&s3.ListObjectVersionsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(st),
	})
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

	j, err := json.Marshal(result.Versions)
	if err != nil {
		log.Errorf("Failed to marshal json: %v", err)
	}
	io.WriteString(w, string(j))
}
