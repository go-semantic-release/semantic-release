package manager

import (
	"os"
	"strings"

	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/condition"
	"github.com/go-semantic-release/semantic-release/v2/pkg/config"
	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/buildin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/updater"
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
	if os.Getenv("GITLAB_CI") == "true" {
		ciType = "gitlab"
	}
	cic, err := plugin.StartCIConditionPlugin(buildin.GetPluginOpts(condition.CIConditionPluginName, ciType))
	if err != nil {
		return nil, err
	}
	return cic, nil
}

func (m *PluginManager) GetProvider() (provider.Provider, error) {
	providerType := "github"
	if strings.ToLower(m.config.ProviderPlugin) == "gitlab" || os.Getenv("GITLAB_CI") == "true" {
		providerType = "gitlab"
	}

	provider, err := plugin.StartProviderPlugin(buildin.GetPluginOpts(provider.PluginName, providerType))
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (m *PluginManager) GetCommitAnalyzer() (analyzer.CommitAnalyzer, error) {
	ca, err := plugin.StartCommitAnalyzerPlugin(buildin.GetPluginOpts(analyzer.CommitAnalyzerPluginName))
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (m *PluginManager) GetChangelogGenerator() (generator.ChangelogGenerator, error) {
	cg, err := plugin.StartChangelogGeneratorPlugin(buildin.GetPluginOpts(generator.ChangelogGeneratorPluginName))
	if err != nil {
		return nil, err
	}
	return cg, nil
}

func (m *PluginManager) GetUpdater() (updater.Updater, error) {
	npmUpdater, err := plugin.StartFilesUpdaterPlugin(buildin.GetPluginOpts(updater.FilesUpdaterPluginName, "npm"))
	if err != nil {
		return nil, err
	}

	updater := &updater.ChainedUpdater{
		Updaters: []updater.FilesUpdater{npmUpdater},
	}
	return updater, nil
}

func (m *PluginManager) Stop() {
	plugin.KillAllPlugins()
}
