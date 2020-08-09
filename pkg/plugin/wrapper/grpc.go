package wrapper

import (
	"context"
	"errors"

	"github.com/go-semantic-release/semantic-release/pkg/analyzer"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type GRPC struct {
	Type string
	Impl interface{}
	plugin.NetRPCUnsupportedPlugin
}

func (p *GRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	switch p.Type {
	case analyzer.PluginNameCommitAnalyzer:
		analyzer.RegisterCommitAnalyzerPluginServer(s, &analyzer.CommitAnalyzerServer{
			Impl: p.Impl.(analyzer.CommitAnalyzer),
		})
	}
	return nil
}

func (p *GRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	switch p.Type {
	case analyzer.PluginNameCommitAnalyzer:
		return &analyzer.CommitAnalyzerClient{
			Impl: analyzer.NewCommitAnalyzerPluginClient(c),
		}, nil
	}
	return nil, errors.New("unknown plugin type")
}