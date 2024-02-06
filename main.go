package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"time"

	corev2 "github.com/sensu/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Url           []string
	Size          int64
	CertFile      string
	KeyFile       string
	TrustedCAFile string
	Timeout       int64
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-etcd",
			Short:    "Check etcd",
			Keyspace: "sensu.io/plugins/check-etcd/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.SlicePluginConfigOption[string]{
			Path:     "url",
			Argument: "url",
			Default:  []string{"http://127.0.0.1:2379"},
			Usage:    "Url of etcd instance(s)",
			Value:    &plugin.Url,
		},
		&sensu.PluginConfigOption[int64]{
			Path:     "size",
			Argument: "size",
			Default:  1_500_000_000, // Alarm at 1.5G, default DB is set to 2G
			Usage:    "Maximum aatabase Size",
			Value:    &plugin.Size,
		},
		&sensu.PluginConfigOption[string]{
			Path:     "cert-file",
			Argument: "cert-file",
			Usage:    "Path to the cert",
			Value:    &plugin.CertFile,
		},
		&sensu.PluginConfigOption[string]{
			Path:     "key-file",
			Argument: "key-file",
			Usage:    "Path to the key",
			Value:    &plugin.KeyFile,
		},
		&sensu.PluginConfigOption[string]{
			Path:     "trusted-ca-file",
			Argument: "trusted-ca-file",
			Usage:    "Path to the CA file",
			Value:    &plugin.TrustedCAFile,
		},
		&sensu.PluginConfigOption[int64]{
			Path:     "timeout",
			Argument: "timeout",
			Usage:    "Request timeout",
			Default:  5,
			Value:    &plugin.Timeout,
		},
	}
)

func main() {
	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {

	if _, err := os.Stat(plugin.CertFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("could not load certificate(%s): %v", plugin.CertFile, err)
		return sensu.CheckStateCritical, nil
	}

	if _, err := os.Stat(plugin.KeyFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("could not load certificate key(%s): %v", plugin.KeyFile, err)
		return sensu.CheckStateCritical, nil
	}

	if _, err := os.Stat(plugin.TrustedCAFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("could not load CA(%s): %v", plugin.TrustedCAFile, err)
		return sensu.CheckStateCritical, nil
	}

	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	tlsConfig := &tls.Config{}
	if len(plugin.CertFile) > 0 && len(plugin.KeyFile) > 0 && len(plugin.TrustedCAFile) > 0 {
		tlsInfo := transport.TLSInfo{
			CertFile:      plugin.CertFile,
			KeyFile:       plugin.KeyFile,
			TrustedCAFile: plugin.TrustedCAFile,
		}
		tlsConfig, _ = tlsInfo.ClientConfig()
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   plugin.Url,
		DialTimeout: time.Duration(plugin.Timeout) * time.Second,
		TLS:         tlsConfig,
	})

	if err != nil {
		fmt.Printf("could not connect: %s", err)
		return sensu.CheckStateCritical, nil
	}
	defer cli.Close()

	status, err := cli.Status(context.Background(), plugin.Url[0])
	if err != nil {
		fmt.Printf("failed to get status: %s", err)
		return sensu.CheckStateCritical, nil
	}

	if status.DbSize > plugin.Size {
		fmt.Printf("Database exeeding set limit (%d): %d\n", plugin.Size, status.DbSize)
		return sensu.CheckStateCritical, nil
	}
	fmt.Printf("Database is within size limit (%d): %d\n", plugin.Size, status.DbSize)
	return sensu.CheckStateOK, nil
}
