package resolver

import "github.com/go-semantic-release/semantic-release/v2/pkg/plugin"

type PluginDownloadInfo struct {
	URL      string
	Checksum string
	FileName string
	Version  string
}

type BatchPluginDownloadInfo struct {
	URL      string
	Checksum string
}

type Resolver interface {
	ResolvePlugin(*plugin.Info) (*PluginDownloadInfo, error)
	Names() []string
}

type BatchResolver interface {
	// BatchResolvePlugins resolves a list of plugins and returns a single URL to download all plugins as tgz archive.
	BatchResolvePlugins([]*plugin.Info) (*BatchPluginDownloadInfo, error)
}
