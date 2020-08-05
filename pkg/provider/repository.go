package provider

import (
	"regexp"

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
	NewVersion string
	Prerelease bool
	Branch     string
	SHA        string
}

type Repository interface {
	GetInfo() (*RepositoryInfo, error)
	GetCommits(sha string) ([]*semrel.RawCommit, error)
	GetReleases(re *regexp.Regexp) ([]*semrel.Release, error)
	CreateRelease(*RepositoryRelease) error
	Provider() string
}
