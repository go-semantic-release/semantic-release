package buildin

import (
	"os"
	"os/exec"

	"github.com/go-semantic-release/semantic-release/pkg/analyzer"
	caPlugin "github.com/go-semantic-release/semantic-release/pkg/analyzer/plugin"
	"github.com/go-semantic-release/semantic-release/pkg/plugin"
	"github.com/urfave/cli/v2"
)

func GetPluginCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:     analyzer.CommitAnalyzerPluginName,
			Action:   caPlugin.DefaultMain,
			Hidden:   true,
			HideHelp: true,
		},
	}
}

func GetPluginOpts(t string) *plugin.PluginOpts {
	bin, _ := os.Executable()
	return &plugin.PluginOpts{
		Type: analyzer.CommitAnalyzerPluginName,
		Cmd:  exec.Command(bin, t),
	}
}
