package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/christophwitzko/go-semantic-release"
	"github.com/christophwitzko/go-semantic-release/condition"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "github token")
	slug := flag.String("slug", os.Getenv("TRAVIS_REPO_SLUG"), "slug of the repository")
	ghr := flag.Bool("ghr", false, "print ghr flags to stdout")
	noci := flag.Bool("noci", false, "run semantic-release locally")
	dry := flag.Bool("dry", false, "do not create release")
	vFile := flag.Bool("vf", false, "create a .version file")
	flag.Parse()

	logger := log.New(os.Stderr, "[semantic-release]: ", 0)

	if *token == "" {
		logger.Println("github token missing")
		os.Exit(1)
		return
	}
	if *slug == "" {
		logger.Println("slug missing")
		os.Exit(1)
		return
	}

	repo := semrel.NewRepository(context.TODO(), *slug, *token)

	logger.Println("getting default branch...")
	defaultBranch, isPrivate, derr := repo.GetInfo()
	if derr != nil {
		logger.Println(derr)
		os.Exit(1)
		return
	}
	logger.Println("found default branch: " + defaultBranch)

	if !*noci {
		logger.Println("running CI condition...")
		if err := condition.Travis(*token, defaultBranch, isPrivate); err != nil {
			logger.Println(err)
			os.Exit(1)
			return
		}
	}

	logger.Println("getting latest release...")
	release, rerr := repo.GetLatestRelease()
	if rerr != nil {
		logger.Println(rerr)
		os.Exit(1)
		return
	}
	if release.Version == nil {
		logger.Println("found invalid version")
		os.Exit(1)
		return
	}
	logger.Println("found: " + release.Version.String())

	logger.Println("getting commits...")
	commits, err := repo.GetCommits()
	if err != nil {
		logger.Println(err)
		os.Exit(1)
		return
	}

	logger.Println("calculating new version...")
	newVer := semrel.GetNewVersion(commits, release)
	if newVer == nil {
		logger.Println("no change")
		os.Exit(1)
		return
	}
	logger.Println("new version: " + newVer.String())

	if *dry {
		logger.Println("DRY RUN: no release was created")
		os.Exit(1)
		return
	}

	logger.Println("generating release...")
	berr := repo.CreateRelease(commits, release, newVer)
	if berr != nil {
		logger.Println(berr)
		os.Exit(1)
		return
	}

	if *ghr {
		fmt.Printf("-u %s -r %s v%s", repo.Owner, repo.Repo, newVer.String())
	}

	if *vFile {
		werr := ioutil.WriteFile(".version", []byte(newVer.String()), 0644)
		if werr != nil {
			logger.Println(werr)
			os.Exit(1)
			return
		}
	}

	logger.Println("done.")
}
