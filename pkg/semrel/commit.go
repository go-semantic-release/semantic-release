package semrel

import (
	"regexp"
	"strings"
)

var commitPattern = regexp.MustCompile(`^(\w*)(?:\((.*)\))?\: (.*)$`)
var breakingPattern = regexp.MustCompile("BREAKING CHANGES?")

type Change struct {
	Major, Minor, Patch bool
}

type Commit struct {
	SHA     string
	Raw     []string
	Type    string
	Scope   string
	Message string
	Change  Change
}

func NewCommit(sha, msg string) *Commit {
	c := new(Commit)
	c.SHA = sha
	c.Raw = strings.Split(msg, "\n")
	found := commitPattern.FindAllStringSubmatch(c.Raw[0], -1)
	if len(found) < 1 {
		return c
	}
	c.Type = strings.ToLower(found[0][1])
	c.Scope = found[0][2]
	c.Message = found[0][3]
	c.Change = Change{
		Major: breakingPattern.MatchString(msg),
		Minor: c.Type == "feat",
		Patch: c.Type == "fix",
	}
	return c
}
