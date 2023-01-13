package github

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/discovery/resolver"
	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"
)

type GitHubResolver struct {
	ghClient *github.Client
}

func NewResolver() *GitHubResolver {
	var tc *http.Client
	if ghToken := os.Getenv("GITHUB_TOKEN"); ghToken != "" {
		tc = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghToken}))
	}
	return &GitHubResolver{
		ghClient: github.NewClient(tc),
	}
}

func (g *GitHubResolver) githubReleaseToDownloadInfo(repoOwner, repoName string, release *github.RepositoryRelease) (*resolver.PluginDownloadInfo, error) {
	var checksumAsset *github.ReleaseAsset
	var pluginAsset *github.ReleaseAsset
	osArchRe := regexp.MustCompile("(?i)" + runtime.GOOS + "(_|-)" + runtime.GOARCH)
	for _, asset := range release.Assets {
		assetName := asset.GetName()
		if checksumAsset == nil && asset.GetSize() <= 4096 && strings.Contains(strings.ToLower(assetName), "checksum") {
			checksumAsset = asset
		}
		if pluginAsset == nil && osArchRe.MatchString(assetName) {
			pluginAsset = asset
		}
	}

	if pluginAsset == nil {
		return nil, fmt.Errorf("no matching plugin binary was found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	foundChecksum := ""
	if checksumAsset != nil {
		checksumDownload, _, err := g.ghClient.Repositories.DownloadReleaseAsset(context.Background(), repoOwner, repoName, checksumAsset.GetID(), http.DefaultClient)
		if err != nil {
			return nil, err
		}
		checksumData, err := ioutil.ReadAll(checksumDownload)
		checksumDownload.Close()
		if err != nil {
			return nil, err
		}
		for _, l := range strings.Split(string(checksumData), "\n") {
			sl := strings.Split(l, " ")
			if len(sl) < 3 {
				continue
			}
			if sl[2] == pluginAsset.GetName() {
				foundChecksum = sl[0]
			}
		}
	}

	return &resolver.PluginDownloadInfo{
		URL:      pluginAsset.GetBrowserDownloadURL(),
		Checksum: foundChecksum,
		FileName: pluginAsset.GetName(),
		Version:  strings.TrimLeft(release.GetTagName(), "v"),
	}, nil
}

type ghRelease struct {
	version *semver.Version
	release *github.RepositoryRelease
}

type ghReleases []*ghRelease

func (gr ghReleases) Len() int           { return len(gr) }
func (gr ghReleases) Less(i, j int) bool { return gr[j].version.LessThan(gr[i].version) }
func (gr ghReleases) Swap(i, j int)      { gr[i], gr[j] = gr[j], gr[i] }

func (g *GitHubResolver) getAllValidGitHubReleases(repoOwner, repoName string) (ghReleases, error) {
	ret := make(ghReleases, 0)
	opts := &github.ListOptions{Page: 1, PerPage: 100}
	for {
		releases, resp, err := g.ghClient.Repositories.ListReleases(context.Background(), repoOwner, repoName, opts)
		if err != nil {
			return nil, err
		}
		for _, release := range releases {
			if release.GetDraft() {
				continue
			}

			if semverTag, err := semver.NewVersion(release.GetTagName()); err == nil {
				ret = append(ret, &ghRelease{version: semverTag, release: release})
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	sort.Sort(ret)
	return ret, nil
}

func (g *GitHubResolver) ResolvePlugin(pluginInfo *plugin.PluginInfo) (*resolver.PluginDownloadInfo, error) {
	if pluginInfo.RepoSlug == "" {
		pluginInfo.RepoSlug = knownPlugins[pluginInfo.ShortNormalizedName]
	}
	if pluginInfo.RepoSlug == "" {
		return nil, fmt.Errorf("repo slug not found")
	}
	repoOwner, repoName, _ := strings.Cut(pluginInfo.RepoSlug, "/")

	var foundRelease *github.RepositoryRelease
	if pluginInfo.Constraint == nil {
		release, _, err := g.ghClient.Repositories.GetLatestRelease(context.Background(), repoOwner, repoName)
		if err != nil {
			return nil, err
		}
		foundRelease = release
	} else {
		validReleases, err := g.getAllValidGitHubReleases(repoOwner, repoName)
		if err != nil {
			return nil, err
		}
		for _, release := range validReleases {
			if pluginInfo.Constraint.Check(release.version) {
				foundRelease = release.release
				break
			}
		}
	}

	if foundRelease == nil {
		return nil, fmt.Errorf("no matching release was found")
	}

	di, err := g.githubReleaseToDownloadInfo(repoOwner, repoName, foundRelease)
	return di, err
}

func (g *GitHubResolver) Names() []string {
	return []string{"github", "gh"}
}
