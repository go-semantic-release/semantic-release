package commit

import (
	"regexp"
	"strings"

	"github.com/go-semantic-release/semantic-release/pkg/semrel"
)

type Analyzer interface {
	Analyze([]*semrel.RawCommit) []*semrel.Commit
}

var commitPattern = regexp.MustCompile(`^(\w*)(?:\((.*)\))?\: (.*)$`)
var breakingPattern = regexp.MustCompile("BREAKING CHANGES?")

type DefaultAnalyzer struct{}

func (da *DefaultAnalyzer) analyzeSingleCommit(rawCommit *semrel.RawCommit) *semrel.Commit {
	c := new(semrel.Commit)
	c.SHA = rawCommit.SHA
	c.Raw = strings.Split(rawCommit.RawMessage, "\n")
	found := commitPattern.FindAllStringSubmatch(c.Raw[0], -1)
	if len(found) < 1 {
		return c
	}
	c.Type = strings.ToLower(found[0][1])
	c.Scope = found[0][2]
	c.Message = found[0][3]
	c.Change = semrel.Change{
		Major: breakingPattern.MatchString(rawCommit.RawMessage),
		Minor: c.Type == "feat",
		Patch: c.Type == "fix",
	}
	return c
}

func (da *DefaultAnalyzer) Analyze(rawCommits []*semrel.RawCommit) []*semrel.Commit {
	ret := make([]*semrel.Commit, len(rawCommits))
	for i, c := range rawCommits {
		ret[i] = da.analyzeSingleCommit(c)
	}
	return ret
}
