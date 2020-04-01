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
	return os.Getenv("GITHUB_REF")
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
	if !gha.IsBranchRef() {
		return fmt.Errorf("not running on a branch")
	}
	return nil
}
