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
	config    *config.Config
	discovery *discovery.Discovery
}

func New(config *config.Config) (*PluginManager, error) {
	dis, err := discovery.New(config)
	if err != nil {
		return nil, err
	}
	return &PluginManager{
		config:    config,
		discovery: dis,
	}, nil
}

func (m *PluginManager) GetCICondition() (condition.CICondition, error) {
	opts, err := m.discovery.FindPlugin(condition.CIConditionPluginName, m.config.CIConditionPlugin)
	if err != nil {
		return nil, err
	}

	cic, err := plugin.StartPlugin(opts)
	if err != nil {
		return nil, err
	}
	return cic.(condition.CICondition), nil
}

func (m *PluginManager) GetProvider() (provider.Provider, error) {
	opts, err := m.discovery.FindPlugin(provider.PluginName, m.config.ProviderPlugin)
	if err != nil {
		return nil, err
	}

	prov, err := plugin.StartPlugin(opts)
	if err != nil {
		return nil, err
	}
	return prov.(provider.Provider), nil
}

func (m *PluginManager) GetCommitAnalyzer() (analyzer.CommitAnalyzer, error) {
	opts, err := m.discovery.FindPlugin(analyzer.CommitAnalyzerPluginName, m.config.CommitAnalyzerPlugin)
	if err != nil {
		return nil, err
	}

	ca, err := plugin.StartPlugin(opts)
	if err != nil {
		return nil, err
	}
	return ca.(analyzer.CommitAnalyzer), nil
}

func (m *PluginManager) GetChangelogGenerator() (generator.ChangelogGenerator, error) {
	opts, err := m.discovery.FindPlugin(generator.ChangelogGeneratorPluginName, m.config.ChangelogGeneratorPlugin)
	if err != nil {
		return nil, err
	}

	cg, err := plugin.StartPlugin(opts)
	if err != nil {
		return nil, err
	}
	return cg.(generator.ChangelogGenerator), nil
}

func (m *PluginManager) GetChainedUpdater() (*updater.ChainedUpdater, error) {
	updaters := make([]updater.FilesUpdater, 0)
	for _, pl := range m.config.FilesUpdaterPlugins {
		opts, err := m.discovery.FindPlugin(updater.FilesUpdaterPluginName, pl)
		if err != nil {
			return nil, err
		}

		upd, err := plugin.StartPlugin(opts)
		if err != nil {
			return nil, err
		}
		updaters = append(updaters, upd.(updater.FilesUpdater))
	}

	updater := &updater.ChainedUpdater{
		Updaters: updaters,
	}
	return updater, nil
}

func (m *PluginManager) GetChainedHooksExecutor() (*hooks.ChainedHooksExecutor, error) {
	hooksChain := make([]hooks.Hooks, 0)
	for _, pl := range m.config.HooksPlugins {
		opts, err := m.discovery.FindPlugin(hooks.PluginName, pl)
		if err != nil {
			return nil, err
		}

		hp, err := plugin.StartPlugin(opts)
		if err != nil {
			return nil, err
		}
		hooksChain = append(hooksChain, hp.(hooks.Hooks))
	}

	return &hooks.ChainedHooksExecutor{
		HooksChain: hooksChain,
	}, nil
}

func (m *PluginManager) Stop() {
	plugin.KillAllPlugins()
}

func (m *PluginManager) FetchAllPlugins() error {
	pluginMap := map[string]string{
		condition.CIConditionPluginName:        m.config.CIConditionPlugin,
		provider.PluginName:                    m.config.ProviderPlugin,
		analyzer.CommitAnalyzerPluginName:      m.config.CommitAnalyzerPlugin,
		generator.ChangelogGeneratorPluginName: m.config.ChangelogGeneratorPlugin,
	}
	for t, name := range pluginMap {
		_, err := m.discovery.FindPlugin(t, name)
		if err != nil {
			return err
		}
	}

	for _, pl := range m.config.FilesUpdaterPlugins {
		_, err := m.discovery.FindPlugin(updater.FilesUpdaterPluginName, pl)
		if err != nil {
			return err
		}
	}

	for _, pl := range m.config.HooksPlugins {
		_, err := m.discovery.FindPlugin(hooks.PluginName, pl)
		if err != nil {
			return err
		}
	}
	return nil
}
