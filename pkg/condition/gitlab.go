package condition

import (
	"fmt"
	"os"
)

type GitLab struct {
}

func (gl *GitLab) Name() string {
	return "GitLab CI"
}

func (gl *GitLab) GetCurrentBranch() string {
	return os.Getenv("CI_COMMIT_BRANCH")
}

func (gl *GitLab) GetCurrentSHA() string {
	return os.Getenv("CI_COMMIT_SHA")
}

func (gl *GitLab) IsBranchRef() bool {
	return gl.GetCurrentBranch() != ""
}

func (gl *GitLab) RunCondition(config CIConfig) error {
	defaultBranch := config["defaultBranch"].(string)
	if !gl.IsBranchRef() {
		return fmt.Errorf("This test run is not running on a branch build.")
	}
	if branch := gl.GetCurrentBranch(); defaultBranch != "*" && branch != defaultBranch {
		return fmt.Errorf("This test run was triggered on the branch %s, while semantic-release is configured to only publish from %s.", branch, defaultBranch)
	}
	return nil
}
