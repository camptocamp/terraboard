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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var svc *s3.S3
var bucket string

func init() {
	sess := session.Must(session.NewSession())
	svc = s3.New(sess, &aws.Config{})
	bucket = os.Getenv("AWS_BUCKET")
}

func stats(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Stats!")
}

func project(w http.ResponseWriter, r *http.Request) {
	proj := strings.TrimPrefix(r.URL.Path, "/project")
	result, err := svc.GetObjectWithContext(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(proj),
	})
	if err != nil {
		errObj := make(map[string]string)
		errObj["error"] = fmt.Sprintf("State not found: %v", proj)
		errObj["details"] = fmt.Sprintf("%v", err)
		j, _ := json.Marshal(errObj)
		io.WriteString(w, string(j))
		return
	}
	defer result.Body.Close()

	content, _ := ioutil.ReadAll(result.Body)
	io.WriteString(w, string(content))
}

func main() {
	http.HandleFunc("/", stats)
	http.HandleFunc("/project/", project)
	http.ListenAndServe(":80", nil)
}
