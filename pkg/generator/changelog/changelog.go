package changelog

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/go-semantic-release/semantic-release/pkg/semrel"
)

type Config struct {
	Commits       []*semrel.Commit
	LatestRelease *semrel.Release
	NewVersion    *semver.Version
}

type Generator interface {
	Generate(*Config) string
}

type DefaultGenerator struct{}

func trimSHA(sha string) string {
	if len(sha) < 9 {
		return sha
	}
	return sha[:8]
}

func formatCommit(c *semrel.Commit) string {
	ret := "* "
	if c.Scope != "" {
		ret += fmt.Sprintf("**%s:** ", c.Scope)
	}
	ret += fmt.Sprintf("%s (%s)\n", c.Message, trimSHA(c.SHA))
	return ret
}

var typeToText = map[string]string{
	"feat":     "Feature",
	"fix":      "Bug Fixes",
	"perf":     "Performance Improvements",
	"revert":   "Reverts",
	"docs":     "Documentation",
	"style":    "Styles",
	"refactor": "Code Refactoring",
	"test":     "Tests",
	"chore":    "Chores",
	"%%bc%%":   "Breaking Changes",
}

func getSortedKeys(m *map[string]string) []string {
	keys := make([]string, len(*m))
	i := 0
	for k := range *m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func (*DefaultGenerator) Generate(changelogConfig *Config) string {
	ret := fmt.Sprintf("## %s (%s)\n\n", changelogConfig.NewVersion.String(), time.Now().UTC().Format("2006-01-02"))
	typeScopeMap := make(map[string]string)
	for _, commit := range changelogConfig.Commits {
		if changelogConfig.LatestRelease.SHA == commit.SHA {
			break
		}
		if commit.Change.Major {
			typeScopeMap["%%bc%%"] += fmt.Sprintf("%s\n```%s\n```\n", formatCommit(commit), strings.Join(commit.Raw[1:], "\n"))
			continue
		}
		if commit.Type == "" {
			continue
		}
		typeScopeMap[commit.Type] += formatCommit(commit)
	}
	for _, t := range getSortedKeys(&typeScopeMap) {
		msg := typeScopeMap[t]
		typeName, found := typeToText[t]
		if !found {
			typeName = t
		}
		ret += fmt.Sprintf("#### %s\n\n%s\n", typeName, msg)
	}
	return ret
}
