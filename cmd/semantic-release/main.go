package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-semantic-release/semantic-release/pkg/condition"
	"github.com/go-semantic-release/semantic-release/pkg/config"
	"github.com/go-semantic-release/semantic-release/pkg/semrel"
	"github.com/go-semantic-release/semantic-release/pkg/update"
	"github.com/urfave/cli"
)

// SRVERSION is the semantic-release version (added at compile time)
var SRVERSION string

func errorHandler(logger *log.Logger) func(error) {
	return func(err error) {
		if err != nil {
			logger.Println(err)
			os.Exit(1)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "semantic-release"
	app.Usage = "automates the package release workflow including: determining the next version number and generating the change log"
	app.Version = SRVERSION
	app.Flags = config.CliFlags
	app.Action = cliHandler

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func cliHandler(c *cli.Context) error {
	conf := config.NewConfig(c)

	logger := log.New(os.Stderr, "[semantic-release]: ", 0)
	exitIfError := errorHandler(logger)

	repo, err := semrel.NewRepository(context.TODO(), conf.GheHost, conf.Slug, conf.Token)
	exitIfError(err)

	logger.Println("getting default branch...")
	defaultBranch, isPrivate, err := repo.GetInfo()
	exitIfError(err)
	logger.Println("found default branch: " + defaultBranch)
	if isPrivate {
		logger.Println("repo is private")
	}

	currentBranch := condition.GetCurrentBranch()
	if currentBranch == "" {
		exitIfError(fmt.Errorf("current branch not found"))
	}
	logger.Println("found current branch: " + currentBranch)

	if conf.BetaRelease.MaintainedVersion != "" && currentBranch == defaultBranch {
		exitIfError(fmt.Errorf("maintained version not allowed on default branch"))
	}

	if conf.BetaRelease.MaintainedVersion != "" {
		logger.Println("found maintained version: " + conf.BetaRelease.MaintainedVersion)
		defaultBranch = "*"
	}

	currentSha := condition.GetCurrentSHA()
	logger.Println("found current sha: " + currentSha)

	if !conf.Noci {
		logger.Println("running CI condition...")
		exitIfError(condition.Travis(conf.Token, defaultBranch, isPrivate || conf.TravisCom))
	}

	logger.Println("getting latest release...")
	release, err := repo.GetLatestRelease(conf.BetaRelease.MaintainedVersion)
	exitIfError(err)
	logger.Println("found version: " + release.Version.String())

	if strings.Contains(conf.BetaRelease.MaintainedVersion, "-") && release.Version.Prerelease() == "" {
		exitIfError(fmt.Errorf("no pre-release for this version possible"))
	}

	logger.Println("getting commits...")
	commits, err := repo.GetCommits(currentSha)
	exitIfError(err)

	logger.Println("calculating new version...")
	newVer := semrel.GetNewVersion(commits, release)
	if newVer == nil {
		exitIfError(errors.New("no change"))
	}
	logger.Println("new version: " + newVer.String())

	if conf.Dry {
		exitIfError(errors.New("DRY RUN: no release was created"))
	}

	logger.Println("generating changelog...")
	changelog := semrel.GetChangelog(commits, release, newVer)
	if conf.Changelog != "" {
		exitIfError(ioutil.WriteFile(conf.Changelog, []byte(changelog), 0644))
	}

	logger.Println("creating release...")
	exitIfError(repo.CreateRelease(changelog, newVer, conf.Prerelease, currentBranch, currentSha))

	if conf.Ghr {
		exitIfError(ioutil.WriteFile(".ghr", []byte(fmt.Sprintf("-u %s -r %s v%s", repo.Owner, repo.Repo, newVer.String())), 0644))
	}

	if conf.Vf {
		exitIfError(ioutil.WriteFile(".version", []byte(newVer.String()), 0644))
	}

	if conf.Update != "" {
		exitIfError(update.Apply(conf.Update, newVer.String()))
	}

	logger.Println("done.")
	return nil
}
