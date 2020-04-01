package condition

import (
	"os"
	"testing"
)

func TestGithubValid(t *testing.T) {
	os.Setenv("GITHUB_REF", "")
	gha := GitHubActions{}
	err := gha.RunCondition(CIConfig{})
	if err == nil {
		t.Fail()
	}
}
