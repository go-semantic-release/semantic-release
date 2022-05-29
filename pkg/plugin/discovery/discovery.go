package discovery

import (
	"errors"
	"fmt"

	"github.com/go-semantic-release/semantic-release/v2/pkg/config"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver/registry"
)

type Discovery struct {
	config    *config.Config
	resolvers map[string]resolver.Resolver
}

func New(config *config.Config) (*Discovery, error) {
	registryResolver := registry.NewResolver()
	return &Discovery{
		config: config,
		resolvers: map[string]resolver.Resolver{
			"default":  registryResolver,
			"registry": registryResolver,
		},
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
	if err := setAndEnsurePluginPath(pInfo); err != nil {
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
