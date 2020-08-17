package analyzer

import (
	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
)

type CommitAnalyzer interface {
	Init(map[string]string) error
	Name() string
	Version() string
	Analyze([]*semrel.RawCommit) []*semrel.Commit
}
