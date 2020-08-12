package buildin

import (
	"os"
	"os/exec"
	"strings"

	defaultGenerator "github.com/go-semantic-release/changelog-generator-default/pkg/generator"
	defaultAnalyzer "github.com/go-semantic-release/commit-analyzer-cz/pkg/analyzer"
	defaultCI "github.com/go-semantic-release/condition-default/pkg/condition"
	githubCI "github.com/go-semantic-release/condition-github/pkg/condition"
	gitlabCI "github.com/go-semantic-release/condition-gitlab/pkg/condition"
	npmUpdater "github.com/go-semantic-release/files-updater-npm/pkg/updater"
	providerGithub "github.com/go-semantic-release/provider-github/pkg/provider"
	providerGitlab "github.com/go-semantic-release/provider-gitlab/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/condition"
	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/updater"
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
			Name: condition.CIConditionPluginName + "_default",
			Action: func(c *cli.Context) error {
				plugin.Serve(&plugin.ServeOpts{
					CICondition: func() condition.CICondition {
						return &defaultCI.DefaultCI{}
					},
				})
				return nil
			},
			Hidden:   true,
			HideHelp: true,
		},
		{
			Name: condition.CIConditionPluginName + "_github",
			Action: func(c *cli.Context) error {
				plugin.Serve(&plugin.ServeOpts{
					CICondition: func() condition.CICondition {
						return &githubCI.GitHubActions{}
					},
				})
				return nil
			},
			Hidden:   true,
			HideHelp: true,
		},
		{
			Name: condition.CIConditionPluginName + "_gitlab",
			Action: func(c *cli.Context) error {
				plugin.Serve(&plugin.ServeOpts{
					CICondition: func() condition.CICondition {
						return &gitlabCI.GitLab{}
					},
				})
				return nil
			},
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
			Name: updater.FilesUpdaterPluginName + "_npm",
			Action: func(c *cli.Context) error {
				plugin.Serve(&plugin.ServeOpts{
					FilesUpdater: func() updater.FilesUpdater {
						return &npmUpdater.Updater{}
					},
				})
				return nil
			},
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
