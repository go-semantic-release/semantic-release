package plugin

import (
	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/urfave/cli/v2"
)

func Main(c *cli.Context) error {
	plugin.Serve(&plugin.ServeOpts{
		CommitAnalyzer: func() analyzer.CommitAnalyzer {
			return &DefaultCommitAnalyzer{}
		},
	})
	return nil
}
