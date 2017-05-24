package condition

import (
	"errors"
	"fmt"
	"github.com/Ableton/go-travis"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func Travis(token, defaultBranch string, private bool) error {
	logger := log.New(os.Stderr, "[condition-travis]: ", 0)
	if os.Getenv("TRAVIS") != "true" {
		return errors.New("semantic-release didn’t run on Travis CI and therefore a new version won’t be published.")
	}
	if val, ok := os.LookupEnv("TRAVIS_PULL_REQUEST"); ok && val != "false" {
		return errors.New("This test run was triggered by a pull request and therefore a new version won’t be published.")
	}
	if os.Getenv("TRAVIS_TAG") != "" {
		return errors.New("This test run was triggered by a git tag and therefore a new version won’t be published.")
	}
	if branch := os.Getenv("TRAVIS_BRANCH"); defaultBranch != "*" && branch != defaultBranch {
		return fmt.Errorf("This test run was triggered on the branch %s, while semantic-release is configured to only publish from %s.", branch, defaultBranch)
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

	buildId, _ := strconv.ParseUint(os.Getenv("TRAVIS_BUILD_ID"), 10, 32)
	currentJobId, _ := strconv.ParseUint(os.Getenv("TRAVIS_JOB_ID"), 10, 32)
	if buildId < 1 || currentJobId < 1 {
		return errors.New("could not parse TRAVIS_BUILD_ID/TRAVIS_JOB_ID")
	}

	endpoint := travis.TRAVIS_API_DEFAULT_URL
	if private {
		logger.Println("repo is private")
		endpoint = travis.TRAVIS_API_PRO_URL
	}

	client := travis.NewClient(endpoint, "")
	client.Headers["User-Agent"] = "Travis"
	if _, _, err := client.Authentication.UsingGithubToken(token); err != nil {
		return err
	}

	for i := 1; i <= 100; i++ {
		jobs, _, err := client.Jobs.ListFromBuild(uint(buildId))
		if err != nil {
			return err
		}

		successes := 0
		for _, job := range jobs {
			if job.Id == uint(currentJobId) || job.AllowFailure || job.State == "passed" {
				successes++
				continue
			}

			if job.State == "created" || job.State == "started" {
				logger.Printf("Aborting attempt %d, because job %s is still pending.\n", i, job.Number)
				break
			}

			if job.State == "errored" || job.State == "failed" {
				logger.Printf("Aborting attempt %d. Job %s failed.\n", i, job.Number)
				return errors.New("In this test run not all jobs passed and therefore a new version won’t be published.")
			}
		}
		if successes >= len(jobs) {
			logger.Printf("Success at attempt %d. All %d jobs passed.\n", i, successes)
			return nil
		}
		time.Sleep(3 * time.Second)
	}
	return errors.New("Timeout. Could not get accumulated results after 100 attempts.")
}
