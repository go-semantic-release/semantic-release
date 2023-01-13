package discovery

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
)

const PluginDir = ".semrel"

var osArchDir = runtime.GOOS + "_" + runtime.GOARCH

func setAndEnsurePluginPath(pluginInfo *plugin.PluginInfo) error {
	pluginPath := path.Join(PluginDir, osArchDir, pluginInfo.NormalizedName)
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		err = os.MkdirAll(pluginPath, 0o755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	pluginInfo.PluginPath = pluginPath
	return nil
}

var ErrPluginNotFound = errors.New("no plugin was found")

func getMatchingVersionDir(pluginInfo *plugin.PluginInfo) (string, error) {
	vDirs, err := os.ReadDir(pluginInfo.PluginPath)
	if err != nil {
		return "", err
	}
	foundVersions := make(semver.Collection, 0)
	for _, f := range vDirs {
		if f.IsDir() {
			fVer, err := semver.NewVersion(f.Name())
			if err != nil {
				continue
			}
			foundVersions = append(foundVersions, fVer)
		}
	}

	if len(foundVersions) == 0 {
		return "", nil
	}
	sort.Sort(sort.Reverse(foundVersions))

	if pluginInfo.Constraint == nil {
		return path.Join(pluginInfo.PluginPath, foundVersions[0].String()), nil
	}

	for _, v := range foundVersions {
		if pluginInfo.Constraint.Check(v) {
			return path.Join(pluginInfo.PluginPath, v.String()), nil
		}
	}
	return "", nil
}

func findPluginLocally(pluginInfo *plugin.PluginInfo) (string, error) {
	vPth, err := getMatchingVersionDir(pluginInfo)
	if err != nil {
		return "", err
	}

	if vPth == "" {
		return "", ErrPluginNotFound
	}

	files, err := os.ReadDir(vPth)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", ErrPluginNotFound
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fInfo, err := f.Info()
		if err != nil {
			return "", fmt.Errorf("failed to get file info for %s: %w", f.Name(), err)
		}
		// check if the file is executable by the user
		if fInfo.Mode()&0o100 == 0 {
			continue
		}
		return path.Join(vPth, f.Name()), nil
	}
	return "", ErrPluginNotFound
}
