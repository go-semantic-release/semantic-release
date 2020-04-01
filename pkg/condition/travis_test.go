package condition

import (
	"testing"
)

func TestTravisValid(t *testing.T) {
	travis := TravisCI{}
	err := travis.RunCondition(CIConfig{"token": "", "defaultBranch": "", "private": false})
	if err == nil {
		t.Fail()
	}
}
