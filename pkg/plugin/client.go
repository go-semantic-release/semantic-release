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

type PluginOpts struct {
	Type       string
	PluginName string
	Cmd        *exec.Cmd
}

var runningClientsMx sync.Mutex
var runningClients = make([]*plugin.Client, 0)

func KillAllPlugins() {
	runningClientsMx.Lock()
	defer runningClientsMx.Unlock()
	for _, c := range runningClients {
		c.Kill()
	}
}

func StartPlugin(opts *PluginOpts) (interface{}, error) {
	runningClientsMx.Lock()
	defer runningClientsMx.Unlock()
	logR, logW := io.Pipe()
	pluginLogger := log.New(os.Stderr, fmt.Sprintf("[%s]: ", opts.PluginName), 0)
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
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: {
				opts.Type: &GRPCWrapper{
					Type: opts.Type,
				},
			},
		},
		Cmd:              opts.Cmd,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           hclog.NewNullLogger(),
		Stderr:           logW,
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, err
	}
	raw, err := rpcClient.Dispense(opts.Type)
	if err != nil {
		client.Kill()
		return nil, err
	}
	runningClients = append(runningClients, client)
	return raw, nil
}
