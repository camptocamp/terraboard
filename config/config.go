package config

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// Config stores the handler's configuration and UI interface parameters
type Config struct {
	Version bool `short:"V" long:"version" description:"Display version."`

	Port int `short:"p" long:"port" description:"Port to listen on." default:"8080"`

	DB struct {
		Host     string `long:"db-host" env:"DB_HOST" description:"Database host." default:"db"`
		User     string `long:"db-user" env:"DB_USER" description:"Database user." default:"gorm"`
		Password string `long:"db-password" env:"DB_PASSWORD" description:"Database password."`
		Name     string `long:"db-name" env:"DB_NAME" description:"Database name." default:"gorm"`
		NoSync   bool   `long:"no-sync" description:"Do not sync database."`
	} `group:"Database Options"`

	S3 struct {
		Bucket        string `long:"s3-bucket" env:"AWS_BUCKET" description:"AWS S3 bucket."`
		DynamoDBTable string `long:"dynamodb-table" env:"AWS_DYNAMODB_TABLE" description:"AWS DynamoDB table for locks."`
	} `group:"AWS S3 Options"`
}

// LoadConfig loads the config from flags & environment
func LoadConfig(version string) *Config {
	var c Config
	parser := flags.NewParser(&c, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	if c.Version {
		fmt.Printf("Terraboard v%v\n", version)
		os.Exit(0)
	}

	return &c
}
