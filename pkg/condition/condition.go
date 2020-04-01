package condition

import (
	"io/ioutil"
	"os"
	"strings"
)

func ReadGitHead() string {
	data, err := ioutil.ReadFile(".git/HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(string(data), "ref: refs/heads/"))
}

func GetDefaultRepoSlug() string {
	if val := os.Getenv("TRAVIS_REPO_SLUG"); val != "" {
		return val
	}
	if val := os.Getenv("GITHUB_REPOSITORY"); val != "" {
		return val
	}
	return ""
}

type CIConfig map[string]interface{}

type CI interface {
	Name() string
	RunCondition(config CIConfig) error
	GetCurrentBranch() string
	GetCurrentSHA() string
}

type DefaultCI struct {
}

func (d DefaultCI) Name() string {
	return "none"
}

func (d DefaultCI) RunCondition(config CIConfig) error {
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
		return &GitHubActions{}
	}
	if os.Getenv("TRAVIS") == "true" {
		return &TravisCI{}
	}
	return &DefaultCI{}
}
