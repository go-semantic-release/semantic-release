package config

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config is a complete set of app configuration
type Config struct {
	Token                           string
	ProviderPlugin                  string
	ProviderOpts                    map[string]string
	CommitAnalyzerPlugin            string
	CommitAnalyzerOpts              map[string]string
	CIConditionPlugin               string
	CIConditionOpts                 map[string]string
	ChangelogGeneratorPlugin        string
	ChangelogGeneratorOpts          map[string]string
	Changelog                       string
	FilesUpdaterPlugins             []string
	FilesUpdaterOpts                map[string]string
	HooksPlugins                    []string
	HooksOpts                       map[string]string
	UpdateFiles                     []string
	Match                           string
	VersionFile                     bool
	Prerelease                      bool
	Ghr                             bool
	NoCI                            bool
	Dry                             bool
	AllowInitialDevelopmentVersions bool
	AllowNoChanges                  bool
	ForceBumpPatchVersion           bool
	MaintainedVersion               string
	PrependChangelog                bool
	DownloadPlugins                 bool
	ShowProgress                    bool
}

func mustGetString(cmd *cobra.Command, name string) string {
	res, err := cmd.Flags().GetString(name)
	if err != nil {
		panic(err)
	}
	return res
}

func mustGetStringArray(cmd *cobra.Command, name string) []string {
	res, err := cmd.Flags().GetStringArray(name)
	if err != nil {
		panic(err)
	}
	return res
}

func mustGetBool(cmd *cobra.Command, name string) bool {
	res, err := cmd.Flags().GetBool(name)
	if err != nil {
		panic(err)
	}
	return res
}

func mergeOpts(v map[string]string, c []string) map[string]string {
	opts := make(map[string]string)
	for k, v := range v {
		opts[k] = v
	}
	for _, opt := range c {
		sOpt := strings.SplitN(opt, "=", 2)
		if len(sOpt) < 2 {
			continue
		}
		opts[strings.ToLower(sOpt[0])] = sOpt[1]
	}
	return opts
}

func NewConfig(cmd *cobra.Command) (*Config, error) {
	provOpts := mergeOpts(
		viper.GetStringMapString("plugins.provider.options"),
		mustGetStringArray(cmd, "provider-opt"))
	caOpts := mergeOpts(
		viper.GetStringMapString("plugins.commit-analyzer.options"),
		mustGetStringArray(cmd, "commit-analyzer-opt"))
	ciOpts := mergeOpts(
		viper.GetStringMapString("plugins.ci-condition.options"),
		mustGetStringArray(cmd, "ci-condition-opt"))
	cgOpts := mergeOpts(
		viper.GetStringMapString("plugins.changelog-generator.options"),
		mustGetStringArray(cmd, "changelog-generator-opt"))
	fuOpts := mergeOpts(
		viper.GetStringMapString("plugins.files-updater.options"),
		mustGetStringArray(cmd, "files-updater-opt"))
	hoOpts := mergeOpts(
		viper.GetStringMapString("plugins.hooks.options"),
		mustGetStringArray(cmd, "hooks-opt"))

	conf := &Config{
		Token:                           mustGetString(cmd, "token"),
		ProviderPlugin:                  viper.GetString("plugins.provider.name"),
		ProviderOpts:                    provOpts,
		CommitAnalyzerPlugin:            viper.GetString("plugins.commit-analyzer.name"),
		CommitAnalyzerOpts:              caOpts,
		CIConditionPlugin:               viper.GetString("plugins.ci-condition.name"),
		CIConditionOpts:                 ciOpts,
		ChangelogGeneratorPlugin:        viper.GetString("plugins.changelog-generator.name"),
		ChangelogGeneratorOpts:          cgOpts,
		Changelog:                       mustGetString(cmd, "changelog"),
		FilesUpdaterPlugins:             viper.GetStringSlice("plugins.files-updater.names"),
		FilesUpdaterOpts:                fuOpts,
		HooksPlugins:                    viper.GetStringSlice("plugins.hooks.names"),
		HooksOpts:                       hoOpts,
		UpdateFiles:                     mustGetStringArray(cmd, "update"),
		Match:                           mustGetString(cmd, "match"),
		VersionFile:                     mustGetBool(cmd, "version-file"),
		Prerelease:                      mustGetBool(cmd, "prerelease"),
		Ghr:                             mustGetBool(cmd, "ghr"),
		NoCI:                            mustGetBool(cmd, "no-ci"),
		Dry:                             mustGetBool(cmd, "dry"),
		AllowInitialDevelopmentVersions: mustGetBool(cmd, "allow-initial-development-versions"),
		AllowNoChanges:                  mustGetBool(cmd, "allow-no-changes"),
		ForceBumpPatchVersion:           mustGetBool(cmd, "force-bump-patch-version"),
		MaintainedVersion:               viper.GetString("maintainedVersion"),
		PrependChangelog:                mustGetBool(cmd, "prepend-changelog"),
		DownloadPlugins:                 mustGetBool(cmd, "download-plugins"),
		ShowProgress:                    mustGetBool(cmd, "show-progress"),
	}
	return conf, nil
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func detectCI() string {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return "github"
	}
	if os.Getenv("GITLAB_CI") == "true" {
		return "gitlab"
	}
	return "default"
}

func defaultProvider() string {
	if os.Getenv("GITLAB_CI") == "true" {
		return "gitlab"
	}
	return "github"
}

func InitConfig(cmd *cobra.Command) error {
	cmd.Flags().StringP("token", "t", "", "provider token")
	cmd.Flags().String("provider", defaultProvider(), "provider plugin name")
	cmd.Flags().StringArray("provider-opt", []string{}, "options that are passed to the provider plugin")
	cmd.Flags().String("commit-analyzer", "default", "commit-analyzer plugin name")
	cmd.Flags().StringArray("commit-analyzer-opt", []string{}, "options that are passed to the commit-analyzer plugin")
	cmd.Flags().String("ci-condition", detectCI(), "ci-condition plugin name")
	cmd.Flags().StringArray("ci-condition-opt", []string{}, "options that are passed to the ci-condition plugin")
	cmd.Flags().String("changelog-generator", "default", "changelog-generator plugin name")
	cmd.Flags().StringArray("changelog-generator-opt", []string{}, "options that are passed to the changelog-generator plugin")
	cmd.Flags().String("changelog", "", "creates a changelog file")
	cmd.Flags().StringSlice("files-updater", []string{}, "files-updater plugin names")
	cmd.Flags().StringArray("files-updater-opt", []string{}, "options that are passed to the files-updater plugins")
	cmd.Flags().StringSlice("hooks", []string{}, "hooks plugin names")
	cmd.Flags().StringArray("hooks-opt", []string{}, "options that are passed to the hooks plugins")
	cmd.Flags().StringArrayP("update", "u", []string{}, "updates the version of a certain files")
	cmd.Flags().String("match", "", "only consider tags matching the given glob(7) pattern, excluding the \"refs/tags/\" prefix.")
	cmd.Flags().String("maintained-version", "", "set the maintained version as base for new releases")
	cmd.Flags().BoolP("version-file", "f", false, "create a .version file with the new version")
	cmd.Flags().Bool("prerelease", false, "flags the release as a prerelease")
	cmd.Flags().Bool("ghr", false, "create a .ghr file with the parameters for ghr")
	cmd.Flags().Bool("no-ci", false, "run semantic-release locally")
	cmd.Flags().Bool("dry", false, "do not create release")
	cmd.Flags().Bool("allow-initial-development-versions", false, "semantic-release will start your initial development release at 0.1.0")
	cmd.Flags().Bool("allow-no-changes", false, "exit with code 0 if no changes are found, useful if semantic-release is automatically run")
	cmd.Flags().Bool("force-bump-patch-version", false, "increments the patch version if no changes are found")
	cmd.Flags().Bool("prepend-changelog", false, "if the changelog file already exist the new changelog is prepended")
	cmd.Flags().Bool("download-plugins", false, "downloads all required plugins if needed")
	cmd.Flags().Bool("show-progress", false, "shows the plugin download progress")
	cmd.Flags().SortFlags = true

	viper.AddConfigPath(".")
	viper.SetConfigName(".semrelrc")
	viper.SetConfigType("json")

	must(viper.BindPFlag("maintainedVersion", cmd.Flags().Lookup("maintained-version")))
	must(viper.BindEnv("maintainedVersion", "MAINTAINED_VERSION"))

	must(viper.BindPFlag("plugins.provider.name", cmd.Flags().Lookup("provider")))
	must(viper.BindPFlag("plugins.commit-analyzer.name", cmd.Flags().Lookup("commit-analyzer")))
	must(viper.BindPFlag("plugins.ci-condition.name", cmd.Flags().Lookup("ci-condition")))
	must(viper.BindPFlag("plugins.changelog-generator.name", cmd.Flags().Lookup("changelog-generator")))
	must(viper.BindPFlag("plugins.files-updater.names", cmd.Flags().Lookup("files-updater")))
	must(viper.BindPFlag("plugins.hooks.names", cmd.Flags().Lookup("hooks")))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return nil
}
