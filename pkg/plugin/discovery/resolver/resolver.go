package resolver

import "github.com/go-semantic-release/semantic-release/v2/pkg/plugin"

type PluginDownloadInfo struct {
	URL      string
	Checksum string
	FileName string
	Version  string
}

type Resolver interface {
	ResolvePlugin(*plugin.Info) (*PluginDownloadInfo, error)
	Names() []string
}
