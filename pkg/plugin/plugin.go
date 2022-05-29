package plugin

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/condition"
	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/updater"
)

type PluginInfo struct {
	Type           string
	Name           string
	NormalizedName string
	Constraint     *semver.Constraints
	Resolver       string
	PluginPath     string
	BinPath        string
}

func normalizedPluginType(t string) string {
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

func GetPluginInfo(pluginType, pluginName string) (*PluginInfo, error) {
	nPluginType := normalizedPluginType(pluginType)
	if nPluginType == "" {
		return nil, fmt.Errorf("invalid plugin type: %s", pluginType)
	}
	resolver := "default"
	name := pluginName
	normalizedName := fmt.Sprintf("%s-%s", nPluginType, pluginName)
	if strings.Contains(name, ":") {
		parts := strings.SplitN(name, ":", 2)
		resolver = parts[0]
		name = parts[1]
		normalizedName = strings.ReplaceAll(normalizedName, ":", "-")
	}
	if strings.Contains(name, "/") {
		parts := strings.SplitN(name, "/", 2)
		name = parts[1]
		normalizedName = strings.ReplaceAll(normalizedName, "/", "-")
	}

	var constraint *semver.Constraints
	if strings.Contains(name, "@") {
		parts := strings.SplitN(name, "@", 2)
		v, err := semver.NewConstraint(parts[1])
		if err != nil {
			return nil, err
		}
		name = parts[0]
		constraint = v
		parts = strings.SplitN(normalizedName, "@", 2)
		normalizedName = parts[0]
	}

	return &PluginInfo{
		Type:           pluginType,
		Name:           name,
		NormalizedName: normalizedName,
		Constraint:     constraint,
		Resolver:       resolver,
	}, nil
}
