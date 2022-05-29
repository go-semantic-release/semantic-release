package discovery

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/condition"
	"github.com/go-semantic-release/semantic-release/v2/pkg/config"
	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/updater"
)

type Discovery struct {
	config *config.Config
}

func New(config *config.Config) (*Discovery, error) {
	return &Discovery{config}, nil
}

func getPluginType(t string) string {
	switch t {
	case analyzer.CommitAnalyzerPluginName:
		return "commit-analyzer"
	case condition.CIConditionPluginName:
		return "condition"
	case generator.ChangelogGeneratorPluginName:
		return "changelog-generator"
	case provider.PluginName:
		return "provider"
	case updater.FilesUpdaterPluginName:
		return "files-updater"
	case hooks.PluginName:
		return "hooks"
	}
	return ""
}

func (d *Discovery) FindPlugin(t, name string) (*plugin.PluginOpts, error) {
	pType := getPluginType(t)
	if pType == "" {
		return nil, errors.New("invalid plugin type")
	}

	var cons *semver.Constraints
	if ve := strings.SplitN(name, "@", 2); len(ve) > 1 {
		v, err := semver.NewConstraint(ve[1])
		if err != nil {
			return nil, err
		}
		name = ve[0]
		cons = v
	}

	pName := strings.ToLower(pType + "-" + name)
	pPath := getPluginPath(pName)
	if err := ensurePluginDir(pPath); err != nil {
		return nil, err
	}

	binPath, err := findPluginLocally(pPath, cons)
	if err != nil {
		binPath, err = fetchPlugin(pName, pPath, cons, d.config.ShowProgress)
		if err != nil {
			return nil, err
		}
	}

	cmd := exec.Command(binPath)
	cmd.SysProcAttr = GetSysProcAttr()

	return &plugin.PluginOpts{
		Type:       t,
		PluginName: pName,
		Cmd:        cmd,
	}, nil
}
