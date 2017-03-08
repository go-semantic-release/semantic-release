package condition

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func Travis(token, defaultBranch string) error {
	if os.Getenv("TRAVIS") != "true" {
		return errors.New("semantic-release didn’t run on Travis CI and therefore a new version won’t be published.")
	}
	if val, ok := os.LookupEnv("TRAVIS_PULL_REQUEST"); ok && val != "false" {
		return errors.New("This test run was triggered by a pull request and therefore a new version won’t be published.")
	}
	if _, ok := os.LookupEnv("TRAVIS_TAG"); ok {
		return errors.New("This test run was triggered by a git tag and therefore a new version won’t be published.")
	}
	if branch := os.Getenv("TRAVIS_BRANCH"); branch != defaultBranch {
		return errors.New(fmt.Sprintf("This test run was triggered on the branch %s, while semantic-release is configured to only publish from %s.", branch, defaultBranch))
	}
	if !strings.HasSuffix(os.Getenv("TRAVIS_JOB_NUMBER"), ".1") {
		return errors.New("This test run is not the build leader and therefore a new version won’t be published.")
	}
	if os.Getenv("TRAVIS_TEST_RESULT") == "1" {
		return errors.New("In this test run not all jobs passed and therefore a new version won’t be published.")
	}
	if os.Getenv("TRAVIS_TEST_RESULT") != "0" {
		return errors.New("Not running in Travis after_success hook.")
	}
	return nil
}
