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

// Config stores the handler's configuration and UI interface parameters
type Config struct {
	Version bool `short:"V" long:"version" description:"Display version."`

	Port int `short:"p" long:"port" yaml:"port" description:"Port to listen on." default:"8080"`

	ConfigFilePath string `long:"config-file" env:"CONFIG_FILE" description:"Config File path" default:""`

	Log struct {
		Level  string `short:"l" long:"log-level" yaml:"level" description:"Set log level ('debug', 'info', 'warn', 'error', 'fatal', 'panic')." env:"TERRABOARD_LOG_LEVEL" default:"info"`
		Format string `long:"log-format" yaml:"format" description:"Set log format ('plain', 'json')." env:"TERRABOARD_LOG_FORMAT" default:"plain"`
	} `group:"Logging Options" yaml:"log"`

	DB struct {
		Host     string `long:"db-host" env:"DB_HOST" yaml:"host" description:"Database host." default:"db"`
		User     string `long:"db-user" env:"DB_USER" yaml:"user" description:"Database user." default:"gorm"`
		Password string `long:"db-password" env:"DB_PASSWORD" yaml:"password" description:"Database password."`
		Name     string `long:"db-name" env:"DB_NAME" yaml:"name" description:"Database name." default:"gorm"`
		NoSync   bool   `long:"no-sync" yaml:"no-sync" description:"Do not sync database."`
	} `group:"Database Options" yaml:"database"`

	AWS struct {
		DynamoDBTable string `long:"dynamodb-table" env:"AWS_DYNAMODB_TABLE" yaml:"dynamodb-table" description:"AWS DynamoDB table for locks."`

		S3 struct {
			Bucket        string `long:"s3-bucket" env:"AWS_BUCKET" yaml:"bucket" description:"AWS S3 bucket."`
			KeyPrefix     string `long:"key-prefix" env:"AWS_KEY_PREFIX" yaml:"key-prefix" description:"AWS Key Prefix."`
			FileExtension string `long:"file-extension" env:"AWS_FILE_EXTENSION" yaml:"file-extension" description:"File extension of state files." default:".tfstate"`
		} `group:"S3 Options" yaml:"s3"`
	} `group:"AWS Options" yaml:"aws"`

	Web struct {
		BaseURL   string `long:"base-url" env:"TERRABOARD_BASE_URL" yaml:"base-url" description:"Base URL." default:"/"`
		LogoutURL string `long:"logout-url" env:"TERRABOARD_LOGOUT_URL" yaml:"logout-url" description:"Logout URL."`
	} `group:"Web" yaml:"web"`
}

// LoadConfigFromYaml loads the config from config file
func (c *Config) LoadConfigFromYaml(configFilePath string) *Config {
	fmt.Printf("Loading config from %s\n", c.ConfigFilePath)
	yamlFile, err := ioutil.ReadFile(configFilePath)
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
			c.LoadConfigFromYaml(c.ConfigFilePath)
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
