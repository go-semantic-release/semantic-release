package provider

import (
	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
)

type Provider interface {
	Init(map[string]string) error
	GetInfo() (*RepositoryInfo, error)
	GetCommits(sha string) ([]*semrel.RawCommit, error)
	GetReleases(re string) ([]*semrel.Release, error)
	CreateRelease(*CreateReleaseConfig) error
	Name() string
}
