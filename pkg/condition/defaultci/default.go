package defaultci

import (
	"io/ioutil"
	"strings"

	"github.com/go-semantic-release/semantic-release/pkg/condition"
	"github.com/go-semantic-release/semantic-release/pkg/plugin"
	"github.com/urfave/cli/v2"
)

func ReadGitHead() string {
	data, err := ioutil.ReadFile(".git/HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(string(data), "ref: refs/heads/"))
}

type DefaultCI struct {
}

func (d *DefaultCI) Name() string {
	return "none"
}

func (d *DefaultCI) RunCondition(map[string]string) error {
	return nil
}

func (d *DefaultCI) GetCurrentBranch() string {
	return ReadGitHead()
}

func (d *DefaultCI) GetCurrentSHA() string {
	return ReadGitHead()
}

func Main(c *cli.Context) error {
	plugin.Serve(&plugin.ServeOpts{
		CICondition: func() condition.CICondition {
			return &DefaultCI{}
		},
	})
	return nil
}
