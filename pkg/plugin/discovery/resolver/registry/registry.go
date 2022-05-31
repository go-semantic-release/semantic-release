package registry

import (
	"errors"
	"fmt"
	"runtime"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver"
)

type RegistryResolver struct {
}

func NewResolver() *RegistryResolver {
	return &RegistryResolver{}
}

func (r *RegistryResolver) ResolvePlugin(pluginInfo *plugin.PluginInfo) (*resolver.PluginDownloadInfo, error) {
	pluginApiRes, err := getPluginInfo(pluginInfo.NormalizedName)
	if err != nil {
		return nil, err
	}

	foundVersion := ""
	if pluginInfo.Constraint == nil {
		foundVersion = pluginApiRes.LatestRelease
	} else {
		versions := make(semver.Collection, 0)
		for v := range pluginApiRes.Versions {
			pv, err := semver.NewVersion(v)
			if err != nil {
				return nil, err
			}
			versions = append(versions, pv)
		}
		sort.Sort(sort.Reverse(versions))
		for _, v := range versions {
			if pluginInfo.Constraint.Check(v) {
				foundVersion = v.String()
				break
			}
		}
	}

	if foundVersion == "" {
		return nil, errors.New("version not found")
	}

	releaseAsset := pluginApiRes.Versions[foundVersion].getMatchingAsset()
	if releaseAsset == nil {
		return nil, fmt.Errorf("a matching plugin was not found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	return &resolver.PluginDownloadInfo{
		URL:      releaseAsset.URL,
		Checksum: releaseAsset.Checksum,
		FileName: releaseAsset.FileName,
		Version:  foundVersion,
	}, nil
}

func (r *RegistryResolver) Names() []string {
	return []string{"registry"}
}
