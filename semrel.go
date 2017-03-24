package semrel

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"regexp"
	"sort"
	"strings"
	"time"
)

var commitPattern = regexp.MustCompile("^(\\w*)(?:\\((.*)\\))?\\: (.*)$")
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

type Release struct {
	SHA     string
	Version *semver.Version
}

type Repository struct {
	Owner  string
	Repo   string
	Ctx    context.Context
	Client *github.Client
}

func NewRepository(ctx context.Context, slug, token string) (*Repository, error) {
	if !strings.Contains(slug, "/") {
		return nil, errors.New("invalid slug")
	}
	repo := new(Repository)
	splited := strings.Split(slug, "/")
	repo.Owner = splited[0]
	repo.Repo = splited[1]
	repo.Ctx = ctx
	repo.Client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)))
	return repo, nil
}

func (repo *Repository) GetInfo() (string, bool, error) {
	r, _, err := repo.Client.Repositories.Get(repo.Ctx, repo.Owner, repo.Repo)
	if err != nil {
		return "", false, err
	}
	return r.GetDefaultBranch(), r.GetPrivate(), nil
}

func parseCommit(commit *github.RepositoryCommit) *Commit {
	c := new(Commit)
	c.SHA = commit.GetSHA()
	c.Raw = strings.Split(commit.Commit.GetMessage(), "\n")
	found := commitPattern.FindAllStringSubmatch(c.Raw[0], -1)
	if len(found) < 1 {
		return c
	}
	c.Type = strings.ToLower(found[0][1])
	c.Scope = found[0][2]
	c.Message = found[0][3]
	c.Change = Change{
		Major: breakingPattern.MatchString(commit.Commit.GetMessage()),
		Minor: c.Type == "feat",
		Patch: c.Type == "fix",
	}
	return c
}

func (repo *Repository) GetCommits() ([]*Commit, error) {
	opts := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	commits, _, err := repo.Client.Repositories.ListCommits(repo.Ctx, repo.Owner, repo.Repo, opts)
	if err != nil {
		return nil, err
	}
	ret := make([]*Commit, len(commits))
	for i, commit := range commits {
		ret[i] = parseCommit(commit)
	}
	return ret, nil
}

func (repo *Repository) GetLatestRelease() (*Release, error) {
	opts := &github.ListOptions{PerPage: 1}
	tags, _, err := repo.Client.Repositories.ListTags(repo.Ctx, repo.Owner, repo.Repo, opts)
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return &Release{"", &semver.Version{}}, nil
	}
	version, verr := semver.NewVersion(tags[0].GetName())
	if verr != nil {
		return nil, verr
	}
	return &Release{tags[0].Commit.GetSHA(), version}, nil
}

func (repo *Repository) CreateRelease(commits []*Commit, latestRelease *Release, newVersion *semver.Version) error {
	tag := fmt.Sprintf("v%s", newVersion.String())
	changelog := GetChangelog(commits, latestRelease, newVersion)
	opts := &github.RepositoryRelease{
		TagName: &tag,
		Body:    &changelog,
	}
	_, _, err := repo.Client.Repositories.CreateRelease(repo.Ctx, repo.Owner, repo.Repo, opts)
	if err != nil {
		return err
	}
	return nil
}

func CaluclateChange(commits []*Commit, latestRelease *Release) Change {
	var change Change
	for _, commit := range commits {
		if latestRelease.SHA == commit.SHA {
			break
		}
		change.Major = change.Major || commit.Change.Major
		change.Minor = change.Minor || commit.Change.Minor
		change.Patch = change.Patch || commit.Change.Patch
	}
	return change
}

func ApplyChange(version *semver.Version, change Change) *semver.Version {
	if version.Major() == 0 {
		change.Major = true
	}
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

func GetNewVersion(commits []*Commit, latestRelease *Release) *semver.Version {
	return ApplyChange(latestRelease.Version, CaluclateChange(commits, latestRelease))
}

func trimSHA(sha string) string {
	if len(sha) < 9 {
		return sha
	}
	return sha[:8]
}

func formatCommit(c *Commit) string {
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

func GetChangelog(commits []*Commit, latestRelease *Release, newVersion *semver.Version) string {
	ret := fmt.Sprintf("## %s (%s)\n\n", newVersion.String(), time.Now().UTC().Format("2006-01-02"))
	typeScopeMap := make(map[string]string)
	for _, commit := range commits {
		if latestRelease.SHA == commit.SHA {
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
