package semrel

import (
	"context"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"regexp"
	"strings"
)

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

type Release struct {
	SHA     string
	Version *semver.Version
}

type Repository struct {
	Owner, Repo string
}

var commitPattern = regexp.MustCompile("^(\\w*)(?:\\((.*)\\))?\\: (.*)$")
var breakingPattern = regexp.MustCompile("BREAKING CHANGES?")

func ParseCommit(commit *github.RepositoryCommit) *Commit {
	c := new(Commit)
	c.SHA = *commit.SHA
	c.Raw = strings.Split(*commit.Commit.Message, "\n")
	found := commitPattern.FindAllStringSubmatch(c.Raw[0], -1)
	if len(found) < 1 {
		return c
	}
	c.Type = strings.ToLower(found[0][1])
	c.Scope = found[0][2]
	c.Message = found[0][3]
	c.Change = Change{
		Major: breakingPattern.MatchString(*commit.Commit.Message),
		Minor: c.Type == "feat",
		Patch: c.Type == "fix",
	}
	return c
}

func GetCommits(ctx context.Context, client *github.Client, repo *Repository) ([]*Commit, error) {
	opts := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	commits, _, err := client.Repositories.ListCommits(ctx, repo.Owner, repo.Repo, opts)
	if err != nil {
		return nil, err
	}
	ret := make([]*Commit, len(commits))
	for i, commit := range commits {
		ret[i] = ParseCommit(commit)
	}
	return ret, nil
}

func GetLatestRelease(ctx context.Context, client *github.Client, repo *Repository) (*Release, error) {
	opts := &github.ListOptions{PerPage: 1}
	tags, _, err := client.Repositories.ListTags(ctx, repo.Owner, repo.Repo, opts)
	if err != nil {
		return nil, err
	}
	v, _ := semver.NewVersion(*tags[0].Name)
	return &Release{*tags[0].Commit.SHA, v}, nil
}

func GetChange(commits []*Commit, release *Release) Change {
	var change Change
	for _, commit := range commits {
		if release.SHA == commit.SHA {
			break
		}
		change.Major = change.Major || commit.Change.Major
		change.Minor = change.Minor || commit.Change.Minor
		change.Patch = change.Patch || commit.Change.Patch
	}
	return change
}

func ApplyChange(version *semver.Version, change Change) *semver.Version {
	var newVersion semver.Version
	switch {
	case change.Major:
		newVersion = version.IncMajor()
		break
	case change.Minor:
		newVersion = version.IncMinor()
		break
	case change.Patch:
		newVersion = version.IncPatch()
		break
	default:
		return nil
	}
	return &newVersion
}
