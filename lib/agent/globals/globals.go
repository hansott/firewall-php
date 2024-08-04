package globals

import (
	. "main/aikido_types"
	"sync"
)

// Local config loaded from EnvironmentConfigPath, that contains info about endpoint, log_level, token, ...
var EnvironmentConfig EnvironmentConfigData

// Cloud config that is obtain as a result from sending events to cloud or pulling the config when it changes
var CloudConfig CloudConfigData

// Config mutex used to sync access to configuration data across the multiple go routines that we run in parallel
var CloudConfigMutex sync.Mutex

// Data about the current machine, computed at init
var Machine MachineData

// List of outgoing hostnames collect from the extensions
var Hostnames = map[string]bool{}

// Hostnames mutex used to sync access to hostnames data across the go routines
var HostnamesMutex sync.Mutex

// List of routes and their methods and count of calls collect from the extensions
// [method][route] = hits
var Routes = map[string]map[string]int{}

// Routes mutex used to sync access to routes data across the go routines
var RoutesMutex sync.Mutex

// Global stats data, including mutex used to sync access to stats data across the go routines
var StatsData StatsDataType

// Rate limiting map, which holds the current rate limiting state for each configured route
var RateLimitingMap = make(map[RateLimitingKey]*RateLimitingValue)

// Rate limiting mutex used to sync access across the go routines
var RateLimitingMutex sync.Mutex
