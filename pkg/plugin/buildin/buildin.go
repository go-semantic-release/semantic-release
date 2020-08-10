package buildin

import (
	"os"
	"os/exec"
	"strings"

	"github.com/go-semantic-release/semantic-release/pkg/analyzer"
	caPlugin "github.com/go-semantic-release/semantic-release/pkg/analyzer/plugin"
	"github.com/go-semantic-release/semantic-release/pkg/condition"
	"github.com/go-semantic-release/semantic-release/pkg/condition/defaultci"
	githubCI "github.com/go-semantic-release/semantic-release/pkg/condition/github"
	gitlabCI "github.com/go-semantic-release/semantic-release/pkg/condition/gitlab"
	"github.com/go-semantic-release/semantic-release/pkg/condition/travis"
	"github.com/go-semantic-release/semantic-release/pkg/plugin"
	"github.com/urfave/cli/v2"
)

func GetPluginCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:     analyzer.CommitAnalyzerPluginName,
			Action:   caPlugin.Main,
			Hidden:   true,
			HideHelp: true,
		},
		{
			Name:     condition.CIConditionPluginName + "_default",
			Action:   defaultci.Main,
			Hidden:   true,
			HideHelp: true,
		},
		{
			Name:     condition.CIConditionPluginName + "_github",
			Action:   githubCI.Main,
			Hidden:   true,
			HideHelp: true,
		},
		{
			Name:     condition.CIConditionPluginName + "_gitlab",
			Action:   gitlabCI.Main,
			Hidden:   true,
			HideHelp: true,
		},
		{
			Name:     condition.CIConditionPluginName + "_travis",
			Action:   travis.Main,
			Hidden:   true,
			HideHelp: true,
		},
	}
}

func GetPluginOpts(t string, suffixNames ...string) *plugin.PluginOpts {
	bin, _ := os.Executable()
	commandName := t
	if len(suffixNames) > 0 {
		commandName = t + "_" + strings.Join(suffixNames, "_")
	}
	return &plugin.PluginOpts{
		Type: t,
		Cmd:  exec.Command(bin, commandName),
	}
}
