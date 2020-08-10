package github

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-semantic-release/semantic-release/pkg/condition"
	"github.com/go-semantic-release/semantic-release/pkg/plugin"
	"github.com/urfave/cli/v2"
)

type GitHubActions struct {
}

func (gha *GitHubActions) Name() string {
	return "GitHub Actions"
}

func (gha *GitHubActions) GetCurrentBranch() string {
	return strings.TrimPrefix(os.Getenv("GITHUB_REF"), "refs/heads/")
}

func (gha *GitHubActions) GetCurrentSHA() string {
	return os.Getenv("GITHUB_SHA")
}

func (gha *GitHubActions) IsBranchRef() bool {
	if val := os.Getenv("GITHUB_REF"); val != "" {
		return strings.HasPrefix(val, "refs/heads/")
	}
	return false
}

func (gha *GitHubActions) RunCondition(config map[string]string) error {
	defaultBranch := config["defaultBranch"]
	if !gha.IsBranchRef() {
		return fmt.Errorf("This test run is not running on a branch build.")
	}
	if branch := gha.GetCurrentBranch(); defaultBranch != "*" && branch != defaultBranch {
		return fmt.Errorf("This test run was triggered on the branch %s, while semantic-release is configured to only publish from %s.", branch, defaultBranch)
	}
	return nil
}

func Main(c *cli.Context) error {
	plugin.Serve(&plugin.ServeOpts{
		CICondition: func() condition.CICondition {
			return &GitHubActions{}
		},
	})
	return nil
}
