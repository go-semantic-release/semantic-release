package discovery

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
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

func extractFileFromTarGz(name, inputFile, outputFile string) error {
	compressedFile, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer compressedFile.Close()

	decompressedFile, err := gzip.NewReader(compressedFile)
	if err != nil {
		return err
	}
	defer decompressedFile.Close()

	tarReader := tar.NewReader(decompressedFile)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("could not extract file")
		}
		if err != nil {
			return err
		}
		if header.Typeflag == tar.TypeReg && strings.HasPrefix(header.Name, name) {
			outFile, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0o755)
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, tarReader)
			outFile.Close()
			if err != nil {
				return err
			}
			return nil
		}
	}
}

var tgzRegexp = regexp.MustCompile(`^(.*)\.(tgz|tar\.gz)$`)

func downloadPlugin(pluginInfo *plugin.Info, downloadInfo *resolver.PluginDownloadInfo, showProgress bool) (string, error) {
	versionDir := path.Join(pluginInfo.PluginPath, downloadInfo.Version)
	targetFile := path.Join(versionDir, downloadInfo.FileName)
	req, err := grab.NewRequest(targetFile, downloadInfo.URL)
	if err != nil {
		return "", err
	}
	if downloadInfo.Checksum != "" {
		sum, decErr := hex.DecodeString(downloadInfo.Checksum)
		if decErr != nil {
			return "", decErr
		}
		req.SetChecksum(sha256.New(), sum, true)
	}

	res := grab.DefaultClient.Do(req)
	if showProgress {
		showDownloadProgressBar(pluginInfo.ShortNormalizedName, res)
	}
	err = res.Err()
	if err != nil {
		return "", err
	}

	tgzMatch := tgzRegexp.FindStringSubmatch(downloadInfo.FileName)
	if len(tgzMatch) > 2 {
		outFile := path.Join(versionDir, tgzMatch[1])
		if err = extractFileFromTarGz(pluginInfo.Name, targetFile, outFile); err != nil {
			return "", err
		}
		targetFile = outFile
	}
	if err := os.Chmod(targetFile, 0o755); err != nil {
		return "", err
	}
	return targetFile, nil
}

//gocyclo:ignore
func downloadBatchPlugins(pluginInfos []*plugin.Info, downloadInfo *resolver.BatchPluginDownloadInfo, showProgress bool) error {
	req, err := grab.NewRequest(PluginDir, downloadInfo.URL)
	if err != nil {
		return err
	}
	if downloadInfo.Checksum != "" {
		sum, decErr := hex.DecodeString(downloadInfo.Checksum)
		if decErr != nil {
			return fmt.Errorf("could not decode checksum: %w", decErr)
		}
		req.SetChecksum(sha256.New(), sum, true)
	}

	res := grab.DefaultClient.Do(req)
	if showProgress {
		showDownloadProgressBar("plugin-archive.tar.gz", res)
	}
	err = res.Err()
	if err != nil {
		return err
	}
	defer os.Remove(res.Filename)

	tgzFile, err := os.Open(res.Filename)
	if err != nil {
		return err
	}
	defer tgzFile.Close()

	gunzip, err := gzip.NewReader(tgzFile)
	if err != nil {
		return err
	}
	defer gunzip.Close()

	tarReader := tar.NewReader(gunzip)
	for {
		header, tarErr := tarReader.Next()
		if errors.Is(tarErr, io.EOF) {
			break
		}
		if tarErr != nil {
			return tarErr
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}

		outFileName := path.Join(PluginDir, header.Name)
		outDirName := path.Dir(outFileName)
		if err = os.MkdirAll(outDirName, 0o755); err != nil {
			return err
		}

		outFile, oErr := os.OpenFile(outFileName, os.O_CREATE|os.O_WRONLY, 0o755)
		if oErr != nil {
			return oErr
		}
		_, cErr := io.Copy(outFile, tarReader)
		_ = outFile.Close()
		if cErr != nil {
			return cErr
		}

		for _, pluginInfo := range pluginInfos {
			if strings.HasPrefix(path.Join(PluginDir, header.Name), pluginInfo.PluginPath) {
				pluginInfo.BinPath = outFileName
			}
		}

	}
	return nil
}
