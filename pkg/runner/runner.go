package runner

import (
	"errors"
	fmt "fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/condition"
	"github.com/go-semantic-release/semantic-release/v2/pkg/config"
	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
	"github.com/go-semantic-release/semantic-release/v2/pkg/updater"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type SemanticRelease struct {
	CI                 condition.CICondition
	Prov               provider.Provider
	HooksExecutor      *hooks.ChainedHooksExecutor
	CommitAnalyzer     analyzer.CommitAnalyzer
	ChangelogGenerator generator.ChangelogGenerator
	Updater            *updater.ChainedUpdater
}

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

func (semantic SemanticRelease) Run(conf *config.Config) {
	logger := log.New(os.Stderr, "[go-semantic-release]: ", 0)
	exitIfError := errorHandler(logger)

	if conf.ProviderOpts["token"] == "" {
		conf.ProviderOpts["token"] = conf.Token
	}
	err := semantic.Prov.Init(conf.ProviderOpts)
	exitIfError(err)

	logger.Println("getting default branch...")
	repoInfo, err := semantic.Prov.GetInfo()
	exitIfError(err)
	logger.Println("found default branch: " + repoInfo.DefaultBranch)
	if repoInfo.Private {
		logger.Println("repo is private")
	}

	currentBranch := semantic.CI.GetCurrentBranch()
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

	currentSha := semantic.CI.GetCurrentSHA()
	logger.Println("found current sha: " + currentSha)

	exitIfError(semantic.HooksExecutor.Init(conf.HooksOpts))

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
		err = semantic.CI.RunCondition(conditionConfig)
		if err != nil {
			herr := semantic.HooksExecutor.NoRelease(&hooks.NoReleaseConfig{
				Reason:  hooks.NoReleaseReason_CONDITION,
				Message: err.Error(),
			})
			if herr != nil {
				logger.Printf("there was an error executing the hooks SemanticRelease: %s", herr.Error())
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
	releases, err := semantic.Prov.GetReleases(matchRegex)
	exitIfError(err)
	release, err := semrel.GetLatestReleaseFromReleases(releases, conf.MaintainedVersion)
	exitIfError(err)
	logger.Println("found version: " + release.Version)

	if strings.Contains(conf.MaintainedVersion, "-") && semver.MustParse(release.Version).Prerelease() == "" {
		exitIfError(fmt.Errorf("no pre-release for this version possible"))
	}

	logger.Println("getting commits...")
	rawCommits, err := semantic.Prov.GetCommits(release.SHA, currentSha)
	exitIfError(err)

	exitIfError(semantic.CommitAnalyzer.Init(conf.ChangelogGeneratorOpts))

	commits := semantic.CommitAnalyzer.Analyze(rawCommits)

	logger.Println("calculating new version...")
	newVer := semrel.GetNewVersion(conf, commits, release)
	if newVer == "" {
		herr := semantic.HooksExecutor.NoRelease(&hooks.NoReleaseConfig{
			Reason:  hooks.NoReleaseReason_NO_CHANGE,
			Message: "",
		})
		if herr != nil {
			logger.Printf("there was an error executing the hooks SemanticRelease: %s", herr.Error())
		}
		if conf.AllowNoChanges {
			logger.Println("no change")
			os.Exit(0)
		} else {
			exitIfError(errors.New("no change"), 65)
		}
	}
	logger.Println("new version: " + newVer)

	if conf.Dry {
		exitIfError(errors.New("DRY RUN: no release was created"), 65)
	}

	logger.Println("generating changelog...")
	exitIfError(semantic.ChangelogGenerator.Init(conf.ChangelogGeneratorOpts))

	changelogRes := semantic.ChangelogGenerator.Generate(&generator.ChangelogGeneratorConfig{
		Commits:       commits,
		LatestRelease: release,
		NewVersion:    newVer,
	})
	if conf.Changelog != "" {
		exitIfError(ioutil.WriteFile(conf.Changelog, []byte(changelogRes), 0644))
	}

	logger.Println("creating release...")
	newRelease := &provider.CreateReleaseConfig{
		Changelog:  changelogRes,
		NewVersion: newVer,
		Prerelease: conf.Prerelease,
		Branch:     currentBranch,
		SHA:        currentSha,
	}
	exitIfError(semantic.Prov.CreateRelease(newRelease))

	if conf.Ghr {
		exitIfError(ioutil.WriteFile(".ghr", []byte(fmt.Sprintf("-u %s -r %s v%s", repoInfo.Owner, repoInfo.Repo, newVer)), 0644))
	}

	if conf.VersionFile {
		exitIfError(ioutil.WriteFile(".version", []byte(newVer), 0644))
	}

	if len(conf.UpdateFiles) > 0 {
		exitIfError(semantic.Updater.Init(conf.FilesUpdaterOpts))

		for _, f := range conf.UpdateFiles {
			exitIfError(semantic.Updater.Apply(f, newVer))
		}
	}

	herr := semantic.HooksExecutor.Success(&hooks.SuccessHookConfig{
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
		logger.Printf("there was an error executing the hooks SemanticRelease: %s", herr.Error())
	}

	logger.Println("done.")
}
