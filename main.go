package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Showmax/go-fqdn"
	corev2 "github.com/sensu/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Url    []string
	Size   int64
	Scheme string
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
			Default:  1_000_000_000, // 1G
			Usage:    "Maximum aatabase Size",
			Value:    &plugin.Size,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "scheme",
			Argument:  "scheme",
			Shorthand: "s",
			Usage:     "Scheme to prepend metric",
			Value:     &plugin.Scheme,
		},
	}
)

func main() {
	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	return sensu.CheckStateOK, nil
}

func GetScheme() string {
	if len(plugin.Scheme) > 0 {
		return plugin.Scheme
	} else {
		realfqdn, err := fqdn.FqdnHostname()
		if err != nil {
			fmt.Printf("failed to get FQDN: %s", err)
		}
		return realfqdn
	}
}

func executeCheck(event *corev2.Event) (int, error) {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   plugin.Url,
		DialTimeout: 5 * time.Second,
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

	// print metrics
	fmt.Printf("etcd_dbsize{hostname=\"%s\"} %d %d\n", GetScheme(), status.DbSize, time.Now().Unix())

	if status.DbSize > plugin.Size {
		fmt.Printf("# Database exeeding set limit (%d): %d\n", plugin.Size, status.DbSize)
		return sensu.CheckStateCritical, nil
	}

	return sensu.CheckStateOK, nil
}
