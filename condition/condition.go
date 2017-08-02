package condition

import (
	"io/ioutil"
	"os"
	"strings"
	"gopkg.in/src-d/go-git.v4"
)

func readGitHead() string {
	data, err := ioutil.ReadFile(".git/HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(string(data), "ref: refs/heads/"))
}

func GetCurrentBranch() string {
	if val := os.Getenv("TRAVIS_BRANCH"); val != "" {
		return val
	}
	return readGitHead()
}

type CommitInfo struct {
	Branch string
	SHA    string
}

func GetCurCommitInfo() (*CommitInfo, error) {


	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}

	headRef, err := repo.Head()

	if err != nil {
		return nil, err
	}

	return &CommitInfo{
		Branch: headRef.Name().Short(),
		SHA:    headRef.Hash().String(),
	}, nil
}
