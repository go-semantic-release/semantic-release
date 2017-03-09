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

var SRVERSION string

func main() {
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "github token")
	slug := flag.String("slug", os.Getenv("TRAVIS_REPO_SLUG"), "slug of the repository")
	ghr := flag.Bool("ghr", false, "create a .ghr file with the parameters for ghr")
	noci := flag.Bool("noci", false, "run semantic-release locally")
	dry := flag.Bool("dry", false, "do not create release")
	vFile := flag.Bool("vf", false, "create a .version file")
	flag.Parse()

	logger := log.New(os.Stderr, "[semantic-release]: ", 0)
	logger.Println("cli version: " + SRVERSION)

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

	repo, nerr := semrel.NewRepository(context.TODO(), *slug, *token)
	if nerr != nil {
		logger.Println(nerr)
		os.Exit(1)
		return
	}

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
	commits, cerr := repo.GetCommits()
	if cerr != nil {
		logger.Println(cerr)
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

	logger.Println("creating release...")
	berr := repo.CreateRelease(commits, release, newVer)
	if berr != nil {
		logger.Println(berr)
		os.Exit(1)
		return
	}

	if *ghr {
		gerr := ioutil.WriteFile(".ghr", []byte(fmt.Sprintf("-u %s -r %s v%s", repo.Owner, repo.Repo, newVer.String())), 0644)
		if gerr != nil {
			logger.Println(gerr)
			os.Exit(1)
			return
		}
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
