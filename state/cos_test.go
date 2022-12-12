package state

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/camptocamp/terraboard/config"
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
	cosInstance := genCOSInstance(t)

	states, err := cosInstance.GetStates()
	if err != nil {
		t.Error("GetStates failed, reason:", err)
	}
	if len(states) == 0 {
		t.Error("States was expected but was empty actually!")
	}
}

func TestCosGetVersions(t *testing.T) {
	cosInstance := genCOSInstance(t)

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
	cosInstance := genCOSInstance(t)

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

func genCOSInstance(t *testing.T) *COS {
	cosConfig := config.COSConfig{
		Bucket:    bucketName,
		Region:    "ap-guangzhou",
		KeyPrefix: "terraform/state/",
		SecretId:  accessKey,
		SecretKey: secretKey,
	}

	exts := []CosExt{}
	cosInstance, err := NewCOS(cosConfig, exts...)
	if err != nil {
		t.Error("NewCOS failed, reason:", err)
	}

	if cosInstance == nil {
		t.Error("COS instances are nil.")
	}
	return cosInstance
}
