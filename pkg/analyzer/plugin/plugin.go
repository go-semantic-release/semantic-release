package plugin

import (
	defaultAnalyzer "github.com/go-semantic-release/commit-analyzer-cz/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/urfave/cli/v2"
)

func Main(c *cli.Context) error {
	plugin.Serve(&plugin.ServeOpts{
		CommitAnalyzer: func() analyzer.CommitAnalyzer {
			return &defaultAnalyzer.DefaultCommitAnalyzer{}
		},
	})
	return nil
}
