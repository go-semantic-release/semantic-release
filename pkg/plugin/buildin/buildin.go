package buildin

import (
	caPlugin "github.com/go-semantic-release/semantic-release/pkg/analyzer/plugin"
	"github.com/urfave/cli/v2"
)

func GetPluginCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:     "commit-analyzer-plugin",
			Action:   caPlugin.DefaultMain,
			Hidden:   true,
			HideHelp: true,
		},
	}
}
