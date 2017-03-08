package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/christophwitzko/go-semantic-release"
	"log"
	"os"
)

func main() {
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "github token")
	slug := flag.String("slug", os.Getenv("TRAVIS_REPO_SLUG"), "slug of the repository")
	ghr := flag.Bool("ghr", false, "print ghr flags to stdout")
	flag.Parse()

	logger := log.New(os.Stderr, "[semantic-release]: ", 0)

	if *token == "" {
		logger.Println("github token missing")
		return
	}
	if *slug == "" {
		logger.Println("slug missing")
		return
	}

	repo := semrel.NewRepository(context.TODO(), *slug, *token)

	logger.Println("getting latest release...")
	release, rerr := repo.GetLatestRelease()
	if rerr != nil {
		logger.Println(rerr)
		return
	}
	if release.Version == nil {
		logger.Println("found invalid version")
		return
	}
	logger.Println("found: " + release.Version.String())
	logger.Println("getting commits...")
	commits, err := repo.GetCommits()
	if err != nil {
		logger.Println(err)
		return
	}
	logger.Println("calculating new version...")
	newVer := semrel.GetNewVersion(commits, release)
	if newVer == nil {
		logger.Println("no change")
		return
	}
	logger.Println("new version: " + newVer.String())
	logger.Println("generate release...")
	berr := repo.CreateRelease(commits, release, newVer)
	if berr != nil {
		logger.Println(berr)
		return
	}
	if *ghr {
		fmt.Printf("-u %s -r %s v%s", repo.Owner, repo.Repo, newVer.String())
	}
	logger.Println("done.")
}
