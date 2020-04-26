package condition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubValid(t *testing.T) {
	gha := GitHubActions{}
	err := gha.RunCondition(CIConfig{"defaultBranch": ""})
	assert.EqualError(t, err, "This test run is not running on a branch build.")
}
