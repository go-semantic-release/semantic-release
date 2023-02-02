package registryV1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

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
	LatestRelease string
	Versions      map[string]*apiPluginRelease
}

func getPluginInfo(endpoint, name string) (*apiPlugin, error) {
	res, err := http.Get(fmt.Sprintf("%s/plugins/%s.json", endpoint, name))
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
