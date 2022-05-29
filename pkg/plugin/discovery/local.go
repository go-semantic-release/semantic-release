package discovery

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sort"

	"github.com/Masterminds/semver/v3"
)

const PluginDir = ".semrel"

var osArchDir = runtime.GOOS + "_" + runtime.GOARCH

func getPluginPath(name string) string {
	pElem := append([]string{PluginDir}, osArchDir, name)
	return path.Join(pElem...)
}

func ensurePluginDir(pth string) error {
	_, err := os.Stat(pth)
	if os.IsNotExist(err) {
		return os.MkdirAll(pth, 0755)
	}
	return err
}

func getMatchingVersionDir(pth string, cons *semver.Constraints) (string, error) {
	vDirs, err := ioutil.ReadDir(pth)
	if err != nil {
		return "", err
	}
	foundVers := make(semver.Collection, 0)
	for _, f := range vDirs {
		if f.IsDir() {
			fVer, err := semver.NewVersion(f.Name())
			if err != nil {
				continue
			}
			foundVers = append(foundVers, fVer)
		}
	}

	if len(foundVers) == 0 {
		return "", errors.New("no installed version found")
	}
	sort.Sort(sort.Reverse(foundVers))

	if cons == nil {
		return path.Join(pth, foundVers[0].String()), nil
	}

	for _, v := range foundVers {
		if cons.Check(v) {
			return path.Join(pth, v.String()), nil
		}
	}
	return "", errors.New("no matching version found")
}

func findPluginLocally(pth string, cons *semver.Constraints) (string, error) {
	vPth, err := getMatchingVersionDir(pth, cons)
	if err != nil {
		return "", err
	}

	files, err := ioutil.ReadDir(vPth)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", errors.New("no plugins found")
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if f.Mode()&0100 == 0 {
			continue
		}
		return path.Join(vPth, f.Name()), nil
	}
	return "", errors.New("no matching plugin found")
}
