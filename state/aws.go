package state

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	aws_sdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/camptocamp/terraboard/config"
	"github.com/hashicorp/terraform/states/statefile"
	log "github.com/sirupsen/logrus"
)

// AWS is a state provider type, leveraging S3 and DynamoDB
type AWS struct {
	svc           *s3.S3
	dynamoSvc     *dynamodb.DynamoDB
	bucket        string
	dynamoTable   string
	keyPrefix     string
	fileExtension string
}

// NewAWS creates an AWS object
func NewAWS(c *config.Config) AWS {
	sess := session.Must(session.NewSession())

	awsConfig := aws_sdk.NewConfig()

	if len(c.AWS.APPRoleArn) > 0 {
		log.Debugf("Using %s role", c.AWS.APPRoleArn)
		creds := stscreds.NewCredentials(sess, c.AWS.APPRoleArn)
		awsConfig.WithCredentials(creds)
	}

	if e := c.AWS.Endpoint; e != "" {
		awsConfig.WithEndpoint(e)
	}

	return AWS{
		svc:           s3.New(sess, awsConfig),
		bucket:        c.AWS.S3.Bucket,
		keyPrefix:     c.AWS.S3.KeyPrefix,
		fileExtension: c.AWS.S3.FileExtension,
		dynamoSvc:     dynamodb.New(sess, awsConfig),
		dynamoTable:   c.AWS.DynamoDBTable,
	}
}

// GetLocks returns a map of locks by State path
func (a *AWS) GetLocks() (locks map[string]LockInfo, err error) {
	if a.dynamoTable == "" {
		err = fmt.Errorf("no dynamoDB table provided, not getting locks")
		return
	}

	results, err := a.dynamoSvc.Scan(&dynamodb.ScanInput{
		TableName: &a.dynamoTable,
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
	infoPrefix := fmt.Sprintf("%s/", a.bucket)
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

// GetStates returns a slice of State files in the S3 bucket
func (a *AWS) GetStates() (states []string, err error) {
	log.WithFields(log.Fields{
		"bucket": a.bucket,
		"prefix": a.keyPrefix,
	}).Debug("Listing states from S3")
	result, err := a.svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws_sdk.String(a.bucket),
		Prefix: &a.keyPrefix,
	})
	if err != nil {
		return states, err
	}

	var keys []string
	for _, obj := range result.Contents {
		for _, ext := range strings.Split(a.fileExtension, ",") {
			if strings.HasSuffix(*obj.Key, ext) {
				keys = append(keys, *obj.Key)
			}
		}
	}
	states = keys
	log.WithFields(log.Fields{
		"bucket": a.bucket,
		"prefix": a.keyPrefix,
		"states": len(states),
	}).Debug("Found states from S3")
	return states, nil
}

// GetState retrieves a single State from the S3 bucket
func (a *AWS) GetState(st, versionID string) (sf *statefile.File, err error) {
	log.WithFields(log.Fields{
		"path":       st,
		"version_id": versionID,
	}).Info("Retrieving state from S3")
	input := &s3.GetObjectInput{
		Bucket: aws_sdk.String(a.bucket),
		Key:    aws_sdk.String(st),
	}
	if versionID != "" {
		input.VersionId = &versionID
	}
	result, err := a.svc.GetObjectWithContext(context.Background(), input)
	if err != nil {
		log.WithFields(log.Fields{
			"path":       st,
			"version_id": versionID,
			"error":      err,
		}).Error("Error retrieving state from S3")
		errObj := make(map[string]string)
		errObj["error"] = fmt.Sprintf("State file not found: %v", st)
		errObj["details"] = fmt.Sprintf("%v", err)
		j, _ := json.Marshal(errObj)
		return sf, fmt.Errorf("%s", string(j))
	}
	defer result.Body.Close()

	sf, err = statefile.Read(result.Body)

	if sf == nil {
		return sf, fmt.Errorf("Failed to find state")
	}

	return
}

// GetVersions returns a slice of Version objects
func (a *AWS) GetVersions(state string) (versions []Version, err error) {
	versions = []Version{}
	result, err := a.svc.ListObjectVersions(&s3.ListObjectVersionsInput{
		Bucket: aws_sdk.String(a.bucket),
		Prefix: aws_sdk.String(state),
	})
	if err != nil {
		return
	}

	for _, v := range result.Versions {
		versions = append(versions, Version{
			ID:           *v.VersionId,
			LastModified: *v.LastModified,
		})
	}

	return
}
