package condition

import (
	"os"
)

func GetCurrentBranch() string {
	// TODO: support other CIs
	return os.Getenv("TRAVIS_BRANCH")
}
