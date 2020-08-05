package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/go-semantic-release/semantic-release/pkg/analyzer/commit"
	"github.com/go-semantic-release/semantic-release/pkg/condition"
	"github.com/go-semantic-release/semantic-release/pkg/config"
	"github.com/go-semantic-release/semantic-release/pkg/generator/changelog"
	"github.com/go-semantic-release/semantic-release/pkg/provider"
	"github.com/go-semantic-release/semantic-release/pkg/provider/github"
	"github.com/go-semantic-release/semantic-release/pkg/provider/gitlab"
	"github.com/go-semantic-release/semantic-release/pkg/semrel"
	"github.com/go-semantic-release/semantic-release/pkg/updater"
	_ "github.com/go-semantic-release/semantic-release/pkg/updater/npm"
	"github.com/urfave/cli/v2"
)

// SRVERSION is the semantic-release version (added at compile time)
var SRVERSION string

func errorHandler(logger *log.Logger) func(error, ...int) {
	return func(err error, exitCode ...int) {
		if err != nil {
			logger.Println(err)
			if len(exitCode) == 1 {
				os.Exit(exitCode[0])
				return
			}
			os.Exit(1)
		}
	}
}

func main() {
	app := &cli.App{
		Name:     "semantic-release",
		Usage:    "automates the package release workflow including: determining the next version number and generating the change log",
		Version:  SRVERSION,
		Commands: nil,
		Flags:    config.CliFlags,
		Action:   cliHandler,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(1)
	}
}

func cliHandler(c *cli.Context) error {

	logger := log.New(os.Stderr, "[semantic-release]: ", 0)
	exitIfError := errorHandler(logger)

	conf, err := config.NewConfig(c)
	exitIfError(err)

	ci := condition.NewCI()
	logger.Printf("detected CI: %s\n", ci.Name())

	var repo provider.Repository

	if conf.GitLab {
		repo, err = gitlab.NewRepository(c.Context, conf.GitLabBaseURL, conf.Token, ci.GetCurrentBranch(), conf.GitLabProjectID)
	} else {
		repo, err = github.NewRepository(c.Context, conf.GheHost, conf.Slug, conf.Token)
	}

	logger.Printf("releasing on: %s\n", repo.Provider())

	exitIfError(err)

	logger.Println("getting default branch...")
	repoInfo, err := repo.GetInfo()
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

	if conf.BetaRelease.MaintainedVersion != "" && currentBranch == repoInfo.DefaultBranch {
		exitIfError(fmt.Errorf("maintained version not allowed on default branch"))
	}

	if conf.BetaRelease.MaintainedVersion != "" {
		logger.Println("found maintained version: " + conf.BetaRelease.MaintainedVersion)
		repoInfo.DefaultBranch = "*"
	}

	currentSha := ci.GetCurrentSHA()
	logger.Println("found current sha: " + currentSha)

	if !conf.Noci {
		logger.Println("running CI condition...")
		config := map[string]interface{}{
			"token":         conf.Token,
			"defaultBranch": repoInfo.DefaultBranch,
			"private":       repoInfo.Private || conf.TravisCom,
		}
		exitIfError(ci.RunCondition(config), 66)
	}

	logger.Println("getting latest release...")
	var matchRegex *regexp.Regexp
	match := strings.TrimSpace(conf.Match)
	if match != "" {
		logger.Printf("getting latest release matching %s...", match)
		matchRegex = regexp.MustCompile("^" + match)
	}
	releases, err := repo.GetReleases(matchRegex)
	exitIfError(err)
	release, err := semrel.GetLatestReleaseFromReleases(releases, conf.BetaRelease.MaintainedVersion)
	exitIfError(err)
	logger.Println("found version: " + release.Version)

	if strings.Contains(conf.BetaRelease.MaintainedVersion, "-") && semver.MustParse(release.Version).Prerelease() == "" {
		exitIfError(fmt.Errorf("no pre-release for this version possible"))
	}

	logger.Println("getting commits...")
	rawCommits, err := repo.GetCommits(currentSha)
	exitIfError(err)

	commitAnalyzer := &commit.DefaultAnalyzer{}
	commits := commitAnalyzer.Analyze(rawCommits)

	logger.Println("calculating new version...")
	newVer := semrel.GetNewVersion(conf, commits, release)
	if newVer == "" {
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
	changelogGenerator := &changelog.DefaultGenerator{}
	changelogRes := changelogGenerator.Generate(&changelog.Config{
		Commits:       commits,
		LatestRelease: release,
		NewVersion:    newVer,
	})
	if conf.Changelog != "" {
		exitIfError(ioutil.WriteFile(conf.Changelog, []byte(changelogRes), 0644))
	}

	logger.Println("creating release...")
	newRelease := &provider.RepositoryRelease{
		Changelog:  changelogRes,
		NewVersion: newVer,
		Prerelease: conf.Prerelease,
		Branch:     currentBranch,
		SHA:        currentSha,
	}
	exitIfError(repo.CreateRelease(newRelease))

	if conf.Ghr {
		exitIfError(ioutil.WriteFile(".ghr", []byte(fmt.Sprintf("-u %s -r %s v%s", repoInfo.Owner, repoInfo.Repo, newVer)), 0644))
	}

	if conf.Vf {
		exitIfError(ioutil.WriteFile(".version", []byte(newVer), 0644))
	}

	if conf.Update != "" {
		exitIfError(updater.Apply(conf.Update, newVer))
	}

	logger.Println("done.")
	return nil
}
