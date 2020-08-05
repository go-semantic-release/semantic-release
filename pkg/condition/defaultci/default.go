package defaultci

import (
	"io/ioutil"
	"strings"
)

func ReadGitHead() string {
	data, err := ioutil.ReadFile(".git/HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(string(data), "ref: refs/heads/"))
}

type DefaultCI struct {
}

func (d DefaultCI) Name() string {
	return "none"
}

func (d DefaultCI) RunCondition(map[string]interface{}) error {
	return nil
}

func (d DefaultCI) GetCurrentBranch() string {
	return ReadGitHead()
}

func (d DefaultCI) GetCurrentSHA() string {
	return ReadGitHead()
}
