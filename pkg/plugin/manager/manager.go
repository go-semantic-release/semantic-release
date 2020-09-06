package manager

import (
	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/condition"
	"github.com/go-semantic-release/semantic-release/v2/pkg/config"
	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery"
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
	opts, err := discovery.FindPlugin(condition.CIConditionPluginName, m.config.CIConditionPlugin)
	if err != nil {
		return nil, err
	}

	cic, err := plugin.StartCIConditionPlugin(opts)
	if err != nil {
		return nil, err
	}
	return cic, nil
}

func (m *PluginManager) GetProvider() (provider.Provider, error) {
	opts, err := discovery.FindPlugin(provider.PluginName, m.config.ProviderPlugin)
	if err != nil {
		return nil, err
	}

	provider, err := plugin.StartProviderPlugin(opts)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (m *PluginManager) GetCommitAnalyzer() (analyzer.CommitAnalyzer, error) {
	opts, err := discovery.FindPlugin(analyzer.CommitAnalyzerPluginName, m.config.CommitAnalyzerPlugin)
	if err != nil {
		return nil, err
	}

	ca, err := plugin.StartCommitAnalyzerPlugin(opts)
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (m *PluginManager) GetChangelogGenerator() (generator.ChangelogGenerator, error) {
	opts, err := discovery.FindPlugin(generator.ChangelogGeneratorPluginName, m.config.ChangelogGeneratorPlugin)
	if err != nil {
		return nil, err
	}

	cg, err := plugin.StartChangelogGeneratorPlugin(opts)
	if err != nil {
		return nil, err
	}
	return cg, nil
}

func (m *PluginManager) GetChainedUpdater() (*updater.ChainedUpdater, error) {
	updaters := make([]updater.FilesUpdater, 0)
	for _, pl := range m.config.FilesUpdaterPlugins {
		opts, err := discovery.FindPlugin(updater.FilesUpdaterPluginName, pl)
		if err != nil {
			return nil, err
		}

		upd, err := plugin.StartFilesUpdaterPlugin(opts)
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

func (m *PluginManager) GetChainedHooksExecutor() (*hooks.ChainedHooksExecutor, error) {
	hooksChain := make([]hooks.Hooks, 0)
	for _, pl := range m.config.HooksPlugins {
		opts, err := discovery.FindPlugin(hooks.PluginName, pl)
		if err != nil {
			return nil, err
		}

		hp, err := plugin.StartHooksPlugin(opts)
		if err != nil {
			return nil, err
		}
		hooksChain = append(hooksChain, hp)
	}

	return &hooks.ChainedHooksExecutor{
		HooksChain: hooksChain,
	}, nil
}

func (m *PluginManager) Stop() {
	plugin.KillAllPlugins()
}
