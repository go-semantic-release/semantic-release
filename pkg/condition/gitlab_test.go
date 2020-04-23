package condition

import (
	"testing"
)

func TestGitlabValid(t *testing.T) {
	gl := Gitlab{}
	err := gl.RunCondition(CIConfig{"defaultBranch": ""})
	if err == nil {
		t.Fail()
	}
}
