package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
	"go.uber.org/multierr"

	"github.com/sardinasystems/sensu-go-haproxy-check/haproxy"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	SocketPath             string
	AllServices            bool
	Service                string
	MissingOk              bool
	MissingFail            bool
	WarningPercent         float32
	CriticalPercent        float32
	SessionWarningPercent  float32
	SessionCriticalPercent float32
	// BackendSessionWarningPercent  float32
	// BackendSessionCriticalPercent float32
	MinWarningCount  int
	MinCriticalCount int
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
		{
			Path:      "socket",
			Env:       "HAPROXY_SOCKET",
			Argument:  "socket",
			Shorthand: "S",
			Default:   "/var/run/haproxy.sock",
			Usage:     "Path to haproxy control socket",
			Value:     &plugin.SocketPath,
		},
		{
			Path:      "service",
			Env:       "HAPROXY_SERVICE",
			Argument:  "service",
			Shorthand: "s",
			Default:   "",
			Usage:     "Service name to check",
			Value:     &plugin.Service,
		},
		{
			Path:      "all_services",
			Env:       "HAPROXY_ALL_SERVICES",
			Argument:  "all-services",
			Shorthand: "A",
			Default:   false,
			Usage:     "Check all services",
			Value:     &plugin.AllServices,
		},
		{
			Path:      "missing_ok",
			Env:       "HAPROXY_MISSING_OK",
			Argument:  "missing-ok",
			Shorthand: "m",
			Default:   false,
			Usage:     "Service missing is Ok",
			Value:     &plugin.MissingOk,
		},
		{
			Path:      "missing_fail",
			Env:       "HAPROXY_MISSING_FAIL",
			Argument:  "missing-fail",
			Shorthand: "f",
			Default:   false,
			Usage:     "Service missing is Fail",
			Value:     &plugin.MissingFail,
		},
		{
			Path:      "warning_percent",
			Env:       "HAPROXY_WARNING_PERCENT",
			Argument:  "warning-percent",
			Shorthand: "w",
			Default:   float32(50.0),
			Usage:     "Warning percent",
			Value:     &plugin.WarningPercent,
		},
		{
			Path:      "critical_percent",
			Env:       "HAPROXY_CRITICAL_PERCENT",
			Argument:  "critical-percent",
			Shorthand: "c",
			Default:   float32(25.0),
			Usage:     "Critical percent",
			Value:     &plugin.CriticalPercent,
		},
		{
			Path:      "session_warning_percent",
			Env:       "HAPROXY_SESSION_WARNING_PERCENT",
			Argument:  "session-warning-percent",
			Shorthand: "W",
			Default:   float32(75.0),
			Usage:     "Session Limit Warning percent",
			Value:     &plugin.SessionWarningPercent,
		},
		{
			Path:      "session_critical_percent",
			Env:       "HAPROXY_SESSION_CRITICAL_PERCENT",
			Argument:  "session-critical-percent",
			Shorthand: "C",
			Default:   float32(90.0),
			Usage:     "Session Limit Critical percent",
			Value:     &plugin.SessionCriticalPercent,
		},
		// {
		// 	Path:      "backend_session_warning_percent",
		// 	Env:       "HAPROXY_BACKEND_SESSION_WARNING_PERCENT",
		// 	Argument:  "backend-session-warning-percent",
		// 	Shorthand: "b",
		// 	Default:   0,
		// 	Usage:     "Per Backend Session Limit Warning percent",
		// 	Value:     &plugin.BackendSessionWarningPercent,
		// },
		// {
		// 	Path:      "backend_session_critical_percent",
		// 	Env:       "HAPROXY_BACKEND_SESSION_CRITICAL_PERCENT",
		// 	Argument:  "backend-session-critical-percent",
		// 	Shorthand: "B",
		// 	Default:   0,
		// 	Usage:     "Per Backend Session Limit Critical percent",
		// 	Value:     &plugin.BackendSessionCriticalPercent,
		// },
		{
			Path:      "min_warning_count",
			Env:       "HAPROXY_MIN_WARNING_COUNT",
			Argument:  "min-warning-count",
			Shorthand: "M",
			Default:   0,
			Usage:     "Minimum server Warning count",
			Value:     &plugin.MinWarningCount,
		},
		{
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
	} else if fi.Mode()&os.ModeSocket == 0 {
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
	stats, err := haproxy.GetStats(plugin.SocketPath)
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("Failed to get service stats: %w", err)
	}

	// Leave only selected services
	for key := range stats {
		if plugin.AllServices || key == plugin.Service {
			continue
		}

		delete(stats, key)
	}

	// No services
	if len(stats) == 0 {
		log.Printf("No service: %s", plugin.Service)
		if plugin.MissingFail {
			return sensu.CheckStateCritical, nil
		} else if plugin.MissingOk {
			return sensu.CheckStateOK, nil
		}

		return sensu.CheckStateUnknown, nil
	}

	ret := sensu.CheckStateOK
	err = nil
	for pxname, stat := range stats {
		newret, err2 := checkService(pxname, stat)
		err = multierr.Append(err, err2)
		if newret > ret {
			ret = newret
		}
	}

	return ret, err
}

func checkService(pxname string, svc haproxy.StatService) (int, error) {
	servers := svc.Servers()

	// Ignore FRONTEND-only entries
	if len(servers) == 0 && plugin.AllServices {
		return sensu.CheckStateOK, nil
	}

	upCount := 0
	failedNames := make([]string, 0)

	for _, s := range servers {
		if s.IsUp() {
			upCount++
		} else {
			failedNames = append(failedNames, s.LogName())
		}
	}

	upPercent := 100.0 * float32(upCount) / float32(len(servers))

	criticalSesions := servers.Filter(func(s haproxy.StatLine) bool {
		return s.Slim > 0 && s.SessionLimitPercentage() > plugin.SessionCriticalPercent
	})

	warningSesions := servers.Filter(func(s haproxy.StatLine) bool {
		return s.Slim > 0 && s.SessionLimitPercentage() > plugin.SessionWarningPercent
	})

	log.Printf("UP: %.0f%% of #%d %s services", upPercent, len(servers), pxname)
	if len(failedNames) > 0 {
		log.Printf("DOWN: %s", strings.Join(failedNames, ", "))
	}

	if len(servers) < plugin.MinCriticalCount {
		return sensu.CheckStateCritical, nil
	} else if upPercent < plugin.CriticalPercent {
		return sensu.CheckStateCritical, nil
	} else if len(criticalSesions) > 0 {
		log.Printf("Active sessions critical:")
		for _, s := range criticalSesions {
			log.Printf("\t%s: %d of %d (%.0f%%) sessions", s.LogName(), s.Scur, s.Slim, s.SessionLimitPercentage())
		}
		return sensu.CheckStateCritical, nil
	}

	if len(servers) < plugin.MinWarningCount {
		return sensu.CheckStateWarning, nil
	} else if upPercent < plugin.WarningPercent {
		return sensu.CheckStateWarning, nil
	} else if len(criticalSesions) > 0 {
		log.Printf("Active sessions warning:")
		for _, s := range warningSesions {
			log.Printf("\t%s: %d of %d (%.0f%%) sessions", s.LogName(), s.Scur, s.Slim, s.SessionLimitPercentage())
		}
		return sensu.CheckStateWarning, nil
	}

	return sensu.CheckStateOK, nil
}
