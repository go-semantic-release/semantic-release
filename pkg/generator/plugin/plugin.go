package plugin

import (
	"github.com/go-semantic-release/semantic-release/pkg/plugin"
	"github.com/urfave/cli/v2"
)

func Main(c *cli.Context) error {
	plugin.Serve(&plugin.ServeOpts{})
	return nil
}
