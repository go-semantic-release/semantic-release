package discovery

import (
	"errors"
	"fmt"

	"github.com/go-semantic-release/semantic-release/v2/pkg/config"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver/github"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver/registry"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver/registryV1"
)

type Discovery struct {
	config    *config.Config
	resolvers map[string]resolver.Resolver
}

func loadResolvers(resolvers ...resolver.Resolver) (map[string]resolver.Resolver, error) {
	resolversMap := make(map[string]resolver.Resolver)
	for _, r := range resolvers {
		for _, name := range r.Names() {
			if name == "default" {
				return nil, fmt.Errorf("resolver name default is reserved")
			}
			if _, ok := resolversMap[name]; ok {
				return nil, fmt.Errorf("resolver %s already exists", name)
			}
			resolversMap[name] = r
		}
	}
	return resolversMap, nil
}

func New(config *config.Config) (*Discovery, error) {
	resolvers, err := loadResolvers(
		registry.NewResolver(config.PluginResolverEndpoint),
		registryV1.NewResolver(config.PluginResolverEndpoint),
		github.NewResolver(),
	)
	if err != nil {
		return nil, err
	}
	// use the registry resolver as default
	resolvers["default"] = resolvers[config.PluginResolver]

	if resolvers["default"] == nil {
		return nil, fmt.Errorf("resolver %s does not exist", config.PluginResolver)
	}

	return &Discovery{
		config:    config,
		resolvers: resolvers,
	}, nil
}

func (d *Discovery) fetchPlugin(pluginInfo *plugin.Info) (string, error) {
	pluginResolver, ok := d.resolvers[pluginInfo.Resolver]
	if !ok {
		return "", fmt.Errorf("resolver %s not found", pluginInfo.Resolver)
	}

	downloadInfo, err := pluginResolver.ResolvePlugin(pluginInfo)
	if err != nil {
		return "", err
	}

	return downloadPlugin(pluginInfo, downloadInfo, d.config.ShowProgress)
}

func (d *Discovery) IsBatchResolver(resolverName string) bool {
	_, ok := d.resolvers[resolverName].(resolver.BatchResolver)
	return ok
}

func (d *Discovery) FindPluginByPluginInfo(pInfo *plugin.Info) error {
	err := setAndEnsurePluginPath(pInfo)
	if err != nil {
		return err
	}

	binPath, err := findPluginLocally(pInfo)
	if errors.Is(err, ErrPluginNotFound) {
		binPath, err = d.fetchPlugin(pInfo)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	pInfo.BinPath = binPath
	return nil
}

func (d *Discovery) FindPlugin(t, name string) (*plugin.Info, error) {
	pInfo, err := plugin.GetPluginInfo(t, name)
	if err != nil {
		return nil, err
	}
	err = d.FindPluginByPluginInfo(pInfo)
	if err != nil {
		return nil, err
	}
	return pInfo, nil
}

func (d *Discovery) FindPluginsWithBatchResolver(resolverName string, pInfos []*plugin.Info) error {
	if !d.IsBatchResolver(resolverName) {
		return fmt.Errorf("resolver %s does not support batch resolving", resolverName)
	}
	missingPlugins := make([]*plugin.Info, 0)
	for _, pInfo := range pInfos {
		err := setAndEnsurePluginPath(pInfo)
		if err != nil {
			return err
		}

		binPath, err := findPluginLocally(pInfo)
		if errors.Is(err, ErrPluginNotFound) {
			missingPlugins = append(missingPlugins, pInfo)
			continue
		} else if err != nil {
			return err
		}
		pInfo.BinPath = binPath
	}

	batchResolver := d.resolvers[resolverName].(resolver.BatchResolver)
	batchDownloadInfo, err := batchResolver.BatchResolvePlugins(missingPlugins)
	if err != nil {
		return err
	}

	return downloadBatchPlugins(missingPlugins, batchDownloadInfo, d.config.ShowProgress)
}
