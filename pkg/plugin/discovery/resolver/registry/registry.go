package registry

import (
	"context"
	"fmt"
	"runtime"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-semantic-release/plugin-registry/pkg/client"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver"
)

const DefaultEndpoint = "https://registry-staging.go-semantic-release.xyz/api/v2"

type Resolver struct {
	client *client.Client
}

func NewResolver() *Resolver {
	return &Resolver{
		client: client.New(DefaultEndpoint),
	}
}

func (r *Resolver) ResolvePlugin(pluginInfo *plugin.Info) (*resolver.PluginDownloadInfo, error) {
	getPluginRes, err := r.client.GetPlugin(context.Background(), pluginInfo.ShortNormalizedName)
	if err != nil {
		return nil, err
	}

	osArch := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	if pluginInfo.Constraint == nil {
		foundAsset := getPluginRes.LatestRelease.Assets[osArch]
		if foundAsset == nil {
			return nil, fmt.Errorf("a matching plugin was not found for %s/%s", runtime.GOOS, runtime.GOARCH)
		}
		return &resolver.PluginDownloadInfo{
			URL:      foundAsset.URL,
			Checksum: foundAsset.Checksum,
			FileName: foundAsset.FileName,
			Version:  getPluginRes.LatestRelease.Version,
		}, nil
	}

	foundVersion := ""
	versions := make(semver.Collection, 0)
	for _, v := range getPluginRes.Versions {
		pv, sErr := semver.NewVersion(v)
		if sErr != nil {
			return nil, sErr
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
	if foundVersion == "" {
		return nil, fmt.Errorf("no matching version was found")
	}

	pluginRelease, err := r.client.GetPluginRelease(context.Background(), pluginInfo.ShortNormalizedName, foundVersion)
	if err != nil {
		return nil, err
	}
	foundAsset := pluginRelease.Assets[osArch]
	if foundAsset == nil {
		return nil, fmt.Errorf("a matching plugin was not found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	return &resolver.PluginDownloadInfo{
		URL:      foundAsset.URL,
		Checksum: foundAsset.Checksum,
		FileName: foundAsset.FileName,
		Version:  getPluginRes.LatestRelease.Version,
	}, nil
}

func (r *Resolver) BatchResolvePlugins(pluginInfos []*plugin.Info) (*resolver.BatchPluginDownloadInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *Resolver) Names() []string {
	// TODO: this should be registry when the registry is ready
	return []string{"registry-beta"}
}
