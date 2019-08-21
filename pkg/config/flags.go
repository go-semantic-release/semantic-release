package config

import "github.com/urfave/cli"

// CliFlags cli flags
var CliFlags = []cli.Flag{
	cli.StringFlag{
		Name:     "token",
		Usage:    "github token",
		EnvVar:   "GITHUB_TOKEN",
		Required: true,
	},
	cli.StringFlag{
		Name:     "slug",
		Usage:    "slug of the repository",
		EnvVar:   "TRAVIS_REPO_SLUG",
		Required: true,
	},
	cli.StringFlag{
		Name:   "changelog",
		Usage:  "creates a changelog file",
		EnvVar: "CHANGELOG",
	},
	cli.BoolFlag{
		Name:   "ghr",
		Usage:  "create a .ghr file with the parameters for ghr",
		EnvVar: "GHR",
	},
	cli.BoolFlag{
		Name:   "noci",
		Usage:  "run semantic-release locally",
		EnvVar: "NOCI",
	},
	cli.BoolFlag{
		Name:   "dry",
		Usage:  "do not create release",
		EnvVar: "DRY",
	},
	cli.BoolFlag{
		Name:   "vf",
		Usage:  "create a .version file",
		EnvVar: "VF",
	},
	cli.StringFlag{
		Name:   "update",
		Usage:  "updates the version of a certain file",
		EnvVar: "UPDATE",
	},
	cli.StringFlag{
		Name:   "ghe-host",
		Usage:  "github enterprise host",
		EnvVar: "GITHUB_ENTERPRISE_HOST",
	},
	cli.BoolFlag{
		Name:   "prerelease",
		Usage:  "flags the release as a prerelease",
		EnvVar: "PRERELEASE",
	},
	cli.BoolFlag{
		Name:   "travis-com",
		Usage:  "force semantic-release to use the travis-ci.com API endpoint",
		EnvVar: "TRAVIS_COM",
	},
}
