package discovery

import (
	"errors"
	"fmt"

	"github.com/go-semantic-release/semantic-release/v2/pkg/config"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver/github"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver/registry"
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
	resolvers, err := loadResolvers(registry.NewResolver(), github.NewResolver())
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

func (d *Discovery) fetchPlugin(pluginInfo *plugin.PluginInfo) (string, error) {
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

func (d *Discovery) FindPlugin(t, name string) (*plugin.PluginInfo, error) {
	pInfo, err := plugin.GetPluginInfo(t, name)
	if err != nil {
		return nil, err
	}

	err = setAndEnsurePluginPath(pInfo)
	if err != nil {
		return nil, err
	}

	binPath, err := findPluginLocally(pInfo)
	if errors.Is(err, ErrPluginNotFound) {
		binPath, err = d.fetchPlugin(pInfo)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	pInfo.BinPath = binPath
	return pInfo, nil
}
