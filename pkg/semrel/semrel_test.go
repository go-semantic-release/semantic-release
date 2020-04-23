package semrel

import (
	"strings"
	"testing"

	"github.com/Masterminds/semver"
)

func TestCaluclateChange(t *testing.T) {
	commits := []*Commit{
		{SHA: "a", Change: Change{true, false, false}},
		{SHA: "b", Change: Change{false, true, false}},
		{SHA: "c", Change: Change{false, false, true}},
	}
	change := CaluclateChange(commits, &Release{})
	if !change.Major || !change.Minor || !change.Patch {
		t.Fail()
	}
	change = CaluclateChange(commits, &Release{SHA: "a"})
	if change.Major || change.Minor || change.Patch {
		t.Fail()
	}
	version, _ := semver.NewVersion("1.0.0")
	newVersion := GetNewVersion(commits, &Release{SHA: "b", Version: version})
	if newVersion.String() != "2.0.0" {
		t.Fail()
	}
}

func TestApplyChange(t *testing.T) {
	version, _ := semver.NewVersion("1.0.0")
	newVersion := ApplyChange(version, Change{false, false, false})
	if newVersion != nil {
		t.Fail()
	}
	newVersion = ApplyChange(version, Change{false, false, true})
	if newVersion.String() != "1.0.1" {
		t.Fail()
	}
	newVersion = ApplyChange(version, Change{false, true, true})
	if newVersion.String() != "1.1.0" {
		t.Fail()
	}
	newVersion = ApplyChange(version, Change{true, true, true})
	if newVersion.String() != "2.0.0" {
		t.Fail()
	}
	version, _ = semver.NewVersion("0.1.0")
	newVersion = ApplyChange(version, Change{})
	if newVersion.String() != "1.0.0" {
		t.Fail()
	}
	version, _ = semver.NewVersion("2.0.0-beta")
	newVersion = ApplyChange(version, Change{true, true, true})
	if newVersion.String() != "2.0.0-beta.1" {
		t.Fail()
	}
	version, _ = semver.NewVersion("2.0.0-beta.2")
	newVersion = ApplyChange(version, Change{true, true, true})
	if newVersion.String() != "2.0.0-beta.3" {
		t.Fail()
	}
	version, _ = semver.NewVersion("2.0.0-beta.1.1")
	newVersion = ApplyChange(version, Change{true, true, true})
	if newVersion.String() != "2.0.0-beta.2" {
		t.Fail()
	}
}

func TestGetChangelog(t *testing.T) {
	commits := []*Commit{
		{},
		{SHA: "123456789", Type: "feat", Scope: "app", Message: "commit message"},
		{SHA: "abcd", Type: "fix", Scope: "", Message: "commit message"},
		{SHA: "12345678", Type: "yolo", Scope: "swag", Message: "commit message"},
		{SHA: "12345678", Type: "chore", Scope: "", Message: "commit message", Raw: []string{"", "BREAKING CHANGE: test"}, Change: Change{Major: true}},
		{SHA: "stop", Type: "chore", Scope: "", Message: "not included"},
	}
	latestRelease := &Release{SHA: "stop"}
	newVersion, _ := semver.NewVersion("2.0.0")
	changelog := GetChangelog(commits, latestRelease, newVersion)
	if !strings.Contains(changelog, "* **app:** commit message (12345678)") ||
		!strings.Contains(changelog, "* commit message (abcd)") ||
		!strings.Contains(changelog, "#### yolo") ||
		!strings.Contains(changelog, "```BREAKING CHANGE: test\n```") ||
		strings.Contains(changelog, "not included") {
		t.Fail()
	}
}

func compareCommit(c *Commit, t, s string, change Change) bool {
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
