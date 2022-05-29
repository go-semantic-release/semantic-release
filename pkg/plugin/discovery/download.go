package discovery

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver"
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

func downloadPlugin(pluginInfo *plugin.PluginInfo, downloadInfo *resolver.PluginDownloadInfo, showProgress bool) (string, error) {
	targetPath := path.Join(pluginInfo.PluginPath, downloadInfo.Version, downloadInfo.FileName)
	req, err := grab.NewRequest(targetPath, downloadInfo.URL)
	if err != nil {
		return "", err
	}
	if downloadInfo.Checksum != "" {
		sum, err := hex.DecodeString(downloadInfo.Checksum)
		if err != nil {
			return "", err
		}
		req.SetChecksum(sha256.New(), sum, true)
	}

	res := grab.DefaultClient.Do(req)
	if showProgress {
		showDownloadProgressBar(pluginInfo.NormalizedName, res)
	}
	if err := res.Err(); err != nil {
		return "", err
	}
	if err := os.Chmod(res.Filename, 0755); err != nil {
		return "", err
	}
	return res.Filename, nil
}
