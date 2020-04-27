package condition

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitlabValid(t *testing.T) {
	gl := GitLab{}
	os.Setenv("CI_COMMIT_BRANCH", "")
	err := gl.RunCondition(CIConfig{"defaultBranch": ""})
	assert.EqualError(t, err, "This test run is not running on a branch build.")
}
