package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func TestSetLogging_debug(t *testing.T) {
	c := Config{}
	c.Log.Level = "debug"
	c.Log.Format = "plain"
	err := c.SetupLogging()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if log.GetLevel() != log.DebugLevel {
		t.Fatalf("Expected debug level, got %v", log.GetLevel())
	}
}

func TestLoadConfig(t *testing.T) {
	t.Skip("Skipping this test since go-flags can't parse properlly flags with go test command")

	c := LoadConfig("1.0.0")
	compareConfig := Config{
		Log: LogConfig{
			Level:  "info",
			Format: "plain",
		},
		ConfigFilePath: "",
		DB: DBConfig{
			Host:         "db",
			Port:         5432,
			User:         "gorm",
			Password:     "",
			Name:         "gorm",
			SSLMode:      "require",
			NoSync:       false,
			SyncInterval: 1,
		},
		AWS: []AWSConfig{
			{
				AccessKey:       "",
				SecretAccessKey: "",
				DynamoDBTable:   "",
				S3: []S3BucketConfig{{
					Bucket:         "",
					KeyPrefix:      "",
					FileExtension:  []string{".tfstate"},
					ForcePathStyle: false,
				}},
			},
		},
		TFE: []TFEConfig{
			{
				Address:      "",
				Token:        "",
				Organization: "",
			},
		},
		GCP: []GCPConfig{
			{
				GCSBuckets: nil,
				GCPSAKey:   "",
			},
		},
		Gitlab: []GitlabConfig{
			{
				Address: "https://gitlab.com",
				Token:   "",
			},
		},
		Web: WebConfig{
			Port:        8080,
			SwaggerPort: 8081,
			BaseURL:     "/",
			LogoutURL:   "",
		},
	}

	if !reflect.DeepEqual(*c, compareConfig) {
		t.Errorf(
			"TestLoadConfig() -> \n\ngot:\n%v,\n\nwant:\n%v",
			spew.Sdump(*c),
			spew.Sdump(compareConfig),
		)
	}
}

func TestLoadConfigFromYaml(t *testing.T) {
	var config Config
	os.Setenv("AWS_DEFAULT_REGION", "test-region")
	defer os.Unsetenv("AWS_DEFAULT_REGION")
	config.LoadConfigFromYaml("config_test.yml")
	compareConfig := Config{
		Log: LogConfig{
			Level:  "error",
			Format: "json",
		},
		ConfigFilePath: "config_test.yml",
		DB: DBConfig{
			Host:         "postgres",
			Port:         15432,
			User:         "terraboard-user",
			Password:     "terraboard-pass",
			Name:         "terraboard-db",
			SSLMode:      "",
			NoSync:       true,
			SyncInterval: 0,
		},
		AWS: []AWSConfig{
			{
				AccessKey:       "root",
				SecretAccessKey: "mypassword",
				DynamoDBTable:   "terraboard-dynamodb",
				Region:          "test-region",
				S3: []S3BucketConfig{{
					Bucket:         "terraboard-bucket",
					KeyPrefix:      "test/",
					FileExtension:  []string{".tfstate"},
					ForcePathStyle: true,
				}},
			},
		},
		TFE: []TFEConfig{
			{
				Address:      "https://tfe.example.com",
				Token:        "foo",
				Organization: "bar",
			},
		},
		GCP: []GCPConfig{
			{
				GCSBuckets: []string{"my-bucket-1", "my-bucket-2"},
				GCPSAKey:   "/path/to/key",
			},
		},
		Gitlab: []GitlabConfig{
			{
				Address: "https://gitlab.example.com",
				Token:   "foo",
			},
		},
		Web: WebConfig{
			Port:        39090,
			SwaggerPort: 8081,
			BaseURL:     "/test/",
			LogoutURL:   "/test-logout",
		},
	}

	if !reflect.DeepEqual(config, compareConfig) {
		t.Errorf(
			"TestLoadConfig() -> \n\ngot:\n%v,\n\nwant:\n%v",
			spew.Sdump(config),
			spew.Sdump(compareConfig),
		)
	}
}

func TestSetLogging_info(t *testing.T) {
	c := Config{}
	c.Log.Level = "info"
	c.Log.Format = "plain"
	err := c.SetupLogging()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if log.GetLevel() != log.InfoLevel {
		t.Fatalf("Expected info level, got %v", log.GetLevel())
	}
}

func TestSetLogging_warn(t *testing.T) {
	c := Config{}
	c.Log.Level = "warn"
	c.Log.Format = "plain"
	err := c.SetupLogging()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if log.GetLevel() != log.WarnLevel {
		t.Fatalf("Expected warn level, got %v", log.GetLevel())
	}
}

func TestSetLogging_error(t *testing.T) {
	c := Config{}
	c.Log.Level = "error"
	c.Log.Format = "plain"
	err := c.SetupLogging()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if log.GetLevel() != log.ErrorLevel {
		t.Fatalf("Expected error level, got %v", log.GetLevel())
	}
}

func TestSetLogging_fatal(t *testing.T) {
	c := Config{}
	c.Log.Level = "fatal"
	c.Log.Format = "plain"
	err := c.SetupLogging()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if log.GetLevel() != log.FatalLevel {
		t.Fatalf("Expected fatal level, got %v", log.GetLevel())
	}
}

func TestSetLogging_panic(t *testing.T) {
	c := Config{}
	c.Log.Level = "panic"
	c.Log.Format = "plain"
	err := c.SetupLogging()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if log.GetLevel() != log.PanicLevel {
		t.Fatalf("Expected panic level, got %v", log.GetLevel())
	}
}

func TestSetLogging_wronglevel(t *testing.T) {
	c := Config{}
	c.Log.Level = "wrong"
	c.Log.Format = "plain"
	err := c.SetupLogging()

	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedError := "Wrong log level 'wrong'"

	if err.Error() != expectedError {
		t.Fatalf("Expected %s, got %s", expectedError, err.Error())
	}
}

func TestSetLogging_json(t *testing.T) {
	c := Config{}
	c.Log.Level = "debug"
	c.Log.Format = "json"
	err := c.SetupLogging()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSetLogging_wrongformat(t *testing.T) {
	c := Config{}
	c.Log.Level = "debug"
	c.Log.Format = "yaml"
	err := c.SetupLogging()

	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedError := "Wrong log format 'yaml'"

	if err.Error() != expectedError {
		t.Fatalf("Expected %s, got %s", expectedError, err.Error())
	}
}
