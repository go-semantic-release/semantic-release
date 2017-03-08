package semrel

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"regexp"
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

func NewRepository(ctx context.Context, slug, token string) *Repository {
	repo := new(Repository)
	splited := strings.Split(slug, "/")
	repo.Owner = splited[0]
	repo.Repo = splited[1]
	repo.Ctx = ctx
	repo.Client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)))
	return repo
}

func (repo *Repository) GetDefaultBranch() (string, error) {
	r, _, err := repo.Client.Repositories.Get(repo.Ctx, repo.Owner, repo.Repo)
	if err != nil {
		return "", err
	}
	return r.GetDefaultBranch(), nil
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
	v, _ := semver.NewVersion(tags[0].GetName())
	return &Release{tags[0].Commit.GetSHA(), v}, nil
}

func (repo *Repository) CreateRelease(commits []*Commit, latestRelease *Release, newVersion *semver.Version) error {
	tag := fmt.Sprintf("v%s", newVersion.String())
	sha := commits[0].SHA
	changelog := GetChangelog(commits, latestRelease, newVersion)
	opts := &github.RepositoryRelease{
		TagName:         &tag,
		TargetCommitish: &sha,
		Body:            &changelog,
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

func GetChangelog(commits []*Commit, latestRelease *Release, newVersion *semver.Version) string {
	ret := fmt.Sprintf("## %s (%s)\n\n", newVersion.String(), time.Now().UTC().Format("2006-01-02"))
	for _, commit := range commits {
		if latestRelease.SHA == commit.SHA {
			break
		}
		if commit.Type == "" {
			continue
		}
		ret += fmt.Sprintf("%s (%s)\n", commit.Raw[0], commit.SHA[:8])
	}
	return ret
}
