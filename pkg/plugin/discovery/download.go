package discovery

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/cavaliergopher/grab/v3"
	"github.com/schollz/progressbar/v3"
)

func showDownloadProgressBar(name string, res *grab.Response) {
	bar := progressbar.NewOptions64(
		res.Size(),
		progressbar.OptionSetDescription(name),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(40),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetPredictTime(false),
	)
	t := time.NewTicker(100 * time.Millisecond)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-t.C:
				_ = bar.Set64(res.BytesComplete())
			case <-res.Done:
				_ = bar.Finish()
				t.Stop()
				done <- struct{}{}
				return
			}
		}
	}()
	<-done
}

func downloadPlugin(name, targetPath, downloadUrl, checksum string, showProgress bool) (string, error) {
	req, err := grab.NewRequest(targetPath, downloadUrl)
	if err != nil {
		return "", err
	}
	if checksum != "" {
		sum, err := hex.DecodeString(checksum)
		if err != nil {
			return "", err
		}
		req.SetChecksum(sha256.New(), sum, true)
	}

	res := grab.DefaultClient.Do(req)
	if showProgress {
		showDownloadProgressBar(name, res)
	}
	if err := res.Err(); err != nil {
		return "", err
	}
	if err := os.Chmod(res.Filename, 0755); err != nil {
		return "", err
	}
	return res.Filename, nil
}

func fetchPlugin(name, pth string, cons *semver.Constraints, showProgress bool) (string, error) {
	pluginInfo, err := getPluginInfo(name)
	if err != nil {
		return "", err
	}

	foundVersion := ""
	if cons == nil {
		foundVersion = pluginInfo.LatestRelease
	} else {
		versions := make(semver.Collection, 0)
		for v := range pluginInfo.Versions {
			pv, err := semver.NewVersion(v)
			if err != nil {
				return "", err
			}
			versions = append(versions, pv)
		}
		sort.Sort(sort.Reverse(versions))
		for _, v := range versions {
			if cons.Check(v) {
				foundVersion = v.String()
				break
			}
		}
	}

	if foundVersion == "" {
		return "", errors.New("version not found")
	}

	releaseAsset := pluginInfo.Versions[foundVersion].getMatchingAsset()
	if releaseAsset == nil {
		return "", fmt.Errorf("a matching plugin was not found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	targetPath := path.Join(pth, foundVersion, releaseAsset.FileName)
	return downloadPlugin(name, targetPath, releaseAsset.URL, releaseAsset.Checksum, showProgress)
}
