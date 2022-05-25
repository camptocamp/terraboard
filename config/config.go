package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	tfversion "github.com/hashicorp/terraform/version"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type configFlags struct {
	Version bool `short:"V" long:"version" description:"Display version."`

	ConfigFilePath string `short:"c" long:"config-file" env:"CONFIG_FILE" description:"Config File path"`

	Provider ProviderConfig `group:"General Provider Options" yaml:"provider"`

	Log LogConfig `group:"Logging Options" yaml:"log"`

	DB DBConfig `group:"Database Options" yaml:"database"`

	AWS AWSConfig `group:"AWS Options" yaml:"aws"`

	S3 S3BucketConfig `group:"S3 Options" yaml:"s3"`

	TFE TFEConfig `group:"Terraform Enterprise Options" yaml:"tfe"`

	GCP GCPConfig `group:"Google Cloud Platform Options" yaml:"gcp"`

	Gitlab GitlabConfig `group:"GitLab Options" yaml:"gitlab"`

	Web WebConfig `group:"Web" yaml:"web"`
}

// LogConfig stores the log configuration
type LogConfig struct {
	Level  string `short:"l" long:"log-level" env:"TERRABOARD_LOG_LEVEL" yaml:"level" description:"Set log level ('debug', 'info', 'warn', 'error', 'fatal', 'panic')." default:"info"`
	Format string `long:"log-format" yaml:"format" env:"TERRABOARD_LOG_FORMAT" description:"Set log format ('plain', 'json')." default:"plain"`
}

// DBConfig stores the database configuration
type DBConfig struct {
	Host         string `long:"db-host" env:"DB_HOST" yaml:"host" description:"Database host." default:"db"`
	Port         uint16 `long:"db-port" env:"DB_PORT" yaml:"port" description:"Database port." default:"5432"`
	User         string `long:"db-user" env:"DB_USER" yaml:"user" description:"Database user." default:"gorm"`
	Password     string `long:"db-password" env:"DB_PASSWORD" yaml:"password" description:"Database password."`
	Name         string `long:"db-name" env:"DB_NAME" yaml:"name" description:"Database name." default:"gorm"`
	SSLMode      string `long:"db-sslmode" env:"DB_SSLMODE" yaml:"sslmode" description:"Database SSL mode." default:"require"`
	NoSync       bool   `long:"no-sync" yaml:"no-sync" description:"Do not sync database."`
	SyncInterval uint16 `long:"sync-interval" yaml:"sync-interval" description:"DB sync interval (in minutes)" default:"1"`
}

// S3BucketConfig stores the S3 bucket configuration
type S3BucketConfig struct {
	Bucket         string   `long:"s3-bucket" env:"AWS_BUCKET" yaml:"bucket" description:"AWS S3 bucket."`
	KeyPrefix      string   `long:"key-prefix" env:"AWS_KEY_PREFIX" yaml:"key-prefix" description:"AWS Key Prefix."`
	FileExtension  []string `long:"file-extension" env:"AWS_FILE_EXTENSION" env-delim:"," yaml:"file-extension" description:"File extension(s) of state files." default:".tfstate"`
	ForcePathStyle bool     `long:"force-path-style" env:"AWS_FORCE_PATH_STYLE" yaml:"force-path-style" description:"Force path style S3 bucket calls."`
}

// AWSConfig stores the DynamoDB table and S3 Bucket configuration
type AWSConfig struct {
	AccessKey       string           `long:"aws-access-key" env:"AWS_ACCESS_KEY_ID" yaml:"access-key" description:"AWS account access key."`
	SecretAccessKey string           `long:"aws-secret-access-key" env:"AWS_SECRET_ACCESS_KEY" yaml:"secret-access-key" description:"AWS secret account access key."`
	SessionToken    string           `long:"aws-session-token" env:"AWS_SESSION_TOKEN" yaml:"session-token" description:"AWS session token."`
	DynamoDBTable   string           `long:"dynamodb-table" env:"AWS_DYNAMODB_TABLE" yaml:"dynamodb-table" description:"AWS DynamoDB table for locks."`
	S3              []S3BucketConfig `group:"S3 Options" yaml:"s3"`
	Endpoint        string           `long:"aws-endpoint" env:"AWS_ENDPOINT" yaml:"endpoint" description:"AWS endpoint."`
	Region          string           `long:"aws-region" env:"AWS_REGION" yaml:"region" description:"AWS region."`
	APPRoleArn      string           `long:"aws-role-arn" env:"APP_ROLE_ARN" yaml:"app-role-arn" description:"Role ARN to Assume."`
	ExternalID      string           `long:"aws-external-id" env:"AWS_EXTERNAL_ID" yaml:"external-id" description:"External ID to use when assuming role."`
}

// TFEConfig stores the Terraform Enterprise configuration
type TFEConfig struct {
	Address      string `long:"tfe-address" env:"TFE_ADDRESS" yaml:"address" description:"Terraform Enterprise address for states access"`
	Token        string `long:"tfe-token" env:"TFE_TOKEN" yaml:"token" description:"Terraform Enterprise Token for states access"`
	Organization string `long:"tfe-organization" env:"TFE_ORGANIZATION" yaml:"organization" description:"Terraform Enterprise organization for states access"`
}

// GCPConfig stores the Google Cloud configuration
type GCPConfig struct {
	HTTPClient *http.Client
	GCSBuckets []string `long:"gcs-bucket" yaml:"gcs-bucket" description:"Google Cloud bucket to search"`
	GCPSAKey   string   `long:"gcp-sa-key-path" env:"GCP_SA_KEY_PATH" yaml:"gcp-sa-key-path" description:"The path to the service account to use to connect to Google Cloud Platform"`
}

// GitlabConfig stores the GitLab configuration
type GitlabConfig struct {
	Address string `long:"gitlab-address" env:"GITLAB_ADDRESS" yaml:"address" description:"GitLab address (root)" default:"https://gitlab.com"`
	Token   string `long:"gitlab-token" env:"GITLAB_TOKEN" yaml:"token" description:"Token to authenticate upon GitLab"`
}

// WebConfig stores the UI interface parameters
type WebConfig struct {
	Port        uint16 `short:"p" long:"port" env:"TERRABOARD_PORT" yaml:"port" description:"Port to listen on." default:"8080"`
	SwaggerPort uint16 `long:"swagger-port" env:"TERRABOARD_SWAGGER_PORT" yaml:"swagger-port" description:"Port for swagger to listen on." default:"8081"`
	BaseURL     string `long:"base-url" env:"TERRABOARD_BASE_URL" yaml:"base-url" description:"Base URL." default:"/"`
	LogoutURL   string `long:"logout-url" env:"TERRABOARD_LOGOUT_URL" yaml:"logout-url" description:"Logout URL."`
}

// ProviderConfig stores genral provider parameters
type ProviderConfig struct {
	NoVersioning bool `long:"no-versioning" env:"TERRABOARD_NO_VERSIONING" yaml:"no-versioning" description:"Disable versioning support from Terraboard (useful for S3 compatible providers like MinIO)"`
	NoLocks      bool `long:"no-locks" env:"TERRABOARD_NO_LOCKS" yaml:"no-locks" description:"Disable locks support from Terraboard (useful for S3 compatible providers like MinIO)"`
}

// Config stores the handler's configuration and UI interface parameters
type Config struct {
	Version bool `short:"V" long:"version" description:"Display version."`

	ConfigFilePath string `short:"c" long:"config-file" env:"CONFIG_FILE" description:"Config File path"`

	Provider ProviderConfig `group:"General Provider Options" yaml:"provider"`

	Log LogConfig `group:"Logging Options" yaml:"log"`

	DB DBConfig `group:"Database Options" yaml:"database"`

	AWS []AWSConfig `group:"AWS Options" yaml:"aws"`

	TFE []TFEConfig `group:"Terraform Enterprise Options" yaml:"tfe"`

	GCP []GCPConfig `group:"Google Cloud Platform Options" yaml:"gcp"`

	Gitlab []GitlabConfig `group:"GitLab Options" yaml:"gitlab"`

	Web WebConfig `group:"Web" yaml:"web"`
}

// LoadConfigFromYaml loads the config from config file
func (c *Config) LoadConfigFromYaml(filename string) *Config {
	fmt.Printf("Loading config from %s\n", filename)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	yamlFile = []byte(os.ExpandEnv(string(yamlFile)))
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}

	c.ConfigFilePath = filename
	return c
}

// Parse flags and env variables to given struct using go-flags
// parser
func parseStructFlagsAndEnv() configFlags {
	var tmpConfig configFlags
	parser := flags.NewParser(&tmpConfig, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatalf("Failed to parse flags: %s", err)
	}

	return tmpConfig
}

// LoadConfig loads the config from flags & environment
func LoadConfig(version string) *Config {
	var c Config
	parsedConfig := parseStructFlagsAndEnv()

	if parsedConfig.Version {
		fmt.Printf("Terraboard %v (built for Terraform v%v)\n", version, tfversion.Version)
		os.Exit(0)
	}

	c = Config{
		Version:        parsedConfig.Version,
		ConfigFilePath: parsedConfig.ConfigFilePath,
		Provider:       parsedConfig.Provider,
		Log:            parsedConfig.Log,
		DB:             parsedConfig.DB,
		AWS:            []AWSConfig{parsedConfig.AWS},
		TFE:            []TFEConfig{parsedConfig.TFE},
		GCP:            []GCPConfig{parsedConfig.GCP},
		Gitlab:         []GitlabConfig{parsedConfig.Gitlab},
		Web:            parsedConfig.Web,
	}
	c.AWS[0].S3 = append(c.AWS[0].S3, parsedConfig.S3)

	if parsedConfig.ConfigFilePath != "" {
		if _, err := os.Stat(parsedConfig.ConfigFilePath); err == nil {
			c.LoadConfigFromYaml(parsedConfig.ConfigFilePath)
		} else {
			fmt.Printf("File %s doesn't exists!\n", c.ConfigFilePath)
			os.Exit(1)
		}
	}

	return &c
}

// SetupLogging sets up logging for Terraboard
func (c Config) SetupLogging() (err error) {
	switch c.Log.Level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		errMsg := fmt.Sprintf("Wrong log level '%v'", c.Log.Level)
		return errors.New(errMsg)
	}

	switch c.Log.Format {
	case "plain":
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		errMsg := fmt.Sprintf("Wrong log format '%v'", c.Log.Format)
		return errors.New(errMsg)
	}

	return
}
