package manager

import (
	"os"

	"github.com/go-semantic-release/semantic-release/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/pkg/condition"
	"github.com/go-semantic-release/semantic-release/pkg/condition/defaultci"
	githubCI "github.com/go-semantic-release/semantic-release/pkg/condition/github"
	gitlabCI "github.com/go-semantic-release/semantic-release/pkg/condition/gitlab"
	"github.com/go-semantic-release/semantic-release/pkg/condition/travis"
	"github.com/go-semantic-release/semantic-release/pkg/config"
	"github.com/go-semantic-release/semantic-release/pkg/generator"
	"github.com/go-semantic-release/semantic-release/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/pkg/plugin/buildin"
	"github.com/go-semantic-release/semantic-release/pkg/provider"
	"github.com/go-semantic-release/semantic-release/pkg/provider/github"
	"github.com/go-semantic-release/semantic-release/pkg/provider/gitlab"
	"github.com/go-semantic-release/semantic-release/pkg/updater"
	"github.com/go-semantic-release/semantic-release/pkg/updater/npm"
)

type PluginManager struct {
	config *config.Config
}

func New(config *config.Config) (*PluginManager, error) {
	return &PluginManager{config}, nil
}

func (m *PluginManager) GetCICondition() (condition.CICondition, error) {
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

func (m *PluginManager) GetProvider() (provider.Repository, error) {
	if m.config.GitLab {
		return &gitlab.GitLabRepository{}, nil
	}
	return &github.GitHubRepository{}, nil
}

func (m *PluginManager) GetCommitAnalyzer() (analyzer.CommitAnalyzer, error) {
	ca, err := plugin.StartCommitAnalyzerPlugin(buildin.GetPluginOpts(analyzer.CommitAnalyzerPluginName))

	if err != nil {
		return nil, err
	}

	return ca, nil
}

func (m *PluginManager) GetChangelogGenerator() (generator.ChangelogGenerator, error) {
	return &generator.DefaultChangelogGenerator{}, nil
}

func (m *PluginManager) GetUpdater() (updater.Updater, error) {
	updater := &updater.ChainedUpdater{
		Updaters: []updater.FilesUpdater{
			&npm.Updater{},
		},
	}
	return updater, nil
}

func (m *PluginManager) Stop() {
	plugin.KillAllPlugins()
}
