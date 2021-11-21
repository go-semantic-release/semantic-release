package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Masterminds/semver/v3"
	"github.com/go-semantic-release/semantic-release/v2/pkg/config"
	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/manager"
	"github.com/go-semantic-release/semantic-release/v2/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
	"github.com/spf13/cobra"
)

// SRVERSION is the semantic-release version (added at compile time)
var SRVERSION string

var exitHandler func()

func errorHandler(logger *log.Logger) func(error, ...int) {
	return func(err error, exitCode ...int) {
		if err != nil {
			logger.Println(err)
			if exitHandler != nil {
				exitHandler()
			}
			if len(exitCode) == 1 {
				os.Exit(exitCode[0])
				return
			}
			os.Exit(1)
		}
	}
}

func main() {
	cmd := &cobra.Command{
		Use:     "semantic-release",
		Short:   "semantic-release - fully automated package/module/image publishing",
		Run:     cliHandler,
		Version: SRVERSION,
	}

	config.SetFlags(cmd)
	cobra.OnInitialize(func() {
		err := config.InitConfig(cmd)
		if err != nil {
			fmt.Printf("\nConfig error: %s\n", err.Error())
			os.Exit(1)
			return
		}
	})
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(1)
	}
}

func cliHandler(cmd *cobra.Command, args []string) {
	logger := log.New(os.Stderr, "[go-semantic-release]: ", 0)
	exitIfError := errorHandler(logger)

	logger.Printf("version: %s\n", SRVERSION)

	conf, err := config.NewConfig(cmd)
	exitIfError(err)

	pluginManager, err := manager.New(conf)
	exitIfError(err)
	exitHandler = func() {
		logger.Println("stopping plugins...")
		pluginManager.Stop()
	}
	defer exitHandler()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		exitIfError(errors.New("terminating..."))
	}()

	if conf.DownloadPlugins {
		exitIfError(pluginManager.FetchAllPlugins())
		logger.Println("all plugins are downloaded")
		os.Exit(0)
	}

	ci, err := pluginManager.GetCICondition()
	exitIfError(err)
	logger.Printf("ci-condition plugin: %s@%s\n", ci.Name(), ci.Version())

	prov, err := pluginManager.GetProvider()
	exitIfError(err)
	logger.Printf("provider plugin: %s@%s\n", prov.Name(), prov.Version())

	if conf.ProviderOpts["token"] == "" {
		conf.ProviderOpts["token"] = conf.Token
	}
	err = prov.Init(conf.ProviderOpts)
	exitIfError(err)

	logger.Println("getting default branch...")
	repoInfo, err := prov.GetInfo()
	exitIfError(err)
	logger.Println("found default branch: " + repoInfo.DefaultBranch)
	if repoInfo.Private {
		logger.Println("repo is private")
	}

	currentBranch := ci.GetCurrentBranch()
	if currentBranch == "" {
		exitIfError(fmt.Errorf("current branch not found"))
	}
	logger.Println("found current branch: " + currentBranch)

	if conf.MaintainedVersion != "" && currentBranch == repoInfo.DefaultBranch {
		exitIfError(fmt.Errorf("maintained version not allowed on default branch"))
	}

	if conf.MaintainedVersion != "" {
		logger.Println("found maintained version: " + conf.MaintainedVersion)
		repoInfo.DefaultBranch = "*"
	}

	currentSha := ci.GetCurrentSHA()
	logger.Println("found current sha: " + currentSha)

	hooksExecutor, err := pluginManager.GetChainedHooksExecutor()
	exitIfError(err)

	hooksNames := hooksExecutor.GetNameVersionPairs()
	if len(hooksNames) > 0 {
		logger.Printf("hooks plugins: %s\n", strings.Join(hooksNames, ", "))
	}

	exitIfError(hooksExecutor.Init(conf.HooksOpts))

	if !conf.NoCI {
		logger.Println("running CI condition...")
		conditionConfig := map[string]string{
			"token":         conf.Token,
			"defaultBranch": repoInfo.DefaultBranch,
			"private":       fmt.Sprintf("%t", repoInfo.Private),
		}
		for k, v := range conf.CIConditionOpts {
			conditionConfig[k] = v
		}
		err = ci.RunCondition(conditionConfig)
		if err != nil {
			herr := hooksExecutor.NoRelease(&hooks.NoReleaseConfig{
				Reason:  hooks.NoReleaseReason_CONDITION,
				Message: err.Error(),
			})
			if herr != nil {
				logger.Printf("there was an error executing the hooks plugins: %s", herr.Error())
			}
			exitIfError(err, 66)
		}

	}

	logger.Println("getting latest release...")
	matchRegex := ""
	match := strings.TrimSpace(conf.Match)
	if match != "" {
		logger.Printf("getting latest release matching %s...", match)
		matchRegex = "^" + match
	}
	releases, err := prov.GetReleases(matchRegex)
	exitIfError(err)
	release, err := semrel.GetLatestReleaseFromReleases(releases, conf.MaintainedVersion)
	exitIfError(err)
	logger.Println("found version: " + release.Version)

	if strings.Contains(conf.MaintainedVersion, "-") && semver.MustParse(release.Version).Prerelease() == "" {
		exitIfError(fmt.Errorf("no pre-release for this version possible"))
	}

	logger.Println("getting commits...")
	rawCommits, err := prov.GetCommits(release.SHA, currentSha)
	exitIfError(err)

	logger.Println("analyzing commits...")
	commitAnalyzer, err := pluginManager.GetCommitAnalyzer()
	exitIfError(err)
	logger.Printf("commit-analyzer plugin: %s@%s\n", commitAnalyzer.Name(), commitAnalyzer.Version())
	exitIfError(commitAnalyzer.Init(conf.ChangelogGeneratorOpts))

	commits := commitAnalyzer.Analyze(rawCommits)

	logger.Println("calculating new version...")
	newVer := semrel.GetNewVersion(conf, commits, release)
	if newVer == "" {
		herr := hooksExecutor.NoRelease(&hooks.NoReleaseConfig{
			Reason:  hooks.NoReleaseReason_NO_CHANGE,
			Message: "",
		})
		if herr != nil {
			logger.Printf("there was an error executing the hooks plugins: %s", herr.Error())
		}
		errNoChange := errors.New("no change")
		if conf.AllowNoChanges {
			exitIfError(errNoChange, 0)
		} else {
			exitIfError(errNoChange, 65)
		}
	}
	logger.Println("new version: " + newVer)

	logger.Println("generating changelog...")
	changelogGenerator, err := pluginManager.GetChangelogGenerator()
	exitIfError(err)
	logger.Printf("changelog-generator plugin: %s@%s\n", changelogGenerator.Name(), changelogGenerator.Version())
	exitIfError(changelogGenerator.Init(conf.ChangelogGeneratorOpts))

	changelogRes := changelogGenerator.Generate(&generator.ChangelogGeneratorConfig{
		Commits:       commits,
		LatestRelease: release,
		NewVersion:    newVer,
	})
	if conf.Changelog != "" {
		oldFile := make([]byte, 0)
		if conf.PrependChangelog {
			oldFileData, err := ioutil.ReadFile(conf.Changelog)
			if err == nil {
				oldFile = append([]byte("\n"), oldFileData...)
			}
		}
		changelogData := append([]byte(changelogRes), oldFile...)
		exitIfError(ioutil.WriteFile(conf.Changelog, changelogData, 0644))
	}

	if conf.Dry {
		if conf.VersionFile {
			exitIfError(ioutil.WriteFile(".version-unreleased", []byte(newVer), 0644))
		}
		exitIfError(errors.New("DRY RUN: no release was created"), 0)
	}

	logger.Println("creating release...")
	newRelease := &provider.CreateReleaseConfig{
		Changelog:  changelogRes,
		NewVersion: newVer,
		Prerelease: conf.Prerelease,
		Branch:     currentBranch,
		SHA:        currentSha,
	}
	exitIfError(prov.CreateRelease(newRelease))

	if conf.Ghr {
		exitIfError(ioutil.WriteFile(".ghr", []byte(fmt.Sprintf("-u %s -r %s v%s", repoInfo.Owner, repoInfo.Repo, newVer)), 0644))
	}

	if conf.VersionFile {
		exitIfError(ioutil.WriteFile(".version", []byte(newVer), 0644))
	}

	if len(conf.UpdateFiles) == 0 && len(conf.FilesUpdaterPlugins) > 0 {
		logger.Println("warning: file update plugins found but no files marked for update. You may be missing the update flag, e.g. --update package.json")
	}

	if len(conf.UpdateFiles) > 0 {
		logger.Println("updating files...")
		updater, err := pluginManager.GetChainedUpdater()
		exitIfError(err)
		logger.Printf("files-updater plugins: %s\n", strings.Join(updater.GetNameVersionPairs(), ", "))
		exitIfError(updater.Init(conf.FilesUpdaterOpts))

		for _, f := range conf.UpdateFiles {
			exitIfError(updater.Apply(f, newVer))
		}
	}

	herr := hooksExecutor.Success(&hooks.SuccessHookConfig{
		Commits:     commits,
		PrevRelease: release,
		NewRelease: &semrel.Release{
			SHA:     currentSha,
			Version: newVer,
		},
		Changelog: changelogRes,
		RepoInfo:  repoInfo,
	})

	if herr != nil {
		logger.Printf("there was an error executing the hooks plugins: %s", herr.Error())
	}

	logger.Println("done.")
}
