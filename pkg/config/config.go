package config

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config is a complete set of app configuration
type Config struct {
	Token                           string
	UpdateFiles                     []string
	ProviderPlugin                  string
	ProviderOpts                    map[string]string
	Changelog                       string
	Match                           string
	VersionFile                     bool
	Prerelease                      bool
	Ghr                             bool
	NoCI                            bool
	Dry                             bool
	AllowInitialDevelopmentVersions bool
	AllowNoChanges                  bool
	MaintainedVersion               string
}

func MustGetString(cmd *cobra.Command, name string) string {
	res, err := cmd.Flags().GetString(name)
	if err != nil {
		panic(err)
	}
	return res
}

func MustGetStringArray(cmd *cobra.Command, name string) []string {
	res, err := cmd.Flags().GetStringArray(name)
	if err != nil {
		panic(err)
	}
	return res
}

func MustGetBool(cmd *cobra.Command, name string) bool {
	res, err := cmd.Flags().GetBool(name)
	if err != nil {
		panic(err)
	}
	return res
}

// NewConfig returns a new Config instance
func NewConfig(cmd *cobra.Command) (*Config, error) {
	provOpts := make(map[string]string)
	for k, v := range viper.GetStringMapString("plugins.provider.options") {
		provOpts[k] = v
	}
	for _, opt := range MustGetStringArray(cmd, "provider-opt") {
		sOpt := strings.SplitN(opt, "=", 2)
		if len(sOpt) < 2 {
			continue
		}
		provOpts[strings.ToLower(sOpt[0])] = sOpt[1]
	}

	conf := &Config{
		Token:                           MustGetString(cmd, "token"),
		UpdateFiles:                     MustGetStringArray(cmd, "update"),
		ProviderPlugin:                  viper.GetString("plugins.provider.name"),
		ProviderOpts:                    provOpts,
		Changelog:                       MustGetString(cmd, "changelog"),
		Match:                           MustGetString(cmd, "match"),
		VersionFile:                     MustGetBool(cmd, "version-file"),
		Prerelease:                      MustGetBool(cmd, "prerelease"),
		Ghr:                             MustGetBool(cmd, "ghr"),
		NoCI:                            MustGetBool(cmd, "no-ci"),
		Dry:                             MustGetBool(cmd, "dry"),
		AllowInitialDevelopmentVersions: MustGetBool(cmd, "allow-initial-development-versions"),
		AllowNoChanges:                  MustGetBool(cmd, "allow-no-changes"),
		MaintainedVersion:               viper.GetString("maintainedVersion"),
	}

	return conf, nil
}
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func InitConfig(cmd *cobra.Command) error {
	cmd.Flags().StringP("token", "t", "", "provider token")
	cmd.Flags().StringArrayP("update", "u", []string{}, "updates the version of a certain files")
	cmd.Flags().StringP("provider", "p", "github", "provider token")
	cmd.Flags().StringArrayP("provider-opt", "o", []string{}, "options that are passed to the provider plugin")
	cmd.Flags().String("changelog", "", "creates a changelog file")
	cmd.Flags().String("match", "", "only consider tags matching the given glob(7) pattern, excluding the \"refs/tags/\" prefix.")
	cmd.Flags().String("maintained-version", "", "set the maintained version as base for new releases")
	cmd.Flags().BoolP("version-file", "f", false, "create a .version file with the new version")
	cmd.Flags().Bool("prerelease", false, "flags the release as a prerelease")
	cmd.Flags().Bool("ghr", false, "create a .ghr file with the parameters for ghr")
	cmd.Flags().Bool("no-ci", false, "run semantic-release locally")
	cmd.Flags().Bool("dry", false, "do not create release")
	cmd.Flags().Bool("allow-initial-development-versions", false, "semantic-release will start your initial development release at 0.1.0")
	cmd.Flags().Bool("allow-no-changes", false, "exit with code 0 if no changes are found, useful if semantic-release is automatically run")
	cmd.Flags().SortFlags = true
	Must(viper.BindPFlag("maintainedVersion", cmd.Flags().Lookup("maintained-version")))
	viper.AddConfigPath(".")
	viper.SetConfigName(".semrelrc")
	viper.SetConfigType("json")
	viper.SetDefault("plugins.commit-analyzer.name", "default")
	viper.SetDefault("plugins.ci-condition.name", "default")
	viper.SetDefault("plugins.changelog-generator.name", "default")
	Must(viper.BindEnv("maintainedVersion", "MAINTAINED_VERSION"))
	Must(viper.BindPFlag("plugins.provider.name", cmd.Flags().Lookup("provider")))
	viper.SetDefault("plugins.files-updater.name", "default")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return nil
}
