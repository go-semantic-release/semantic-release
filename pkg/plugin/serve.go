package plugin

import (
	"github.com/go-semantic-release/semantic-release/pkg/analyzer"
	"github.com/hashicorp/go-plugin"
)

var Handshake = plugin.HandshakeConfig{
	MagicCookieKey:   "GO_SEMANTIC_RELEASE_MAGIC_COOKIE",
	MagicCookieValue: "beepboop",
}

type CommitAnalyzerFunc func() analyzer.CommitAnalyzer

//type ChangelogGeneratorFunc func() changelog.Generator

type ServeOpts struct {
	CommitAnalyzer CommitAnalyzerFunc
	//ChangelogGenerator ChangelogGeneratorFunc
}

func Serve(opts *ServeOpts) {
	pluginSet := make(plugin.PluginSet)

	switch {
	case opts.CommitAnalyzer != nil:
		pluginSet[analyzer.CommitAnalyzerPluginName] = &GRPCWrapper{
			Type: analyzer.CommitAnalyzerPluginName,
			Impl: opts.CommitAnalyzer(),
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
