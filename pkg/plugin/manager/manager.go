package manager

import (
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
	cic, err := plugin.StartCIConditionPlugin(buildin.GetPluginOpts(condition.CIConditionPluginName, m.config.CIConditionPlugin))
	if err != nil {
		return nil, err
	}
	return cic, nil
}

func (m *PluginManager) GetProvider() (provider.Provider, error) {
	provider, err := plugin.StartProviderPlugin(buildin.GetPluginOpts(provider.PluginName, m.config.ProviderPlugin))
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
	updaters := make([]updater.FilesUpdater, 0)
	for _, pl := range m.config.FilesUpdaterPlugins {
		upd, err := plugin.StartFilesUpdaterPlugin(buildin.GetPluginOpts(updater.FilesUpdaterPluginName, pl))
		if err != nil {
			return nil, err
		}
		updaters = append(updaters, upd)
	}

	updater := &updater.ChainedUpdater{
		Updaters: updaters,
	}
	return updater, nil
}

func (m *PluginManager) Stop() {
	plugin.KillAllPlugins()
}
