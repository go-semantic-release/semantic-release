package manager

import (
	"os"

	"github.com/go-semantic-release/semantic-release/pkg/analyzer/commit"
	"github.com/go-semantic-release/semantic-release/pkg/condition"
	"github.com/go-semantic-release/semantic-release/pkg/condition/defaultci"
	githubCI "github.com/go-semantic-release/semantic-release/pkg/condition/github"
	gitlabCI "github.com/go-semantic-release/semantic-release/pkg/condition/gitlab"
	"github.com/go-semantic-release/semantic-release/pkg/condition/travis"
	"github.com/go-semantic-release/semantic-release/pkg/config"
	"github.com/go-semantic-release/semantic-release/pkg/generator/changelog"
	"github.com/go-semantic-release/semantic-release/pkg/provider"
	"github.com/go-semantic-release/semantic-release/pkg/provider/github"
	"github.com/go-semantic-release/semantic-release/pkg/provider/gitlab"
)

type Manager struct {
	config *config.Config
}

func New(config *config.Config) (*Manager, error) {
	return &Manager{config}, nil
}

func (m *Manager) GetCICondition() (condition.CI, error) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return &githubCI.GitHubActions{}, nil
	}
	if os.Getenv("TRAVIS") == "true" {
		return &travis.TravisCI{}, nil
	}
	if os.Getenv("GITLAB_CI") == "true" {
		return &gitlabCI.GitLab{}, nil
	}
	return &defaultci.DefaultCI{}, nil
}

func (m *Manager) GetProvider() (provider.Repository, error) {
	if m.config.GitLab {
		return &gitlab.GitLabRepository{}, nil
	}
	return &github.GitHubRepository{}, nil
}

func (m *Manager) GetCommitAnalyzer() (commit.Analyzer, error) {
	return &commit.DefaultAnalyzer{}, nil
}

func (m *Manager) GetChangelogGenerator() (changelog.Generator, error) {
	return &changelog.DefaultGenerator{}, nil
}
