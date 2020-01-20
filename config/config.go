package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

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
	SSLMode      string `long:"db-sslmode" yaml:"sslmode" description:"Database SSL mode." default:"disable"`
	NoSync       bool   `long:"no-sync" yaml:"no-sync" description:"Do not sync database."`
	SyncInterval uint16 `long:"sync-interval" yaml:"sync-interval" description:"DB sync interval (in minutes)" default:"1"`
}

// S3BucketConfig stores the S3 bucket configuration
type S3BucketConfig struct {
	Bucket        string `long:"s3-bucket" env:"AWS_BUCKET" yaml:"bucket" description:"AWS S3 bucket."`
	KeyPrefix     string `long:"key-prefix" env:"AWS_KEY_PREFIX" yaml:"key-prefix" description:"AWS Key Prefix."`
	FileExtension string `long:"file-extension" env:"AWS_FILE_EXTENSION" yaml:"file-extension" description:"File extension of state files." default:".tfstate"`
}

// AWSConfig stores the DynamoDB table and S3 Bucket configuration
type AWSConfig struct {
	DynamoDBTable string         `long:"dynamodb-table" env:"AWS_DYNAMODB_TABLE" yaml:"dynamodb-table" description:"AWS DynamoDB table for locks."`
	S3            S3BucketConfig `group:"S3 Options" yaml:"s3"`
}

// TFEConfig stores the Terraform Enterprise configuration
type TFEConfig struct {
	Address      string `long:"tfe-address" env:"TFE_ADDRESS" yaml:"tfe-address" description:"Terraform Enterprise address for states access"`
	Token        string `long:"tfe-token" env:"TFE_TOKEN" yaml:"tfe-token" description:"Terraform Enterprise Token for states access"`
	Organization string `long:"tfe-organization" env:"TFE_ORGANIZATION" yaml:"tfe-organization" description:"Terraform Enterprise organization for states access"`
}

// WebConfig stores the UI interface parameters
type WebConfig struct {
	Port      uint16 `short:"p" long:"port" env:"TERRABOARD_PORT" yaml:"port" description:"Port to listen on." default:"8080"`
	BaseURL   string `long:"base-url" env:"TERRABOARD_BASE_URL" yaml:"base-url" description:"Base URL." default:"/"`
	LogoutURL string `long:"logout-url" env:"TERRABOARD_LOGOUT_URL" yaml:"logout-url" description:"Logout URL."`
}

// Config stores the handler's configuration and UI interface parameters
type Config struct {
	Version bool `short:"V" long:"version" description:"Display version."`

	ConfigFilePath string `short:"c" long:"config-file" env:"CONFIG_FILE" description:"Config File path"`

	Log LogConfig `group:"Logging Options" yaml:"log"`

	DB DBConfig `group:"Database Options" yaml:"database"`

	AWS AWSConfig `group:"AWS Options" yaml:"aws"`

	TFE TFEConfig `group:"Terraform Enterprise Options" yaml:"tfe"`

	Web WebConfig `group:"Web" yaml:"web"`
}

// LoadConfigFromYaml loads the config from config file
func (c *Config) LoadConfigFromYaml() *Config {
	fmt.Printf("Loading config from %s\n", c.ConfigFilePath)
	yamlFile, err := ioutil.ReadFile(c.ConfigFilePath)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}

	return c
}

// LoadConfig loads the config from flags & environment
func LoadConfig(version string) *Config {
	var c Config
	parser := flags.NewParser(&c, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	if c.ConfigFilePath != "" {
		if _, err := os.Stat(c.ConfigFilePath); err == nil {
			c.LoadConfigFromYaml()
		} else {
			fmt.Printf("File %s doesn't exists!\n", c.ConfigFilePath)
			os.Exit(1)
		}
	}

	if c.Version {
		fmt.Printf("Terraboard v%v\n", version)
		os.Exit(0)
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
