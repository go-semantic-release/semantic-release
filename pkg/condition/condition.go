package condition

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-semantic-release/semantic-release/pkg/condition/github"
	"github.com/go-semantic-release/semantic-release/pkg/condition/gitlab"
	"github.com/go-semantic-release/semantic-release/pkg/condition/travis"
)

func ReadGitHead() string {
	data, err := ioutil.ReadFile(".git/HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(string(data), "ref: refs/heads/"))
}

type CI interface {
	Name() string
	RunCondition(map[string]interface{}) error
	GetCurrentBranch() string
	GetCurrentSHA() string
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

func NewCI() CI {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return &github.GitHubActions{}
	}
	if os.Getenv("TRAVIS") == "true" {
		return &travis.TravisCI{}
	}
	if os.Getenv("GITLAB_CI") == "true" {
		return &gitlab.GitLab{}
	}
	return &DefaultCI{}
}
