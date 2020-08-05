package plugin

import (
	"github.com/go-semantic-release/semantic-release/pkg/analyzer/commit"
	"github.com/go-semantic-release/semantic-release/pkg/plugin/wrapper"
	"github.com/hashicorp/go-plugin"
)

var Handshake = plugin.HandshakeConfig{
	MagicCookieKey:   "GO_SEMANTIC_RELEASE_MAGIC_COOKIE",
	MagicCookieValue: "beepboop",
}

type CommitAnalyzerFunc func() commit.Analyzer

//type ChangelogGeneratorFunc func() changelog.Generator

type ServeOpts struct {
	CommitAnalyzer CommitAnalyzerFunc
	//ChangelogGenerator ChangelogGeneratorFunc
}

func Serve(opts *ServeOpts) {
	pluginSet := make(plugin.PluginSet)

	switch {
	case opts.CommitAnalyzer != nil:
		pluginSet[wrapper.TypeCommitAnalyzer] = &wrapper.GRPC{
			Type: wrapper.TypeCommitAnalyzer,
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
