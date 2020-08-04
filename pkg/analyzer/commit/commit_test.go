package commit

import (
	"fmt"
	"testing"

	"github.com/go-semantic-release/semantic-release/pkg/semrel"
	"github.com/stretchr/testify/require"
)

func compareCommit(c *semrel.Commit, t, s string, change semrel.Change) bool {
	if c.Type != t || c.Scope != s {
		return false
	}
	if c.Change.Major != change.Major ||
		c.Change.Minor != change.Minor ||
		c.Change.Patch != change.Patch {
		return false
	}
	return true
}

func createRawCommit(sha, message string) *semrel.RawCommit {
	return &semrel.RawCommit{
		SHA:        sha,
		RawMessage: message,
	}
}

func TestDefaultAnalyzer(t *testing.T) {
	testCases := []struct {
		RawCommit *semrel.RawCommit
		Type      string
		Scope     string
		Change    semrel.Change
	}{
		{
			createRawCommit("a", "feat: new feature"),
			"feat",
			"",
			semrel.Change{Major: false, Minor: true, Patch: false},
		},
	}

	defaultAnalyzer := &DefaultAnalyzer{}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("AnalyzeCommitMessage: %s", tc.RawCommit.RawMessage), func(t *testing.T) {
			require.True(t, compareCommit(defaultAnalyzer.analyzeSingleCommit(tc.RawCommit), tc.Type, tc.Scope, tc.Change))
		})
	}
}
