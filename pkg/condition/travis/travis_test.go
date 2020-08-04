package travis

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTravisValid(t *testing.T) {
	travis := TravisCI{}
	os.Setenv("TRAVIS", "")
	err := travis.RunCondition(map[string]interface{}{"token": "", "defaultBranch": "", "private": false})
	assert.EqualError(t, err, "semantic-release didn’t run on Travis CI and therefore a new version won’t be published.")
}
