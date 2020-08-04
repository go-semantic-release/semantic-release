package github

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubValid(t *testing.T) {
	gha := GitHubActions{}
	os.Setenv("GITHUB_REF", "")
	err := gha.RunCondition(map[string]interface{}{"defaultBranch": ""})
	assert.EqualError(t, err, "This test run is not running on a branch build.")
}
