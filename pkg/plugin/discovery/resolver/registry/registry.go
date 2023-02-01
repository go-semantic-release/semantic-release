package registry

import (
	"fmt"

	"github.com/go-semantic-release/plugin-registry/pkg/client"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver"
)

const DefaultEndpoint = "https://registry-staging.go-semantic-release.xyz/api/v2"

type Resolver struct {
	client *client.Client
}

func NewResolver() *Resolver {
	return &Resolver{
		client: client.New(DefaultEndpoint),
	}
}

func (r *Resolver) ResolvePlugin(pluginInfo *plugin.Info) (*resolver.PluginDownloadInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *Resolver) Names() []string {
	// TODO: this should be registry when the registry is ready
	return []string{"registry-beta"}
}
