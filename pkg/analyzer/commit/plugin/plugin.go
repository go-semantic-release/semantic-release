package plugin

import (
	"github.com/go-semantic-release/semantic-release/pkg/analyzer/commit"
	"github.com/go-semantic-release/semantic-release/pkg/plugin"
	"github.com/urfave/cli/v2"
)

func DefaultMain(c *cli.Context) error {
	plugin.Serve(&plugin.ServeOpts{
		CommitAnalyzer: func() commit.Analyzer {
			return &commit.DefaultAnalyzer{}
		},
	})
	return nil
}
