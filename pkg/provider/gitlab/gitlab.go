package gitlab

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/go-semantic-release/semantic-release/pkg/provider"
	"github.com/go-semantic-release/semantic-release/pkg/semrel"
	"github.com/xanzy/go-gitlab"
)

type GitLabRepository struct {
	projectID string
	branch    string
	Ctx       context.Context
	client    *gitlab.Client
}

func NewRepository(ctx context.Context, gitlabBaseUrl, token, branch string, projectID string) (*GitLabRepository, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project id is required")
	}

	repo := new(GitLabRepository)
	repo.projectID = projectID
	repo.Ctx = ctx
	repo.branch = branch

	var (
		client *gitlab.Client
		err    error
	)

	if gitlabBaseUrl != "" {
		client, err = gitlab.NewClient(token, gitlab.WithBaseURL(gitlabBaseUrl))
	} else {
		client, err = gitlab.NewClient(token)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	repo.client = client

	return repo, nil
}

func (repo *GitLabRepository) GetInfo() (*provider.RepositoryInfo, error) {
	project, _, err := repo.client.Projects.GetProject(repo.projectID, nil)

	if err != nil {
		return nil, err
	}
	return &provider.RepositoryInfo{
		Owner:         "",
		Repo:          "",
		DefaultBranch: project.DefaultBranch,
		Private:       project.Visibility == gitlab.PrivateVisibility,
	}, nil
}

func (repo *GitLabRepository) GetCommits(sha string) ([]*semrel.RawCommit, error) {
	opts := &gitlab.ListCommitsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
		RefName: gitlab.String(fmt.Sprintf("%s...%s", repo.branch, sha)),
		All:     gitlab.Bool(true),
	}

	allCommits := make([]*semrel.RawCommit, 0)

	for {
		commits, resp, err := repo.client.Commits.ListCommits(repo.projectID, opts)

		if err != nil {
			return nil, err
		}

		for _, commit := range commits {
			allCommits = append(allCommits, &semrel.RawCommit{
				SHA:        commit.ID,
				RawMessage: commit.Message,
			})
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opts.Page = resp.NextPage
	}

	return allCommits, nil
}

func (repo *GitLabRepository) GetReleases(re *regexp.Regexp) ([]*semrel.Release, error) {
	allReleases := make([]*semrel.Release, 0)

	opts := &gitlab.ListTagsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	for {
		tags, resp, err := repo.client.Tags.ListTags(repo.projectID, opts)
		if err != nil {
			return nil, err
		}

		for _, tag := range tags {
			if re != nil && !re.MatchString(tag.Name) {
				continue
			}

			version, err := semver.NewVersion(tag.Name)
			if err != nil {
				continue
			}

			allReleases = append(allReleases, &semrel.Release{
				SHA:     tag.Commit.ID,
				Version: version.String(),
			})
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opts.Page = resp.NextPage
	}

	return allReleases, nil
}

func (repo *GitLabRepository) CreateRelease(release *provider.RepositoryRelease) error {
	tag := fmt.Sprintf("v%s", release.NewVersion)

	// Gitlab does not have any notion of pre-releases
	_, _, err := repo.client.Releases.CreateRelease(repo.projectID, &gitlab.CreateReleaseOptions{
		TagName: &tag,
		Ref:     &release.SHA,
		// TODO: this may been to be wrapped in ```
		Description: &release.Changelog,
	})

	return err
}

func (repo *GitLabRepository) Provider() string {
	return "GitLab"
}
