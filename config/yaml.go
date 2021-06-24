package config

/*********************************************
 * Custom UnmarshalYAML used to define some struct fields
 * default values where go-flags ones aren't applicable
 * (and so makes them optional)
 *********************************************/

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
	raw := GitlabConfig{
		Address: "https://gitlab.com",
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*s = GitlabConfig(raw)
	return nil
}
