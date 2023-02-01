package registryV1

import (
	"errors"
	"fmt"
	"runtime"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver"
)

type Resolver struct{}

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) ResolvePlugin(pluginInfo *plugin.Info) (*resolver.PluginDownloadInfo, error) {
	pluginAPIRes, err := getPluginInfo(pluginInfo.NormalizedName)
	if err != nil {
		return nil, err
	}

	foundVersion := ""
	if pluginInfo.Constraint == nil {
		foundVersion = pluginAPIRes.LatestRelease
	} else {
		versions := make(semver.Collection, 0)
		for v := range pluginAPIRes.Versions {
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

	releaseAsset := pluginAPIRes.Versions[foundVersion].getMatchingAsset()
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

func (r *Resolver) Names() []string {
	return []string{"registry"}
}
