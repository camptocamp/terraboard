package config

/*********************************************
 * Custom UnmarshalYAML used to define some struct fields
 * default values where go-flags ones aren't applicable
 * (and so makes them optional)
 *********************************************/

func (s *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawConfig Config
	raw := rawConfig{
		DB: DBConfig{
			Host:         "db",
			Port:         5432,
			User:         "gorm",
			Name:         "gorm",
			SSLMode:      "require",
			SyncInterval: 1,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "plain",
		},
		Web: WebConfig{
			Port:        8080,
			SwaggerPort: 8081,
			BaseURL:     "/",
		},
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*s = Config(raw)
	return nil
}

func (s *S3BucketConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawS3BucketConfig S3BucketConfig
	raw := rawS3BucketConfig{
		FileExtension: []string{".tfstate"},
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*s = S3BucketConfig(raw)
	return nil
}

func (s *GitlabConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawGitlabConfig GitlabConfig
	raw := rawGitlabConfig{
		Address: "https://gitlab.com",
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*s = GitlabConfig(raw)
	return nil
}
