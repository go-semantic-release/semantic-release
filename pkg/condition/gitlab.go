package condition

import (
	"fmt"
	"os"
)

type Gitlab struct {
}

func (gl *Gitlab) Name() string {
	return "Gitlab"
}

func (gl *Gitlab) GetCurrentBranch() string {
	return os.Getenv("CI_COMMIT_BRANCH")
}

func (gl *Gitlab) GetCurrentSHA() string {
	return os.Getenv("CI_COMMIT_SHA")
}

func (gl *Gitlab) IsBranchRef() bool {
	return gl.GetCurrentBranch() != ""
}

func (gl *Gitlab) RunCondition(config CIConfig) error {
	defaultBranch := config["defaultBranch"].(string)
	if !gl.IsBranchRef() {
		return fmt.Errorf("This test run is not running on a branch build.")
	}
	if branch := gl.GetCurrentBranch(); defaultBranch != "*" && branch != defaultBranch {
		return fmt.Errorf("This test run was triggered on the branch %s, while semantic-release is configured to only publish from %s.", branch, defaultBranch)
	}
	return nil
}
