package commit

import (
	"github.com/go-semantic-release/semantic-release/pkg/semrel"
)

type Analyzer interface {
	Analyze([]*semrel.RawCommit) []*semrel.Commit
}
