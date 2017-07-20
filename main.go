package main

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
)

var svc *s3.S3
var bucket string
var baseUrl string

func init() {
	sess := session.Must(session.NewSession())
	svc = s3.New(sess, &aws.Config{})
	bucket = os.Getenv("AWS_BUCKET")
	baseUrl = os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "/"
	}
}

func addBase(path string) string {
	return fmt.Sprintf("%s%s", baseUrl, path)
}

func trimBase(r *http.Request, prefix string) string {
	return strings.TrimPrefix(r.URL.Path, addBase(prefix))
}

func idx(w http.ResponseWriter, r *http.Request) {
	idx, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Errorf("Failed to open index.html: %v", err)
		// TODO: Return error page
	}
	idxStr := string(idx)
	idxStr = strings.Replace(idxStr, "base href=\"/\"", fmt.Sprintf("base href=\"%s\"", baseUrl), 1)
	io.WriteString(w, idxStr)
}

func states(w http.ResponseWriter, r *http.Request) {
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

func state(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := trimBase(r, "api/state")
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

func history(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st := trimBase(r, "api/history/")
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

func main() {
	http.HandleFunc(baseUrl, idx)
	staticFs := http.FileServer(http.Dir("static"))
	http.Handle(addBase("static/"), http.StripPrefix(addBase("static"), staticFs))
	http.HandleFunc(addBase("api/states"), states)
	http.HandleFunc(addBase("api/state/"), state)
	http.HandleFunc(addBase("api/history/"), history)
	http.ListenAndServe(":80", nil)
}
