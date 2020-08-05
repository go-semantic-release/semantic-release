package provider

import (
	"github.com/go-semantic-release/semantic-release/pkg/semrel"
)

type Repository interface {
	Init(map[string]string) error
	GetInfo() (*RepositoryInfo, error)
	GetCommits(sha string) ([]*semrel.RawCommit, error)
	GetReleases(re string) ([]*semrel.Release, error)
	CreateRelease(*RepositoryRelease) error
	Provider() string
}
