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
	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
	"github.com/go-semantic-release/semantic-release/v2/pkg/updater"
	"github.com/spf13/cobra"
)

type tempCommitAnalyzerInterface interface {
	Analyze(commits []*semrel.RawCommit) []*semrel.Commit
}

type tempCommitAnalyzerWrapper struct {
	tempCommitAnalyzerInterface
}

func (*tempCommitAnalyzerWrapper) Init(m map[string]string) error {
	return nil
}

func (*tempCommitAnalyzerWrapper) Name() string {
	return "default"
}

func (*tempCommitAnalyzerWrapper) Version() string {
	return "dev"
}

type tempCIInterface interface {
	Name() string
	RunCondition(m map[string]string) error
	GetCurrentBranch() string
	GetCurrentSHA() string
}

type tempCIWrapper struct {
	tempCIInterface
}

func (*tempCIWrapper) Version() string {
	return "dev"
}

type tempChangelogGeneratorInterface interface {
	Generate(config *generator.ChangelogGeneratorConfig) string
}

type tempChangelogGeneratorWrapper struct {
	tempChangelogGeneratorInterface
}

func (*tempChangelogGeneratorWrapper) Init(m map[string]string) error {
	return nil
}

func (*tempChangelogGeneratorWrapper) Name() string {
	return "default"
}

func (*tempChangelogGeneratorWrapper) Version() string {
	return "dev"
}

type tempProviderInterface interface {
	Init(map[string]string) error
	Name() string
	GetInfo() (*provider.RepositoryInfo, error)
	GetCommits(sha string) ([]*semrel.RawCommit, error)
	GetReleases(re string) ([]*semrel.Release, error)
	CreateRelease(*provider.CreateReleaseConfig) error
}

type tempProviderWrapper struct {
	tempProviderInterface
}

func (*tempProviderWrapper) Version() string {
	return "dev"
}

type tempFilesUpdaterInterface interface {
	ForFiles() string
	Apply(file, newVersion string) error
}

type tempFilesUpdaterWrapper struct {
	tempFilesUpdaterInterface
}

func (*tempFilesUpdaterWrapper) Init(m map[string]string) error {
	return nil
}

func (*tempFilesUpdaterWrapper) Name() string {
	return "default"
}

func (*tempFilesUpdaterWrapper) Version() string {
	return "dev"
}

func RegisterPluginCommands(cmd *cobra.Command) {
	cmd.AddCommand([]*cobra.Command{
		{
			Use: analyzer.CommitAnalyzerPluginName,
			Run: func(cmd *cobra.Command, args []string) {
				plugin.Serve(&plugin.ServeOpts{
					CommitAnalyzer: func() analyzer.CommitAnalyzer {
						return &tempCommitAnalyzerWrapper{&defaultAnalyzer.DefaultCommitAnalyzer{}}
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
						return &tempCIWrapper{&defaultCI.DefaultCI{}}
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
						return &tempCIWrapper{&githubCI.GitHubActions{}}
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
						return &tempCIWrapper{&gitlabCI.GitLab{}}
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
						return &tempChangelogGeneratorWrapper{&defaultGenerator.DefaultChangelogGenerator{}}
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
						return &tempProviderWrapper{&providerGithub.GitHubRepository{}}
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
						return &tempProviderWrapper{&providerGitlab.GitLabRepository{}}
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
						return &tempFilesUpdaterWrapper{&npmUpdater.Updater{}}
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
