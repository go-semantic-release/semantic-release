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

type Info struct {
	Type                string
	Name                string
	NormalizedName      string
	ShortNormalizedName string
	Constraint          *semver.Constraints
	Resolver            string
	RepoSlug            string
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

var nameNormalizer = strings.NewReplacer(":", "-", "/", "-")

func GetPluginInfo(pluginType, pluginName string) (*Info, error) {
	nPluginType := normalizedPluginType(pluginType)
	if nPluginType == "" {
		return nil, fmt.Errorf("invalid plugin type: %s", pluginType)
	}
	resolver := "default"
	repoSlug := ""
	name := pluginName
	normalizedName := nameNormalizer.Replace(fmt.Sprintf("%s-%s", nPluginType, pluginName))

	if strings.Contains(name, ":") {
		parts := strings.Split(name, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid plugin name format")
		}
		resolver = parts[0]
		name = parts[1]
	}

	if strings.Contains(name, "/") {
		slashParts := strings.Split(name, "/")
		pName := slashParts[len(slashParts)-1]
		// remove version constraint from the slug
		if strings.Contains(pName, "@") {
			nameWithoutVersion, _, _ := strings.Cut(pName, "@")
			repoSlug = strings.Join(slashParts[:len(slashParts)-1], "/") + "/" + nameWithoutVersion
		} else {
			repoSlug = name
		}
		// the last part is the plugin name
		name = pName
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

		// remove version constraint from the normalized name
		normalizedParts := strings.Split(normalizedName, "@")
		normalizedName = strings.Join(normalizedParts[:len(normalizedParts)-1], "@")
	}

	return &Info{
		Type:                pluginType,
		Name:                name,
		NormalizedName:      normalizedName,
		ShortNormalizedName: fmt.Sprintf("%s-%s", nPluginType, name),
		Constraint:          constraint,
		Resolver:            resolver,
		RepoSlug:            repoSlug,
	}, nil
}
