package changelog

import (
	"strings"
	"testing"

	"github.com/go-semantic-release/semantic-release/pkg/semrel"
)

func TestDefaultGenerator(t *testing.T) {
	changelogConfig := &Config{}
	changelogConfig.Commits = []*semrel.Commit{
		{},
		{SHA: "123456789", Type: "feat", Scope: "app", Message: "commit message"},
		{SHA: "abcd", Type: "fix", Scope: "", Message: "commit message"},
		{SHA: "12345678", Type: "yolo", Scope: "swag", Message: "commit message"},
		{SHA: "12345678", Type: "chore", Scope: "", Message: "commit message", Raw: []string{"", "BREAKING CHANGE: test"}, Change: &semrel.Change{Major: true}},
		{SHA: "stop", Type: "chore", Scope: "", Message: "not included"},
	}
	changelogConfig.LatestRelease = &semrel.Release{SHA: "stop"}
	changelogConfig.NewVersion = "2.0.0"
	generator := &DefaultGenerator{}
	changelog := generator.Generate(changelogConfig)
	if !strings.Contains(changelog, "* **app:** commit message (12345678)") ||
		!strings.Contains(changelog, "* commit message (abcd)") ||
		!strings.Contains(changelog, "#### yolo") ||
		!strings.Contains(changelog, "```BREAKING CHANGE: test\n```") ||
		strings.Contains(changelog, "not included") {
		t.Fail()
	}
}
