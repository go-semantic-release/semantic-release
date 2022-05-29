package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

const PluginAPI = "https://plugins.go-semantic-release.xyz/api/v1"

type apiPluginAsset struct {
	FileName string
	URL      string
	OS       string
	Arch     string
	Checksum string
}

type apiPluginRelease struct {
	CreatedAt time.Time
	Assets    []*apiPluginAsset
}

func (r *apiPluginRelease) getMatchingAsset() *apiPluginAsset {
	for _, a := range r.Assets {
		if a.OS == runtime.GOOS && a.Arch == runtime.GOARCH {
			return a
		}
	}
	return nil
}

type apiPlugin struct {
	Type          string
	Name          string
	LatestRelease string
	Versions      map[string]*apiPluginRelease
}

func getPluginInfo(name string) (*apiPlugin, error) {
	res, err := http.Get(fmt.Sprintf("%s/plugins/%s.json", PluginAPI, name))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return nil, fmt.Errorf("plugin not found: %s", name)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New("invalid response")
	}
	var plugin *apiPlugin
	if err := json.NewDecoder(res.Body).Decode(&plugin); err != nil {
		return nil, err
	}
	return plugin, nil
}
