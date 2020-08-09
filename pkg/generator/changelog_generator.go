package generator

import (
	"github.com/go-semantic-release/semantic-release/pkg/semrel"
)

type ChangelogGeneratorConfig struct {
	Commits       []*semrel.Commit
	LatestRelease *semrel.Release
	NewVersion    string
}

type ChangelogGenerator interface {
	Generate(*ChangelogGeneratorConfig) string
}
