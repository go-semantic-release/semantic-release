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
	pluginInfo, err := m.discovery.FindPlugin(condition.CIConditionPluginName, m.config.CiConditionPlugin)
	if err != nil {
		return nil, err
	}

	cic, err := plugin.StartPlugin(pluginInfo)
	if err != nil {
		return nil, err
	}
	return cic.(condition.CICondition), nil
}

func (m *PluginManager) GetProvider() (provider.Provider, error) {
	pluginInfo, err := m.discovery.FindPlugin(provider.PluginName, m.config.ProviderPlugin)
	if err != nil {
		return nil, err
	}

	prov, err := plugin.StartPlugin(pluginInfo)
	if err != nil {
		return nil, err
	}
	return prov.(provider.Provider), nil
}

func (m *PluginManager) GetCommitAnalyzer() (analyzer.CommitAnalyzer, error) {
	pluginInfo, err := m.discovery.FindPlugin(analyzer.CommitAnalyzerPluginName, m.config.CommitAnalyzerPlugin)
	if err != nil {
		return nil, err
	}

	ca, err := plugin.StartPlugin(pluginInfo)
	if err != nil {
		return nil, err
	}
	return ca.(analyzer.CommitAnalyzer), nil
}

func (m *PluginManager) GetChangelogGenerator() (generator.ChangelogGenerator, error) {
	pluginInfo, err := m.discovery.FindPlugin(generator.ChangelogGeneratorPluginName, m.config.ChangelogGeneratorPlugin)
	if err != nil {
		return nil, err
	}

	cg, err := plugin.StartPlugin(pluginInfo)
	if err != nil {
		return nil, err
	}
	return cg.(generator.ChangelogGenerator), nil
}

func (m *PluginManager) GetChainedUpdater() (*updater.ChainedUpdater, error) {
	updaters := make([]updater.FilesUpdater, 0)
	for _, pl := range m.config.FilesUpdaterPlugins {
		pluginInfo, err := m.discovery.FindPlugin(updater.FilesUpdaterPluginName, pl)
		if err != nil {
			return nil, err
		}

		upd, err := plugin.StartPlugin(pluginInfo)
		if err != nil {
			return nil, err
		}
		updaters = append(updaters, upd.(updater.FilesUpdater))
	}

	return &updater.ChainedUpdater{
		Updaters: updaters,
	}, nil
}

func (m *PluginManager) GetChainedHooksExecutor() (*hooks.ChainedHooksExecutor, error) {
	hooksChain := make([]hooks.Hooks, 0)
	for _, pl := range m.config.HooksPlugins {
		pluginInfo, err := m.discovery.FindPlugin(hooks.PluginName, pl)
		if err != nil {
			return nil, err
		}

		hp, err := plugin.StartPlugin(pluginInfo)
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

func (m *PluginManager) getAllPlugins() [][]string {
	plugins := make([][]string, 0, 4)
	// required plugins
	plugins = append(plugins, []string{condition.CIConditionPluginName, m.config.CiConditionPlugin})
	plugins = append(plugins, []string{provider.PluginName, m.config.ProviderPlugin})
	plugins = append(plugins, []string{analyzer.CommitAnalyzerPluginName, m.config.CommitAnalyzerPlugin})
	plugins = append(plugins, []string{generator.ChangelogGeneratorPluginName, m.config.ChangelogGeneratorPlugin})

	// optional plugins
	for _, pl := range m.config.FilesUpdaterPlugins {
		plugins = append(plugins, []string{updater.FilesUpdaterPluginName, pl})
	}
	for _, pl := range m.config.HooksPlugins {
		plugins = append(plugins, []string{hooks.PluginName, pl})
	}
	return plugins
}

func (m *PluginManager) getAllPluginInfos() ([]*plugin.Info, error) {
	infos := make([]*plugin.Info, 0, 4)
	for _, pl := range m.getAllPlugins() {
		pInfo, err := plugin.GetPluginInfo(pl[0], pl[1])
		if err != nil {
			return nil, err
		}
		infos = append(infos, pInfo)
	}
	return infos, nil
}

func (m *PluginManager) checkIfSameResolvers(infos []*plugin.Info) (string, bool) {
	resolver := ""
	for _, info := range infos {
		if resolver == "" {
			resolver = info.Resolver
		} else if resolver != info.Resolver {
			return "", false
		}
	}
	return resolver, true
}

func (m *PluginManager) PrefetchAllPluginsIfBatchIsPossible() (bool, []*plugin.Info, error) {
	pInfos, err := m.getAllPluginInfos()
	if err != nil {
		return false, nil, err
	}

	if m.config.PluginResolverDisableBatchPrefetch {
		return false, pInfos, nil
	}

	if resolver, ok := m.checkIfSameResolvers(pInfos); ok && m.discovery.IsBatchResolver(resolver) {
		// all plugins have the same resolver, and it supports batch resolving
		bErr := m.discovery.FindPluginsWithBatchResolver(resolver, pInfos)
		if bErr != nil {
			return false, pInfos, bErr
		}
		return true, pInfos, nil
	}
	return false, pInfos, nil
}

func (m *PluginManager) FetchAllPlugins() error {
	batchWasPossible, pInfos, err := m.PrefetchAllPluginsIfBatchIsPossible()
	if err != nil && pInfos == nil {
		return err
	}
	if batchWasPossible {
		return nil
	}

	// try to find plugins one by one
	for _, pInfo := range pInfos {
		fErr := m.discovery.FindPluginByPluginInfo(pInfo)
		if fErr != nil {
			return fErr
		}
	}
	return nil
}
