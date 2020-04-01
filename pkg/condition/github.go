package condition

import (
	"fmt"
	"os"
	"strings"
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

func (gha *GitHubActions) RunCondition(config CIConfig) error {
	defaultBranch := config["defaultBranch"].(string)
	if !gha.IsBranchRef() {
		return fmt.Errorf("This test run is not running on a branch build.")
	}
	if branch := gha.GetCurrentBranch(); defaultBranch != "*" && branch != defaultBranch {
		return fmt.Errorf("This test run was triggered on the branch %s, while semantic-release is configured to only publish from %s.", branch, defaultBranch)
	}
	return nil
}
