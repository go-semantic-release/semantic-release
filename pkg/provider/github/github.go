package github

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
	"github.com/google/go-github/v32/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

type GitHubRepository struct {
	owner  string
	repo   string
	client *github.Client
}

func (repo *GitHubRepository) Init(config map[string]string) error {
	gheHost := config["githubEnterpriseHost"]
	slug := config["slug"]
	token := config["token"]
	if !strings.Contains(slug, "/") {
		return errors.New("invalid slug")
	}
	split := strings.Split(slug, "/")
	repo.owner = split[0]
	repo.repo = split[1]
	oauthClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	if gheHost != "" {
		gheUrl := fmt.Sprintf("https://%s/api/v3/", gheHost)
		rClient, err := github.NewEnterpriseClient(gheUrl, gheUrl, oauthClient)
		if err != nil {
			return err
		}
		repo.client = rClient
	} else {
		repo.client = github.NewClient(oauthClient)
	}
	return nil
}

func (repo *GitHubRepository) GetInfo() (*provider.RepositoryInfo, error) {
	r, _, err := repo.client.Repositories.Get(context.Background(), repo.owner, repo.repo)
	if err != nil {
		return nil, err
	}
	return &provider.RepositoryInfo{
		Owner:         r.GetOwner().GetLogin(),
		Repo:          r.GetName(),
		DefaultBranch: r.GetDefaultBranch(),
		Private:       r.GetPrivate(),
	}, nil
}

func (repo *GitHubRepository) GetCommits(sha string) ([]*semrel.RawCommit, error) {
	opts := &github.CommitsListOptions{
		SHA:         sha,
		ListOptions: github.ListOptions{PerPage: 100},
	}
	commits, _, err := repo.client.Repositories.ListCommits(context.Background(), repo.owner, repo.repo, opts)
	if err != nil {
		return nil, err
	}
	ret := make([]*semrel.RawCommit, len(commits))
	for i, commit := range commits {
		ret[i] = &semrel.RawCommit{
			SHA:        commit.GetSHA(),
			RawMessage: commit.Commit.GetMessage(),
		}
	}
	return ret, nil
}

func (repo *GitHubRepository) GetReleases(rawRe string) ([]*semrel.Release, error) {
	re := regexp.MustCompile(rawRe)
	allReleases := make([]*semrel.Release, 0)
	opts := &github.ReferenceListOptions{Ref: "tags", ListOptions: github.ListOptions{PerPage: 100}}
	for {
		refs, resp, err := repo.client.Git.ListMatchingRefs(context.Background(), repo.owner, repo.repo, opts)
		if resp != nil && resp.StatusCode == 404 {
			return allReleases, nil
		}
		if err != nil {
			return nil, err
		}
		for _, r := range refs {
			tag := strings.TrimPrefix(r.GetRef(), "refs/tags/")
			if rawRe != "" && !re.MatchString(tag) {
				continue
			}
			if r.Object.GetType() != "commit" {
				continue
			}
			version, err := semver.NewVersion(tag)
			if err != nil {
				continue
			}
			allReleases = append(allReleases, &semrel.Release{SHA: r.Object.GetSHA(), Version: version.String()})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allReleases, nil
}

func (repo *GitHubRepository) CreateRelease(release *provider.CreateReleaseConfig) error {
	tag := fmt.Sprintf("v%s", release.NewVersion)
	isPrerelease := release.Prerelease || semver.MustParse(release.NewVersion).Prerelease() != ""

	if release.Branch != release.SHA {
		ref := "refs/tags/" + tag
		tagOpts := &github.Reference{
			Ref:    &ref,
			Object: &github.GitObject{SHA: &release.SHA},
		}
		_, _, err := repo.client.Git.CreateRef(context.Background(), repo.owner, repo.repo, tagOpts)
		if err != nil {
			return err
		}
	}

	opts := &github.RepositoryRelease{
		TagName:         &tag,
		Name:            &tag,
		TargetCommitish: &release.Branch,
		Body:            &release.Changelog,
		Prerelease:      &isPrerelease,
	}
	_, _, err := repo.client.Repositories.CreateRelease(context.Background(), repo.owner, repo.repo, opts)
	return err
}

func (repo *GitHubRepository) Name() string {
	return "GitHub"
}

func Main(c *cli.Context) error {
	plugin.Serve(&plugin.ServeOpts{
		Provider: func() provider.Provider {
			return &GitHubRepository{}
		},
	})
	return nil
}
