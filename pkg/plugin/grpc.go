package plugin

import (
	"context"
	"errors"

	"github.com/go-semantic-release/semantic-release/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/pkg/condition"
	"github.com/go-semantic-release/semantic-release/pkg/generator"
	"github.com/go-semantic-release/semantic-release/pkg/provider"
	"github.com/go-semantic-release/semantic-release/pkg/updater"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type GRPCWrapper struct {
	Type string
	Impl interface{}
	plugin.NetRPCUnsupportedPlugin
}

func (p *GRPCWrapper) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	switch p.Type {
	case analyzer.CommitAnalyzerPluginName:
		analyzer.RegisterCommitAnalyzerPluginServer(s, &analyzer.CommitAnalyzerServer{
			Impl: p.Impl.(analyzer.CommitAnalyzer),
		})
	case condition.CIConditionPluginName:
		condition.RegisterCIConditionPluginServer(s, &condition.CIConditionServer{
			Impl: p.Impl.(condition.CICondition),
		})
	case generator.ChangelogGeneratorPluginName:
		generator.RegisterChangelogGeneratorPluginServer(s, &generator.ChangelogGeneratorServer{
			Impl: p.Impl.(generator.ChangelogGenerator),
		})
	case provider.PluginName:
		provider.RegisterProviderPluginServer(s, &provider.Server{
			Impl: p.Impl.(provider.Provider),
		})
	case updater.FilesUpdaterPluginName:
		updater.RegisterFilesUpdaterPluginServer(s, &updater.FilesUpdaterServer{
			Impl: p.Impl.(updater.FilesUpdater),
		})
	}

	return nil
}

func (p *GRPCWrapper) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	switch p.Type {
	case analyzer.CommitAnalyzerPluginName:
		return &analyzer.CommitAnalyzerClient{
			Impl: analyzer.NewCommitAnalyzerPluginClient(c),
		}, nil
	case condition.CIConditionPluginName:
		return &condition.CIConditionClient{
			Impl: condition.NewCIConditionPluginClient(c),
		}, nil
	case generator.ChangelogGeneratorPluginName:
		return &generator.ChangelogGeneratorClient{
			Impl: generator.NewChangelogGeneratorPluginClient(c),
		}, nil
	case provider.PluginName:
		return &provider.Client{
			Impl: provider.NewProviderPluginClient(c),
		}, nil
	case updater.FilesUpdaterPluginName:
		return &updater.FilesUpdaterClient{
			Impl: updater.NewFilesUpdaterPluginClient(c),
		}, nil
	}
	return nil, errors.New("unknown plugin type")
}
