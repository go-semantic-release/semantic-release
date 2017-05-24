package condition

import (
	"io/ioutil"
	"os"
	"strings"
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
