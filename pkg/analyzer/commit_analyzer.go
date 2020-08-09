package analyzer

import (
	"github.com/go-semantic-release/semantic-release/pkg/semrel"
)

type CommitAnalyzer interface {
	Analyze([]*semrel.RawCommit) []*semrel.Commit
}
