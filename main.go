package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
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
	Debug            bool
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-go-haproxy-check",
			Short:    "plugin to check haproxy services",
			Keyspace: "sensu.io/plugins/sensu-go-haproxy-check/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "socket",
			Env:       "HAPROXY_SOCKET",
			Argument:  "socket",
			Shorthand: "S",
			Default:   "/var/run/haproxy.sock",
			Usage:     "Path to haproxy control socket",
			Value:     &plugin.SocketPath,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "service",
			Env:       "HAPROXY_SERVICE",
			Argument:  "service",
			Shorthand: "s",
			Default:   "",
			Usage:     "Service name to check",
			Value:     &plugin.Service,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "all_services",
			Env:       "HAPROXY_ALL_SERVICES",
			Argument:  "all-services",
			Shorthand: "A",
			Default:   false,
			Usage:     "Check all services",
			Value:     &plugin.AllServices,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "missing_ok",
			Env:       "HAPROXY_MISSING_OK",
			Argument:  "missing-ok",
			Shorthand: "m",
			Default:   false,
			Usage:     "Service missing is Ok",
			Value:     &plugin.MissingOk,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "missing_fail",
			Env:       "HAPROXY_MISSING_FAIL",
			Argument:  "missing-fail",
			Shorthand: "f",
			Default:   false,
			Usage:     "Service missing is Fail",
			Value:     &plugin.MissingFail,
		},
		&sensu.PluginConfigOption[float32]{
			Path:      "warning_percent",
			Env:       "HAPROXY_WARNING_PERCENT",
			Argument:  "warning-percent",
			Shorthand: "w",
			Default:   float32(50.0),
			Usage:     "Warning percent",
			Value:     &plugin.WarningPercent,
		},
		&sensu.PluginConfigOption[float32]{
			Path:      "critical_percent",
			Env:       "HAPROXY_CRITICAL_PERCENT",
			Argument:  "critical-percent",
			Shorthand: "c",
			Default:   float32(25.0),
			Usage:     "Critical percent",
			Value:     &plugin.CriticalPercent,
		},
		&sensu.PluginConfigOption[float32]{
			Path:      "session_warning_percent",
			Env:       "HAPROXY_SESSION_WARNING_PERCENT",
			Argument:  "session-warning-percent",
			Shorthand: "W",
			Default:   float32(75.0),
			Usage:     "Session Limit Warning percent",
			Value:     &plugin.SessionWarningPercent,
		},
		&sensu.PluginConfigOption[float32]{
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
		&sensu.PluginConfigOption[int]{
			Path:      "min_warning_count",
			Env:       "HAPROXY_MIN_WARNING_COUNT",
			Argument:  "min-warning-count",
			Shorthand: "M",
			Default:   0,
			Usage:     "Minimum server Warning count",
			Value:     &plugin.MinWarningCount,
		},
		&sensu.PluginConfigOption[int]{
			Path:      "min_critical_count",
			Env:       "HAPROXY_MIN_CRITICAL_COUNT",
			Argument:  "min-critical-count",
			Shorthand: "X",
			Default:   0,
			Usage:     "Minimum server Critical count",
			Value:     &plugin.MinCriticalCount,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "debug",
			Env:       "HAPROXY_DEBUG",
			Argument:  "debug",
			Shorthand: "d",
			Default:   false,
			Usage:     "output debugging data",
			Value:     &plugin.Debug,
		},
	}
)

func main() {
	useStdin := false
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error check stdin: %v\n", err)
		panic(err)
	}
	//Check the Mode bitmask for Named Pipe to indicate stdin is connected
	if fi.Mode()&os.ModeNamedPipe != 0 {
		useStdin = true
	}

	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, useStdin)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
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

func executeCheck(event *corev2.Event) (int, error) {
	stats, rawData, err := haproxy.GetStats(plugin.SocketPath)
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("Failed to get service stats: %w", err)
	}

	// Leave only selected services
	pxkeys := make([]string, 0)
	for key := range stats {
		if plugin.AllServices || key == plugin.Service {
			pxkeys = append(pxkeys, key)
			continue
		}

		delete(stats, key)
	}
	sort.Strings(pxkeys)

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
	for _, pxname := range pxkeys {
		stat := stats[pxname]
		newret, err2 := checkService(pxname, stat)
		err = multierr.Append(err, err2)
		if newret > ret {
			ret = newret
		}

		if plugin.Debug && (newret > sensu.CheckStateOK || err2 != nil) {
			b, _ := json.Marshal(&stat)
			log.Print(string(b))
		}
	}

	if plugin.Debug && (ret > sensu.CheckStateOK || err != nil) {
		log.Printf("Raw stat data\n---\n%s", string(rawData))
	}

	return ret, err
}

func checkService(pxname string, svc haproxy.StatService) (int, error) {
	servers := svc.Servers()
	backend, backendOk := svc[haproxy.Backend]

	var backendPtr *haproxy.StatLine
	if backendOk {
		backendPtr = &backend
	}

	// Ignore FRONTEND-only entries
	if len(servers) == 0 && plugin.AllServices {
		return sensu.CheckStateOK, nil
	}

	upCount := 0
	failedNames := make([]string, 0)

	for _, s := range servers {
		if s.IsUp(backendPtr) {
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
