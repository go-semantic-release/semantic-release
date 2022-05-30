package plugin

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

var runningClientsMx sync.Mutex
var runningClients = make([]*plugin.Client, 0)

func KillAllPlugins() {
	runningClientsMx.Lock()
	defer runningClientsMx.Unlock()
	for _, c := range runningClients {
		c.Kill()
	}
}

func StartPlugin(pluginInfo *PluginInfo) (interface{}, error) {
	runningClientsMx.Lock()
	defer runningClientsMx.Unlock()
	logR, logW := io.Pipe()
	pluginLogger := log.New(os.Stderr, fmt.Sprintf("[%s]: ", pluginInfo.ShortNormalizedName), 0)
	go func() {
		logLineScanner := bufio.NewScanner(logR)
		for logLineScanner.Scan() {
			line := logLineScanner.Text()
			// skip JSON logging
			if strings.HasPrefix(line, "{") {
				continue
			}
			pluginLogger.Println(line)
		}
	}()

	cmd := exec.Command(pluginInfo.BinPath)
	cmd.SysProcAttr = GetSysProcAttr()

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: {
				pluginInfo.Type: &GRPCWrapper{
					Type: pluginInfo.Type,
				},
			},
		},
		Cmd:              cmd,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           hclog.NewNullLogger(),
		Stderr:           logW,
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, err
	}
	raw, err := rpcClient.Dispense(pluginInfo.Type)
	if err != nil {
		client.Kill()
		return nil, err
	}
	runningClients = append(runningClients, client)
	return raw, nil
}
