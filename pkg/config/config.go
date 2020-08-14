package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// Config is a complete set of app configuration
type Config struct {
	Token                           string
	Slug                            string
	Changelog                       string
	Ghr                             bool
	Noci                            bool
	Dry                             bool
	Vf                              bool
	Update                          string
	GheHost                         string
	Prerelease                      bool
	Match                           string
	AllowInitialDevelopmentVersions bool
	AllowNoChanges                  bool
	GitLab                          bool
	GitLabBaseURL                   string
	GitLabProjectID                 string
	MaintainedVersion               string
}

// NewConfig returns a new Config instance
func NewConfig(c *cli.Context) (*Config, error) {
	conf := &Config{
		Token:                           c.String("token"),
		Slug:                            c.String("slug"),
		Changelog:                       c.String("changelog"),
		Ghr:                             c.Bool("ghr"),
		Noci:                            c.Bool("noci"),
		Dry:                             c.Bool("dry"),
		Vf:                              c.Bool("vf"),
		Update:                          c.String("update"),
		GheHost:                         c.String("ghe-host"),
		Prerelease:                      c.Bool("prerelease"),
		Match:                           c.String("match"),
		AllowInitialDevelopmentVersions: c.Bool("allow-initial-development-versions"),
		AllowNoChanges:                  c.Bool("allow-no-changes"),
		GitLab:                          c.Bool("gitlab"),
		GitLabBaseURL:                   c.String("gitlab-base-url"),
		GitLabProjectID:                 c.String("gitlab-project-id"),
		MaintainedVersion:               "",
	}

	f, err := os.OpenFile(".semrelrc", os.O_RDONLY, 0)
	if err != nil {
		return conf, nil
	}
	defer f.Close()

	var configFile map[string]string
	if err := json.NewDecoder(f).Decode(&configFile); err != nil {
		return nil, fmt.Errorf("could not parse .semrelrc: %w", err)
	}

	conf.MaintainedVersion = configFile["maintainedVersion"]

	return conf, nil
}
