package analyzer

import (
	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
)

type CommitAnalyzer interface {
	Analyze([]*semrel.RawCommit) []*semrel.Commit
}
