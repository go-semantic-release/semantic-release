package buildin

import (
	"os"
	"os/exec"
	"strings"

	defaultGenerator "github.com/go-semantic-release/changelog-generator-default/pkg/generator"
	defaultAnalyzer "github.com/go-semantic-release/commit-analyzer-cz/pkg/analyzer"
	providerGithub "github.com/go-semantic-release/provider-github/pkg/provider"
	providerGitlab "github.com/go-semantic-release/provider-gitlab/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/condition"
	"github.com/go-semantic-release/semantic-release/v2/pkg/condition/defaultci"
	githubCI "github.com/go-semantic-release/semantic-release/v2/pkg/condition/github"
	gitlabCI "github.com/go-semantic-release/semantic-release/v2/pkg/condition/gitlab"
	"github.com/go-semantic-release/semantic-release/v2/pkg/condition/travis"
	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/updater"
	"github.com/go-semantic-release/semantic-release/v2/pkg/updater/npm"
	"github.com/urfave/cli/v2"
)

func GetPluginCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: analyzer.CommitAnalyzerPluginName,
			Action: func(c *cli.Context) error {
				plugin.Serve(&plugin.ServeOpts{
					CommitAnalyzer: func() analyzer.CommitAnalyzer {
						return &defaultAnalyzer.DefaultCommitAnalyzer{}
					},
				})
				return nil
			},
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
		{
			Name: generator.ChangelogGeneratorPluginName,
			Action: func(c *cli.Context) error {
				plugin.Serve(&plugin.ServeOpts{
					ChangelogGenerator: func() generator.ChangelogGenerator {
						return &defaultGenerator.DefaultChangelogGenerator{}
					},
				})
				return nil
			},
			Hidden:   true,
			HideHelp: true,
		},
		{
			Name: provider.PluginName + "_github",
			Action: func(c *cli.Context) error {
				plugin.Serve(&plugin.ServeOpts{
					Provider: func() provider.Provider {
						return &providerGithub.GitHubRepository{}
					},
				})
				return nil
			},
			Hidden:   true,
			HideHelp: true,
		},
		{
			Name: provider.PluginName + "_gitlab",
			Action: func(c *cli.Context) error {
				plugin.Serve(&plugin.ServeOpts{
					Provider: func() provider.Provider {
						return &providerGitlab.GitLabRepository{}
					},
				})
				return nil
			},
			Hidden:   true,
			HideHelp: true,
		},
		{
			Name:     updater.FilesUpdaterPluginName + "_npm",
			Action:   npm.Main,
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
