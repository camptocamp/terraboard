package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/camptocamp/terraboard/config"
	"github.com/hashicorp/terraform/terraform"
)

var svc *s3.S3
var dynamoSvc *dynamodb.DynamoDB
var bucket string
var dynamoTable string
var keyPrefix string

func Setup(c *config.Config) {
	sess := session.Must(session.NewSession())
	svc = s3.New(sess, &aws.Config{})
	bucket = c.S3.Bucket
    keyPrefix = c.S3.KeyPrefix

	dynamoSvc = dynamodb.New(sess, &aws.Config{})
	dynamoTable = c.S3.DynamoDBTable
}

type LockInfo struct {
	ID        string
	Operation string
	Info      string
	Who       string
	Version   string
	Created   *time.Time
	Path      string
}

type Lock struct {
	LockID string
	Info   string
}

func GetLocks() (locks map[string]LockInfo, err error) {
	if dynamoTable == "" {
		err = fmt.Errorf("No dynamoDB table provided. Not getting locks.")
		return
	}

	results, err := dynamoSvc.Scan(&dynamodb.ScanInput{
		TableName: &dynamoTable,
	})
	if err != nil {
		return locks, err
	}

	var lockList []Lock
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &lockList)
	if err != nil {
		return locks, err
	}

	locks = make(map[string]LockInfo)
	infoPrefix := fmt.Sprintf("%s/", bucket)
	for _, lock := range lockList {
		if lock.Info != "" {
			var info LockInfo
			err = json.Unmarshal([]byte(lock.Info), &info)
			if err != nil {
				return locks, err
			}

			locks[strings.TrimPrefix(info.Path, infoPrefix)] = info
		}
	}
	return
}

func GetStates() (states []string, err error) {
	result, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
        Prefix: &keyPrefix,
	})
	if err != nil {
		return states, err
	}

	var keys []string
	for _, obj := range result.Contents {
		if strings.HasSuffix(*obj.Key, ".tfstate") {
			keys = append(keys, *obj.Key)
		}
	}
	states = keys
	return states, nil
}

func GetVersions(prefix string) (versions []*s3.ObjectVersion, err error) {
	result, err := svc.ListObjectVersions(&s3.ListObjectVersionsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return versions, err
	}

	return result.Versions, nil
}

func GetState(st, versionId string) (state *terraform.State, err error) {
	log.WithFields(log.Fields{
		"path":       st,
		"version_id": versionId,
	}).Info("Retrieving state from S3")
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(st),
	}
	if versionId != "" {
		input.VersionId = &versionId
	}
	result, err := svc.GetObjectWithContext(context.Background(), input)
	if err != nil {
		log.WithFields(log.Fields{
			"path":       st,
			"version_id": versionId,
			"error":      err,
		}).Error("Error retrieving state from S3")
		errObj := make(map[string]string)
		errObj["error"] = fmt.Sprintf("State file not found: %v", st)
		errObj["details"] = fmt.Sprintf("%v", err)
		j, _ := json.Marshal(errObj)
		return state, fmt.Errorf("%s", string(j))
	}
	defer result.Body.Close()

	content, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"path":       st,
			"version_id": versionId,
			"error":      err,
		}).Error("Error reading state from S3")
		errObj := make(map[string]string)
		errObj["error"] = fmt.Sprintf("Failed to read S3 response: %v", st)
		errObj["details"] = fmt.Sprintf("%v", err)
		j, _ := json.Marshal(errObj)
		return state, fmt.Errorf("%s", string(j))
	}

	json.Unmarshal(content, &state)

	return
}
