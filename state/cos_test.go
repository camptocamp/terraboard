package state

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/camptocamp/terraboard/config"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var (
	bucketName string
	accessKey  string
	secretKey  string
)

func TestMain(m *testing.M) {
	bucketName = "BucketName-APPID"
	accessKey = os.Getenv("COS_SECRET_ID")
	secretKey = os.Getenv("COS_SECRET_KEY")

	if accessKey == "" {
		fmt.Println("AccessKey is empty, use the example AK...")
		accessKey = "ACCESSKEYEXAMPLE"
	}

	if secretKey == "" {
		fmt.Println("SecretKey is empty, use the example SK...")
		secretKey = "SECRETKEYEXAMPLE/XXXXX/SECRETKEYEXAMPLE"
	}

	fmt.Println("Test begin...")
	m.Run()
	fmt.Println("Test end.")
}

func TestCosNewCOSOk(t *testing.T) {
	cosConfig := config.COSConfig{
		Bucket:    bucketName,
		Region:    "ap-guangzhou",
		KeyPrefix: "terraform/state/",
		SecretId:  accessKey,
		SecretKey: secretKey,
	}

	exts := []CosExt{}
	t.Logf("[TEST]: cosConfig:[%v], exts:[%v]", cosConfig, exts)
	cosInstance, err := NewCOS(cosConfig, exts...)
	if err != nil {
		t.Error("NewCOS failed, reason:", err)
	}

	if cosInstance == nil {
		t.Error("COS instance is nil")
	}
}

func TestCosNewCOSWithOutBucket(t *testing.T) {
	cosInstance, err := NewCOS(
		config.COSConfig{
			Bucket:    "",
			Region:    "ap-guangzhou",
			KeyPrefix: "terraform/state/",
			SecretId:  accessKey,
			SecretKey: secretKey,
		},
		nil,
	)

	if cosInstance != nil {
		t.Error("COS instance should be nil")
	}
	if err == nil {
		t.Error("An error is expected to output.")
	}
}

func TestCosNewCOSWithOutAKSK(t *testing.T) {
	exts := []CosExt{}
	cosConfig := config.COSConfig{
		Bucket:    "test",
		Region:    "ap-guangzhou",
		KeyPrefix: "terraform/state/",
		SecretId:  "",
		SecretKey: "",
	}
	_, err := NewCOS(cosConfig, exts...)

	if err == nil {
		t.Error("An error is expected to output.")
		return
	}
	if !strings.Contains(err.Error(), "missing SecretId or SecretKey") {
		t.Error("Missing the expected log output")
	}
}

func TestCosNewCOSWithOutToken(t *testing.T) {
	_, err := NewCOS(
		config.COSConfig{
			Bucket:      bucketName,
			Region:      "ap-guangzhou",
			KeyPrefix:   "terraform/state/",
			SecretToken: "",
		},
		nil,
	)

	if err == nil {
		t.Error("An error is expected to output.")
		return
	}
	if !strings.Contains(err.Error(), "missing SecretId or SecretKey") {
		t.Error("Missing the expected log output")
	}
}

func TestCosNewCOSCollection(t *testing.T) {
	provider := config.ProviderConfig{
		NoVersioning: false,
		NoLocks:      false,
	}
	config := genConfig4COS(provider)

	instances, err := NewCOSCollection(&config)
	if err != nil {
		t.Error("NewCOSCollection failed, reason:", err)
	}

	if instances == nil || len(instances) != 1 {
		t.Error("COS instances are nil or not the expected number")
	}
}

func TestCosGetLocksWithNoLocks(t *testing.T) {
	provider := config.ProviderConfig{
		NoVersioning: false,
		NoLocks:      true,
	}
	config := genConfig4COS(provider)

	cosInstances, err := NewCOSCollection(&config)
	if err != nil {
		t.Error("NewCOSCollection failed, reason:", err.Error())
	}

	if cosInstances == nil || len(cosInstances) != 1 {
		t.Error("COS instances are nil or not the expected number")
	}

	locks, _ := cosInstances[0].GetLocks()
	if len(locks) != 0 {
		t.Error("Locks should be empty due to NoLocks option")
	}
}

func TestCosGetVersionWithNoVersioning(t *testing.T) {
	provider := config.ProviderConfig{
		NoVersioning: true,
		NoLocks:      false,
	}
	config := genConfig4COS(provider)

	cosInstances, err := NewCOSCollection(&config)
	if err != nil {
		t.Error("NewCOSCollection failed, reason:", err)
	}

	if cosInstances == nil || len(cosInstances) != 1 {
		t.Error("COS instances are nil or not the expected number")
	}

	versions, _ := cosInstances[0].GetVersions("test")

	if len(versions) != 1 {
		t.Error("Expected one versions")
	}
}

func TestCosGetStates(t *testing.T) {
	cosInstance := &COS{
		bucket: BucketServiceMock{},
	}

	states, err := cosInstance.GetStates()
	if err != nil {
		t.Error("GetStates failed, reason:", err)
	}
	if len(states) == 0 {
		t.Error("States was expected but was empty actually!")
	}
}

func TestCosGetVersions(t *testing.T) {
	cosInstance := &COS{
		bucket: BucketServiceMock{},
	}

	states, err := cosInstance.GetStates()
	if err != nil {
		t.Error("GetStates failed, reason:", err)
	}
	if len(states) == 0 {
		t.Error("States was expected but was empty actually!")
	}

	versions, err := cosInstance.GetVersions(states[0])
	if err != nil {
		t.Error("GetVersions failed, reason:", err)
	}
	if len(versions) == 0 {
		t.Error("Versions was expected but was empty actually!")
	}
}

func TestCosGetState(t *testing.T) {
	cosInstance := &COS{
		bucket: BucketServiceMock{},
		object: ObjectServiceMock{},
	}

	states, _ := cosInstance.GetStates()
	vers, _ := cosInstance.GetVersions(states[0])

	state, err := cosInstance.GetState(states[0], vers[0].ID)
	if err != nil {
		t.Error("GetState failed, reason:", err)
	}
	if state == nil {
		t.Error("The specified State was expected but was nil actually!")
	}
}

func genConfig4COS(provider config.ProviderConfig) config.Config {
	cosConfig := config.COSConfig{
		Bucket:    bucketName,
		Region:    "ap-guangzhou",
		KeyPrefix: "terraform/state/",
		SecretId:  accessKey,
		SecretKey: secretKey,
	}

	config := config.Config{
		COS:            []config.COSConfig{cosConfig},
		Version:        false,
		ConfigFilePath: "",
		Provider:       provider,
		DB:             config.DBConfig{},
		AWS:            []config.AWSConfig{},
		TFE:            []config.TFEConfig{},
		GCP:            []config.GCPConfig{},
		Gitlab:         []config.GitlabConfig{},
		Web:            config.WebConfig{},
	}
	return config
}

type BucketServiceMock struct {
	*cos.BucketService
}

type ObjectServiceMock struct {
	*cos.ObjectService
}

func (b BucketServiceMock) Get(_ context.Context, _ *cos.BucketGetOptions) (*cos.BucketGetResult, *cos.Response, error) {
	return &cos.BucketGetResult{
		Contents: []cos.Object{
			{Key: "test.tfstate"},
			{Key: "test2.tfstate"},
			{Key: "test3.tfstate"},
		},
		IsTruncated: func() bool { b := false; return b }(),
	}, nil, nil
}

func (b BucketServiceMock) GetObjectVersions(_ context.Context, _ *cos.BucketGetObjectVersionsOptions) (*cos.BucketGetObjectVersionsResult, *cos.Response, error) {
	return &cos.BucketGetObjectVersionsResult{
		Version: []cos.ListVersionsResultVersion{
			{Key: "testId", VersionId: "v1", LastModified: time.Now().AddDate(0, 0, -2).Format("2006-01-02 15:04:05")},
			{Key: "testId2", VersionId: "v2", LastModified: time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")},
		},
	}, nil, nil
}

func (o ObjectServiceMock) Get(_ context.Context, _ string, _ *cos.ObjectGetOptions, _ ...string) (*cos.Response, error) {
	return &cos.Response{
		Response: &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"version": 4, "terraform_version": "1.4.5", "serial": 7}`))),
		},
	}, nil
}
