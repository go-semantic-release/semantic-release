package config

import (
	"github.com/urfave/cli/v2"
)

// CliFlags cli flags
var CliFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "token",
		Usage:    "github or gitlab token",
		EnvVars:  []string{"GITHUB_TOKEN", "GH_TOKEN", "GITLAB_TOKEN"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     "slug",
		Usage:    "slug of the repository",
		EnvVars:  []string{"GITHUB_REPOSITORY", "TRAVIS_REPO_SLUG", "CI_PROJECT_PATH_SLUG"},
		Required: true,
	},
	&cli.StringFlag{
		Name:  "changelog",
		Usage: "creates a changelog file",
	},
	&cli.BoolFlag{
		Name:  "ghr",
		Usage: "create a .ghr file with the parameters for ghr",
	},
	&cli.BoolFlag{
		Name:  "noci",
		Usage: "run semantic-release locally",
	},
	&cli.BoolFlag{
		Name:  "dry",
		Usage: "do not create release",
	},
	&cli.BoolFlag{
		Name:  "vf",
		Usage: "create a .version file",
	},
	&cli.StringFlag{
		Name:  "update",
		Usage: "updates the version of a certain file",
	},
	&cli.StringFlag{
		Name:    "ghe-host",
		Usage:   "github enterprise host",
		EnvVars: []string{"GITHUB_ENTERPRISE_HOST"},
	},
	&cli.BoolFlag{
		Name:  "prerelease",
		Usage: "flags the release as a prerelease",
	},
	&cli.BoolFlag{
		Name:  "travis-com",
		Usage: "force semantic-release to use the travis-ci.com API endpoint",
	},
	&cli.StringFlag{
		Name:  "match",
		Usage: "Only consider tags matching the given glob(7) pattern, excluding the \"refs/tags/\" prefix.",
	},
	&cli.BoolFlag{
		Name:    "gitlab",
		Usage:   "run semantic-release on GitLab",
		EnvVars: []string{"GITLAB_CI"},
	},
	&cli.StringFlag{
		Name:    "gitlab-base-url",
		Usage:   "GitLab self hosted api path",
		EnvVars: []string{"CI_SERVER_URL"},
	},
	&cli.StringFlag{
		Name:    "gitlab-project-id",
		Usage:   "GitLab project unique id",
		EnvVars: []string{"CI_PROJECT_ID"},
	},
	&cli.BoolFlag{
		Name:  "allow-initial-development-versions",
		Usage: "semantic-release will start your initial development release at 0.1.0",
	},
	&cli.BoolFlag{
		Name:  "allow-no-changes",
		Value: false,
		Usage: "Exit with code 0 if no changes are found, useful if semantic-release is automatically run",
	},
}
