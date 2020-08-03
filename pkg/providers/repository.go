package providers

import (
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/go-semantic-release/semantic-release/pkg/semrel"
)

type RepositoryInfo struct {
	Owner         string
	Repo          string
	DefaultBranch string
	Private       bool
}

type RepositoryRelease struct {
	Changelog  string
	NewVersion *semver.Version
	Prerelease bool
	Branch     string
	SHA        string
}

type Repository interface {
	GetInfo() (*RepositoryInfo, error)
	GetCommits(sha string) ([]*semrel.Commit, error)
	GetReleases(re *regexp.Regexp) (semrel.Releases, error)
	CreateRelease(*RepositoryRelease) error
	Provider() string
}
