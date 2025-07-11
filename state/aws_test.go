package state

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	dynamock "github.com/gusaul/go-dynamock"

	"github.com/camptocamp/terraboard/config"
	"github.com/sirupsen/logrus"
)

// TODO: tests for the AWS features of the state package

func TestNewAWS(t *testing.T) {
	awsInstance := NewAWS(
		config.AWSConfig{
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Region:          "us-east-1",
			Endpoint:        "http://localhost:8000",
		},
		config.S3BucketConfig{
			Bucket: "test",
		},
		false,
		false,
	)

	if awsInstance == nil || awsInstance.svc == nil {
		t.Error("AWS instance is nil")
	}
}

func TestNewAWSNoBucket(t *testing.T) {
	awsInstance := NewAWS(
		config.AWSConfig{
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Region:          "us-east-1",
			Endpoint:        "http://localhost:8000",
		},
		config.S3BucketConfig{
			Bucket: "",
		},
		false,
		false,
	)

	if awsInstance != nil {
		t.Error("AWS instance should be nil")
	}
}

func TestNewAWSNoCredentials(t *testing.T) {
	// Remove os.Exit call to logrus Fatal
	logrus.StandardLogger().ExitFunc = func(code int) {}

	// Redirect logrus output to a buffer
	buf := &bytes.Buffer{}
	logrus.SetOutput(buf)

	_ = NewAWS(
		config.AWSConfig{
			AccessKey:       "",
			SecretAccessKey: "",
			Region:          "us-east-1",
			Endpoint:        "http://localhost:8000",
		},
		config.S3BucketConfig{
			Bucket: "test",
		},
		false,
		false,
	)

	// Test output
	t.Log(buf)
	if strings.Count(buf.String(), "\n") > 0 {
		t.Error("Expected no log output")
	}
}

func TestNewAWSWithAPPRoleArn(t *testing.T) {
	awsInstance := NewAWS(
		config.AWSConfig{
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Region:          "us-east-1",
			Endpoint:        "http://localhost:8000",
			APPRoleArn:      "arn:aws:iam::123456789012:role/app-role",
			ExternalID:      "123456789",
		},
		config.S3BucketConfig{
			Bucket: "test",
		},
		false,
		false,
	)

	if awsInstance == nil || awsInstance.svc == nil {
		t.Error("AWS instance is nil")
	}
}

func TestNewAWSCollection(t *testing.T) {
	config := config.Config{
		AWS: []config.AWSConfig{
			{
				AccessKey:       "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				Region:          "us-east-1",
				Endpoint:        "http://localhost:8000",
				S3: []config.S3BucketConfig{
					{
						Bucket: "test",
					},
					{
						Bucket: "test2",
					},
				},
			},
		},
		Version:        false,
		ConfigFilePath: "",
		Provider: config.ProviderConfig{
			NoVersioning: false,
			NoLocks:      false,
		},
		DB:     config.DBConfig{},
		TFE:    []config.TFEConfig{},
		GCP:    []config.GCPConfig{},
		Gitlab: []config.GitlabConfig{},
		Web:    config.WebConfig{},
	}
	instances := NewAWSCollection(&config)

	if instances == nil || len(instances) != 2 {
		t.Error("AWS instances are nil or not the expected number")
	}
}

func TestGetLocksEmpty(t *testing.T) {
	awsInstance := NewAWS(
		config.AWSConfig{
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Region:          "us-east-1",
			Endpoint:        "http://localhost:8000",
			DynamoDBTable:   "test-locks",
		},
		config.S3BucketConfig{
			Bucket: "test",
		},
		false,
		false,
	)
	dyna, mock := dynamock.New()
	awsInstance.dynamoSvc = dyna

	mock.ExpectScan().Table(awsInstance.dynamoTable).WillReturns(dynamodb.ScanOutput{})

	locks, err := awsInstance.GetLocks()
	if err != nil {
		t.Error(err)
	} else if len(locks) != 0 {
		t.Error("Expected no locks")
	}
}

func TestGetLocksNoLocks(t *testing.T) {
	awsInstance := NewAWS(
		config.AWSConfig{
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Region:          "us-east-1",
			Endpoint:        "http://localhost:8000",
			DynamoDBTable:   "test-locks",
		},
		config.S3BucketConfig{
			Bucket: "test",
		},
		true,
		false,
	)

	locks, _ := awsInstance.GetLocks()
	if len(locks) != 0 {
		t.Error("Locks should be empty due to noLocks option")
	}
}

func TestGetLocksNoDynamoTable(t *testing.T) {
	awsInstance := NewAWS(
		config.AWSConfig{
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Region:          "us-east-1",
			Endpoint:        "http://localhost:8000",
			DynamoDBTable:   "",
		},
		config.S3BucketConfig{
			Bucket: "test",
		},
		false,
		false,
	)

	_, err := awsInstance.GetLocks()
	if err == nil {
		t.Error("Err shouldn't be nil due to missing dynamodb table")
	}
}

func TestGetLocks(t *testing.T) {
	awsInstance := NewAWS(
		config.AWSConfig{
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Region:          "us-east-1",
			Endpoint:        "http://localhost:8000",
			DynamoDBTable:   "test-locks",
		},
		config.S3BucketConfig{
			Bucket: "test",
		},
		false,
		false,
	)
	dyna, mock := dynamock.New()
	awsInstance.dynamoSvc = dyna

	mock.ExpectScan().Table(awsInstance.dynamoTable).WillReturns(dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"LockID": {
					N: aws.String("lock1"),
				},
			},
			{
				"LockID": {
					N: aws.String("lock2"),
				},
				"Info": {
					S: aws.String(`{
						"Operation":"test",
						"Who":"testUser",
						"Version":"1.0.0",
						"Path":"test_path"
					 }`),
				},
			},
		},
	})

	locks, err := awsInstance.GetLocks()
	if err != nil {
		t.Error(err)
	} else if len(locks) != 1 {
		t.Error("Expected one lock")
	}
}

type s3Mock struct {
	s3iface.S3API
}

func (s *s3Mock) ListObjectsV2(_ *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return &s3.ListObjectsV2Output{Contents: []*s3.Object{
		{Key: aws.String("test.tfstate")}, {Key: aws.String("test2.tfstate")}, {Key: aws.String("test3.tfstate")}},
		IsTruncated: func() *bool { b := false; return &b }(),
		KeyCount:    func() *int64 { b := int64(3); return &b }(),
		MaxKeys:     func() *int64 { b := int64(1000); return &b }()}, nil
}
func (s *s3Mock) ListObjectVersions(_ *s3.ListObjectVersionsInput) (*s3.ListObjectVersionsOutput, error) {
	return &s3.ListObjectVersionsOutput{
		Versions: []*s3.ObjectVersion{
			{Key: aws.String("testId"), VersionId: aws.String("test"), LastModified: aws.Time(time.Now())},
			{Key: aws.String("testId2"), VersionId: aws.String("test2"), LastModified: aws.Time(time.Now())},
		},
	}, nil
}
func (s *s3Mock) GetObjectWithContext(_ aws.Context, _ *s3.GetObjectInput, _ ...request.Option) (*s3.GetObjectOutput, error) {
	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"Version": 4, "Serial": 3, "TerraformVersion": "0.12.0"}`))),
	}, nil
}

func TestGetStates(t *testing.T) {
	awsInstance := NewAWS(
		config.AWSConfig{
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Region:          "us-east-1",
			Endpoint:        "http://localhost:8000",
			DynamoDBTable:   "test-locks",
		},
		config.S3BucketConfig{
			Bucket:        "test",
			FileExtension: []string{".tfstate"},
		},
		false,
		false,
	)

	mock := s3Mock{}
	awsInstance.svc = &mock

	states, err := awsInstance.GetStates()
	if err != nil {
		t.Error(err)
	} else if len(states) != 3 {
		t.Error("Expected three states")
	}
}

func TestGetState(t *testing.T) {
	awsInstance := NewAWS(
		config.AWSConfig{
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Region:          "us-east-1",
			Endpoint:        "http://localhost:8000",
			DynamoDBTable:   "test-locks",
		},
		config.S3BucketConfig{
			Bucket:        "test",
			FileExtension: []string{".tfstate"},
		},
		false,
		false,
	)

	mock := s3Mock{}
	awsInstance.svc = &mock

	state, err := awsInstance.GetState("test", "ver_test")
	if err != nil {
		t.Error(err)
	} else if state == nil {
		t.Error("Unexpected nil state")
	}
}

func TestGetVersions(t *testing.T) {
	awsInstance := NewAWS(
		config.AWSConfig{
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Region:          "us-east-1",
			Endpoint:        "http://localhost:8000",
			DynamoDBTable:   "test-locks",
		},
		config.S3BucketConfig{
			Bucket:        "test",
			FileExtension: []string{".tfstate"},
		},
		false,
		false,
	)

	mock := s3Mock{}
	awsInstance.svc = &mock

	versions, err := awsInstance.GetVersions("test")
	if err != nil {
		t.Error(err)
	} else if len(versions) != 2 {
		t.Error("Expected 2 versions")
	}
}
