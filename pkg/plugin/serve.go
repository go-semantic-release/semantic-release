package plugin

import (
	"github.com/go-semantic-release/semantic-release/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/pkg/condition"
	"github.com/hashicorp/go-plugin"
)

var Handshake = plugin.HandshakeConfig{
	MagicCookieKey:   "GO_SEMANTIC_RELEASE_MAGIC_COOKIE",
	MagicCookieValue: "beepboop",
}

type CommitAnalyzerFunc func() analyzer.CommitAnalyzer
type CIConditionFunc func() condition.CICondition

type ServeOpts struct {
	CommitAnalyzer CommitAnalyzerFunc
	CICondition    CIConditionFunc
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

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: pluginSet,
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
