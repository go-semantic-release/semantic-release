package plugin

import (
	"os/exec"
	"sync"

	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/condition"
	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/updater"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type PluginOpts struct {
	Type string
	Cmd  *exec.Cmd
}

var runningClientsMx sync.Mutex
var runningClients = make([]*plugin.Client, 0)

func KillAllPlugins() {
	runningClientsMx.Lock()
	defer runningClientsMx.Unlock()
	for _, c := range runningClients {
		c.Kill()
	}
}

func startPlugin(opts *PluginOpts) (interface{}, error) {
	runningClientsMx.Lock()
	defer runningClientsMx.Unlock()
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: {
				opts.Type: &GRPCWrapper{
					Type: opts.Type,
				},
			},
		},
		Cmd:              opts.Cmd,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           hclog.NewNullLogger(),
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, err
	}
	raw, err := rpcClient.Dispense(opts.Type)
	if err != nil {
		client.Kill()
		return nil, err
	}
	runningClients = append(runningClients, client)
	return raw, nil
}

func StartCommitAnalyzerPlugin(opts *PluginOpts) (analyzer.CommitAnalyzer, error) {
	raw, err := startPlugin(opts)
	if err != nil {
		return nil, err
	}
	return raw.(analyzer.CommitAnalyzer), nil
}

func StartCIConditionPlugin(opts *PluginOpts) (condition.CICondition, error) {
	raw, err := startPlugin(opts)
	if err != nil {
		return nil, err
	}
	return raw.(condition.CICondition), nil
}

func StartChangelogGeneratorPlugin(opts *PluginOpts) (generator.ChangelogGenerator, error) {
	raw, err := startPlugin(opts)
	if err != nil {
		return nil, err
	}
	return raw.(generator.ChangelogGenerator), nil
}

func StartProviderPlugin(opts *PluginOpts) (provider.Provider, error) {
	raw, err := startPlugin(opts)
	if err != nil {
		return nil, err
	}
	return raw.(provider.Provider), nil
}

func StartFilesUpdaterPlugin(opts *PluginOpts) (updater.FilesUpdater, error) {
	raw, err := startPlugin(opts)
	if err != nil {
		return nil, err
	}
	return raw.(updater.FilesUpdater), nil
}
