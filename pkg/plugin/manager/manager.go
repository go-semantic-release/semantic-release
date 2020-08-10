package manager

import (
	"os"

	"github.com/go-semantic-release/semantic-release/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/pkg/condition"
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
	ciType := "default"
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		ciType = "github"
	}
	if os.Getenv("TRAVIS") == "true" {
		ciType = "travis"
	}
	if os.Getenv("GITLAB_CI") == "true" {
		ciType = "gitlab"
	}
	cic, err := plugin.StartCIConditionPlugin(buildin.GetPluginOpts(condition.CIConditionPluginName, ciType))
	if err != nil {
		return nil, err
	}
	return cic, nil
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
