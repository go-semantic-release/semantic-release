package semrel

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	gitlab "github.com/xanzy/go-gitlab"
)

type GitLabRepository struct {
	owner     string
	repo      string
	projectID string
	branch    string
	Ctx       context.Context
	client    *gitlab.Client
}

func NewGitLabRepository(ctx context.Context, gitlabBaseUrl, slug, token, branch string, projectID string) (*GitLabRepository, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project id is required")
	}

	repo := new(GitLabRepository)
	repo.projectID = projectID
	repo.Ctx = ctx
	repo.branch = branch

	if strings.Contains(slug, "/") {
		split := strings.Split(slug, "/")
		repo.owner = split[0]
		repo.repo = split[1]
	}

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

func (repo *GitLabRepository) GetInfo() (string, bool, error) {
	project, _, err := repo.client.Projects.GetProject(repo.projectID, nil)

	if err != nil {
		return "", false, err
	}

	return project.DefaultBranch, project.Visibility == gitlab.PrivateVisibility, nil
}

func (repo *GitLabRepository) GetCommits(sha string) ([]*Commit, error) {
	opts := &gitlab.ListCommitsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
		RefName: gitlab.String(fmt.Sprintf("%s...%s", repo.branch, sha)),
		All:     gitlab.Bool(true),
	}

	allCommits := make([]*Commit, 0)

	for {
		commits, resp, err := repo.client.Commits.ListCommits(repo.projectID, opts)

		if err != nil {
			return nil, err
		}

		for _, commit := range commits {
			allCommits = append(allCommits, parseGitlabCommit(commit))
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opts.Page = resp.NextPage
	}

	return allCommits, nil
}

func (repo *GitLabRepository) GetReleases(re *regexp.Regexp) (Releases, error) {
	allReleases := make(Releases, 0)

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

			allReleases = append(allReleases, &Release{
				SHA:     tag.Commit.ID,
				Version: version,
			})
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opts.Page = resp.NextPage
	}

	return allReleases, nil
}

func (repo *GitLabRepository) CreateRelease(changelog string, newVersion *semver.Version, prerelease bool, branch, sha string) error {
	tag := fmt.Sprintf("v%s", newVersion.String())

	// Gitlab does not have any notion of pre-releases
	_, _, err := repo.client.Releases.CreateRelease(repo.projectID, &gitlab.CreateReleaseOptions{
		TagName: &tag,
		Ref:     &sha,
		// TODO: this may been to be wrapped in ```
		Description: &changelog,
	})

	return err
}

func parseGitlabCommit(commit *gitlab.Commit) *Commit {
	c := new(Commit)
	c.SHA = commit.ID
	c.Raw = strings.Split(commit.Message, "\n")
	found := commitPattern.FindAllStringSubmatch(c.Raw[0], -1)
	if len(found) < 1 {
		return c
	}
	c.Type = strings.ToLower(found[0][1])
	c.Scope = found[0][2]
	c.Message = found[0][3]
	c.Change = Change{
		Major: breakingPattern.MatchString(commit.Message),
		Minor: c.Type == "feat",
		Patch: c.Type == "fix",
	}
	return c
}

func (repo *GitLabRepository) Owner() string {
	return repo.owner
}

func (repo *GitLabRepository) Repo() string {
	return repo.repo
}

func (repo *GitLabRepository) Provider() string {
	return "GitLab"
}
