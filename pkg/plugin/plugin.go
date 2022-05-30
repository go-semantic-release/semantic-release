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
	Type                string
	Name                string
	NormalizedName      string
	ShortNormalizedName string
	Constraint          *semver.Constraints
	Resolver            string
	PluginPath          string
	BinPath             string
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
		parts := strings.Split(name, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid plugin name format")
		}
		resolver = parts[0]
		name = parts[1]
		normalizedName = strings.ReplaceAll(normalizedName, ":", "-")
	}

	if strings.Contains(name, "/") {
		parts := strings.Split(name, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid plugin name format")
		}
		name = parts[1]
		normalizedName = strings.ReplaceAll(normalizedName, "/", "-")
	}

	var constraint *semver.Constraints
	if strings.Contains(name, "@") {
		parts := strings.Split(name, "@")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid plugin name format")
		}
		v, err := semver.NewConstraint(parts[1])
		if err != nil {
			return nil, err
		}
		name = parts[0]
		constraint = v
		normalizedName, _, _ = strings.Cut(normalizedName, "@")
	}

	return &PluginInfo{
		Type:                pluginType,
		Name:                name,
		NormalizedName:      normalizedName,
		ShortNormalizedName: fmt.Sprintf("%s-%s", nPluginType, name),
		Constraint:          constraint,
		Resolver:            resolver,
	}, nil
}
