package main

import (
	"fmt"
	"log"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	SocketPath                    string
	AllServices                   bool
	Service                       string
	MissingOk                     bool
	MissingFail                   bool
	WarningPercent                int
	CriticalPercent               int
	SessionWarningPercent         int
	SessionCriticalPercent        int
	BackendSessionWarningPercent  int
	BackendSessionCriticalPercent int
	MinWarningCount               int
	MinCriticalCount              int
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-go-haproxy-check",
			Short:    "plugin to check haproxy services",
			Keyspace: "sensu.io/plugins/sensu-go-haproxy-check/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "socket",
			Env:       "HAPROXY_SOCKET",
			Argument:  "socket",
			Shorthand: "S",
			Default:   "/var/run/haproxy.socket",
			Usage:     "Path to haproxy control socket",
			Value:     &plugin.SocketPath,
		},
		&sensu.PluginConfigOption{
			Path:      "service",
			Env:       "HAPROXY_SERVICE",
			Argument:  "service",
			Shorthand: "s",
			Default:   "",
			Usage:     "Service name to check",
			Value:     &plugin.Service,
		},
		&sensu.PluginConfigOption{
			Path:      "all_services",
			Env:       "HAPROXY_ALL_SERVICES",
			Argument:  "all-services",
			Shorthand: "A",
			Default:   false,
			Usage:     "Check all services",
			Value:     &plugin.AllServices,
		},
		&sensu.PluginConfigOption{
			Path:      "missing_ok",
			Env:       "HAPROXY_MISSING_OK",
			Argument:  "missing-ok",
			Shorthand: "m",
			Default:   false,
			Usage:     "Service missing is Ok",
			Value:     &plugin.MissingOk,
		},
		&sensu.PluginConfigOption{
			Path:      "missing_fail",
			Env:       "HAPROXY_MISSING_FAIL",
			Argument:  "missing-fail",
			Shorthand: "f",
			Default:   false,
			Usage:     "Service missing is Fail",
			Value:     &plugin.MissingFail,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	if len(plugin.Example) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--example or CHECK_EXAMPLE environment variable is required")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	log.Println("executing check with --example", plugin.Example)
	return sensu.CheckStateOK, nil
}
