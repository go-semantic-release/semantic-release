package semrel

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v30/github"
	"golang.org/x/oauth2"
)

type GitHubRepository struct {
	owner  string
	repo   string
	Ctx    context.Context
	Client *github.Client
}

func NewGitHubRepository(ctx context.Context, gheHost, slug, token string) (*GitHubRepository, error) {
	if !strings.Contains(slug, "/") {
		return nil, errors.New("invalid slug")
	}
	repo := new(GitHubRepository)
	split := strings.Split(slug, "/")
	repo.owner = split[0]
	repo.repo = split[1]
	repo.Ctx = ctx
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	if gheHost != "" {
		gheUrl := fmt.Sprintf("https://%s/api/v3/", gheHost)
		rClient, err := github.NewEnterpriseClient(gheUrl, gheUrl, oauthClient)
		if err != nil {
			return nil, err
		}
		repo.Client = rClient
	} else {
		repo.Client = github.NewClient(oauthClient)
	}
	return repo, nil
}

func (repo *GitHubRepository) GetInfo() (string, bool, error) {
	r, _, err := repo.Client.Repositories.Get(repo.Ctx, repo.owner, repo.repo)
	if err != nil {
		return "", false, err
	}
	return r.GetDefaultBranch(), r.GetPrivate(), nil
}

func (repo *GitHubRepository) GetCommits(sha string) ([]*Commit, error) {
	opts := &github.CommitsListOptions{
		SHA:         sha,
		ListOptions: github.ListOptions{PerPage: 100},
	}
	commits, _, err := repo.Client.Repositories.ListCommits(repo.Ctx, repo.owner, repo.repo, opts)
	if err != nil {
		return nil, err
	}
	ret := make([]*Commit, len(commits))
	for i, commit := range commits {
		ret[i] = parseGithubCommit(commit)
	}
	return ret, nil
}

func (repo *GitHubRepository) GetLatestRelease(vrange string, re *regexp.Regexp) (*Release, error) {
	allReleases := make(Releases, 0)
	opts := &github.ReferenceListOptions{Type: "tags", ListOptions: github.ListOptions{PerPage: 100}}
	for {
		refs, resp, err := repo.Client.Git.ListRefs(repo.Ctx, repo.owner, repo.repo, opts)
		if resp != nil && resp.StatusCode == 404 {
			return &Release{"", &semver.Version{}}, nil
		}
		if err != nil {
			return nil, err
		}
		for _, r := range refs {
			tag := strings.TrimPrefix(r.GetRef(), "refs/tags/")
			if re != nil && !re.MatchString(tag) {
				continue
			}
			if r.Object.GetType() != "commit" {
				continue
			}
			version, err := semver.NewVersion(tag)
			if err != nil {
				continue
			}
			allReleases = append(allReleases, &Release{r.Object.GetSHA(), version})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allReleases.GetLatestRelease(vrange)
}

func (repo *GitHubRepository) CreateRelease(changelog string, newVersion *semver.Version, prerelease bool, branch, sha string) error {
	tag := fmt.Sprintf("v%s", newVersion.String())
	isPrerelease := prerelease || newVersion.Prerelease() != ""

	if branch != sha {
		ref := "refs/tags/" + tag
		tagOpts := &github.Reference{
			Ref:    &ref,
			Object: &github.GitObject{SHA: &sha},
		}
		_, _, err := repo.Client.Git.CreateRef(repo.Ctx, repo.owner, repo.repo, tagOpts)
		if err != nil {
			return err
		}
	}

	opts := &github.RepositoryRelease{
		TagName:         &tag,
		Name:            &tag,
		TargetCommitish: &branch,
		Body:            &changelog,
		Prerelease:      &isPrerelease,
	}
	_, _, err := repo.Client.Repositories.CreateRelease(repo.Ctx, repo.owner, repo.repo, opts)
	if err != nil {
		return err
	}
	return nil
}

func parseGithubCommit(commit *github.RepositoryCommit) *Commit {
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

func (repo *GitHubRepository) Owner() string {
	return repo.owner
}

func (repo *GitHubRepository) Repo() string {
	return repo.repo
}

func (repo *GitHubRepository) Provider() string {
	return "GitHub"
}
