package condition

import (
	"testing"
)

func TestGithubValid(t *testing.T) {
	gha := GitHubActions{}
	err := gha.RunCondition(CIConfig{"defaultBranch": ""})
	if err == nil {
		t.Fail()
	}
}
