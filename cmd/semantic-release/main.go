package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/semantic-release/go-semantic-release"
	"github.com/semantic-release/go-semantic-release/condition"
	"io/ioutil"
	"log"
	"os"
)

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
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "github token")
	slug := flag.String("slug", os.Getenv("TRAVIS_REPO_SLUG"), "slug of the repository")
	ghr := flag.Bool("ghr", false, "create a .ghr file with the parameters for ghr")
	noci := flag.Bool("noci", false, "run semantic-release locally")
	dry := flag.Bool("dry", false, "do not create release")
	vFile := flag.Bool("vf", false, "create a .version file")
	showVersion := flag.Bool("version", false, "outputs the semantic-release version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("semantic-release v%s", SRVERSION)
		return
	}

	logger := log.New(os.Stderr, "[semantic-release]: ", 0)
	exitIfError := errorHandler(logger)

	if *token == "" {
		exitIfError(errors.New("github token missing"))
	}
	if *slug == "" {
		exitIfError(errors.New("slug missing"))
	}

	repo, nerr := semrel.NewRepository(context.TODO(), *slug, *token)
	exitIfError(nerr)

	logger.Println("getting default branch...")
	defaultBranch, isPrivate, derr := repo.GetInfo()
	exitIfError(derr)
	logger.Println("found default branch: " + defaultBranch)

	if !*noci {
		logger.Println("running CI condition...")
		exitIfError(condition.Travis(*token, defaultBranch, isPrivate))
	}

	logger.Println("getting latest release...")
	release, rerr := repo.GetLatestRelease()
	exitIfError(rerr)
	logger.Println("found version: " + release.Version.String())

	logger.Println("getting commits...")
	commits, cerr := repo.GetCommits()
	exitIfError(cerr)

	logger.Println("calculating new version...")
	newVer := semrel.GetNewVersion(commits, release)
	if newVer == nil {
		exitIfError(errors.New("no change"))
	}
	logger.Println("new version: " + newVer.String())

	if *dry {
		exitIfError(errors.New("DRY RUN: no release was created"))
	}

	logger.Println("creating release...")
	exitIfError(repo.CreateRelease(commits, release, newVer))

	if *ghr {
		exitIfError(ioutil.WriteFile(".ghr", []byte(fmt.Sprintf("-u %s -r %s v%s", repo.Owner, repo.Repo, newVer.String())), 0644))
	}

	if *vFile {
		exitIfError(ioutil.WriteFile(".version", []byte(newVer.String()), 0644))
	}

	logger.Println("done.")
}
