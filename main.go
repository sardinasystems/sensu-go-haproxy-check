package main

import (
	"fmt"
	"os"
	"path/filepath"

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
		&sensu.PluginConfigOption{
			Path:      "warning_percent",
			Env:       "HAPROXY_WARNING_PERCENT",
			Argument:  "warning-percent",
			Shorthand: "w",
			Default:   50,
			Usage:     "Warning percent",
			Value:     &plugin.WarningPercent,
		},
		&sensu.PluginConfigOption{
			Path:      "critical_percent",
			Env:       "HAPROXY_CRITICAL_PERCENT",
			Argument:  "critical-percent",
			Shorthand: "c",
			Default:   25,
			Usage:     "Critical percent",
			Value:     &plugin.CriticalPercent,
		},
		&sensu.PluginConfigOption{
			Path:      "session_warning_percent",
			Env:       "HAPROXY_SESSION_WARNING_PERCENT",
			Argument:  "session-warning-percent",
			Shorthand: "W",
			Default:   75,
			Usage:     "Session Limit Warning percent",
			Value:     &plugin.SessionWarningPercent,
		},
		&sensu.PluginConfigOption{
			Path:      "session_critical_percent",
			Env:       "HAPROXY_SESSION_CRITICAL_PERCENT",
			Argument:  "session-critical-percent",
			Shorthand: "C",
			Default:   90,
			Usage:     "Session Limit Critical percent",
			Value:     &plugin.SessionCriticalPercent,
		},
		&sensu.PluginConfigOption{
			Path:      "backend_session_warning_percent",
			Env:       "HAPROXY_BACKEND_SESSION_WARNING_PERCENT",
			Argument:  "backend-session-warning-percent",
			Shorthand: "b",
			Default:   0,
			Usage:     "Per Backend Session Limit Warning percent",
			Value:     &plugin.BackendSessionWarningPercent,
		},
		&sensu.PluginConfigOption{
			Path:      "backend_session_critical_percent",
			Env:       "HAPROXY_BACKEND_SESSION_CRITICAL_PERCENT",
			Argument:  "backend-session-critical-percent",
			Shorthand: "B",
			Default:   0,
			Usage:     "Per Backend Session Limit Critical percent",
			Value:     &plugin.BackendSessionCriticalPercent,
		},
		&sensu.PluginConfigOption{
			Path:      "min_warning_count",
			Env:       "HAPROXY_MIN_WARNING_COUNT",
			Argument:  "min-warning-count",
			Shorthand: "M",
			Default:   0,
			Usage:     "Minimum server Warning count",
			Value:     &plugin.MinWarningCount,
		},
		&sensu.PluginConfigOption{
			Path:      "min_critical_count",
			Env:       "HAPROXY_MIN_CRITICAL_COUNT",
			Argument:  "min-critical-count",
			Shorthand: "X",
			Default:   0,
			Usage:     "Minimum server Critical count",
			Value:     &plugin.MinCriticalCount,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	path, err := filepath.Abs(plugin.SocketPath)
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("--socket error: %w", err)
	}

	fi, err := os.Lstat(path)
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("--socket error: %w", err)
	} else if fi.Mode() != os.ModeSocket {
		return sensu.CheckStateUnknown, fmt.Errorf("--socket: %s is not socket: %v", path, fi.Mode())
	}
	plugin.SocketPath = path

	if plugin.Service == "" && !plugin.AllServices {
		return sensu.CheckStateWarning, fmt.Errorf("--service or --all-services are required")
	} else if plugin.Service != "" && plugin.AllServices {
		return sensu.CheckStateWarning, fmt.Errorf("Only one --service or --all-services should be used")
	}

	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {

	stats, err := haproxy.GetStat(plugin.SocketPath)
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("Failed to get service stats: %w", err)
	}

	return sensu.CheckStateOK, nil
}
