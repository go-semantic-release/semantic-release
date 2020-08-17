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
	"github.com/spf13/cobra"
)

func RegisterPluginCommands(cmd *cobra.Command) {
	cmd.AddCommand([]*cobra.Command{
		{
			Use: analyzer.CommitAnalyzerPluginName,
			Run: func(cmd *cobra.Command, args []string) {
				plugin.Serve(&plugin.ServeOpts{
					CommitAnalyzer: func() analyzer.CommitAnalyzer {
						return &defaultAnalyzer.DefaultCommitAnalyzer{}
					},
				})
			},
			Hidden: true,
		},
		{
			Use: condition.CIConditionPluginName + "_default",
			Run: func(cmd *cobra.Command, args []string) {
				plugin.Serve(&plugin.ServeOpts{
					CICondition: func() condition.CICondition {
						return &defaultCI.DefaultCI{}
					},
				})
			},
			Hidden: true,
		},
		{
			Use: condition.CIConditionPluginName + "_github",
			Run: func(cmd *cobra.Command, args []string) {
				plugin.Serve(&plugin.ServeOpts{
					CICondition: func() condition.CICondition {
						return &githubCI.GitHubActions{}
					},
				})
			},
			Hidden: true,
		},
		{
			Use: condition.CIConditionPluginName + "_gitlab",
			Run: func(cmd *cobra.Command, args []string) {
				plugin.Serve(&plugin.ServeOpts{
					CICondition: func() condition.CICondition {
						return &gitlabCI.GitLab{}
					},
				})
			},
			Hidden: true,
		},
		{
			Use: generator.ChangelogGeneratorPluginName,
			Run: func(cmd *cobra.Command, args []string) {
				plugin.Serve(&plugin.ServeOpts{
					ChangelogGenerator: func() generator.ChangelogGenerator {
						return &defaultGenerator.DefaultChangelogGenerator{}
					},
				})
			},
			Hidden: true,
		},
		{
			Use: provider.PluginName + "_github",
			Run: func(cmd *cobra.Command, args []string) {
				plugin.Serve(&plugin.ServeOpts{
					Provider: func() provider.Provider {
						return &providerGithub.GitHubRepository{}
					},
				})
			},
			Hidden: true,
		},
		{
			Use: provider.PluginName + "_gitlab",
			Run: func(cmd *cobra.Command, args []string) {
				plugin.Serve(&plugin.ServeOpts{
					Provider: func() provider.Provider {
						return &providerGitlab.GitLabRepository{}
					},
				})
			},
			Hidden: true,
		},
		{
			Use: updater.FilesUpdaterPluginName + "_npm",
			Run: func(cmd *cobra.Command, args []string) {
				plugin.Serve(&plugin.ServeOpts{
					FilesUpdater: func() updater.FilesUpdater {
						return &npmUpdater.Updater{}
					},
				})
			},
			Hidden: true,
		},
	}...)
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
