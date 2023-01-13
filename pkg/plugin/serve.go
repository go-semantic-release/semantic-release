package plugin

import (
	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/condition"
	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/updater"
	"github.com/hashicorp/go-plugin"
)

var Handshake = plugin.HandshakeConfig{
	MagicCookieKey:   "GO_SEMANTIC_RELEASE_MAGIC_COOKIE",
	MagicCookieValue: "beepboop",
}

type (
	CommitAnalyzerFunc     func() analyzer.CommitAnalyzer
	CIConditionFunc        func() condition.CICondition
	ChangelogGeneratorFunc func() generator.ChangelogGenerator
	ProviderFunc           func() provider.Provider
	FilesUpdaterFunc       func() updater.FilesUpdater
	HooksFunc              func() hooks.Hooks
)

type ServeOpts struct {
	CommitAnalyzer     CommitAnalyzerFunc
	CICondition        CIConditionFunc
	ChangelogGenerator ChangelogGeneratorFunc
	Provider           ProviderFunc
	FilesUpdater       FilesUpdaterFunc
	Hooks              HooksFunc
}

func Serve(opts *ServeOpts) {
	pluginSet := make(plugin.PluginSet)

	if opts.CommitAnalyzer != nil {
		pluginSet[analyzer.CommitAnalyzerPluginName] = &GRPCWrapper{
			Type: analyzer.CommitAnalyzerPluginName,
			Impl: opts.CommitAnalyzer(),
		}
	}

	if opts.CICondition != nil {
		pluginSet[condition.CIConditionPluginName] = &GRPCWrapper{
			Type: condition.CIConditionPluginName,
			Impl: opts.CICondition(),
		}
	}

	if opts.ChangelogGenerator != nil {
		pluginSet[generator.ChangelogGeneratorPluginName] = &GRPCWrapper{
			Type: generator.ChangelogGeneratorPluginName,
			Impl: opts.ChangelogGenerator(),
		}
	}

	if opts.Provider != nil {
		pluginSet[provider.PluginName] = &GRPCWrapper{
			Type: provider.PluginName,
			Impl: opts.Provider(),
		}
	}

	if opts.FilesUpdater != nil {
		pluginSet[updater.FilesUpdaterPluginName] = &GRPCWrapper{
			Type: updater.FilesUpdaterPluginName,
			Impl: opts.FilesUpdater(),
		}
	}

	if opts.Hooks != nil {
		pluginSet[hooks.PluginName] = &GRPCWrapper{
			Type: hooks.PluginName,
			Impl: opts.Hooks(),
		}
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: pluginSet,
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
