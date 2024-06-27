package globals

import (
	. "main/aikido_types"
	"sync"
)

// Aikido cloud token, passed by the extension via gRPC
var Token string

// Aikido log level token, passed by the extension via gRPC
var LogLevel string

// Local config loaded from LocalConfigPath, that contains info about endpoint, log_level, ...
var LocalConfig LocalConfigData

// Cloud config that is obtain as a result from sending events to cloud or pulling the config when it changes
var CloudConfig CloudConfigData

// Config mutex used to sync access to configuration data across the multiple go routines that we run in parallel
var ConfigMutex sync.Mutex

// Data about the current machine, computed at init
var Machine MachineData

// List of outgoing hostnames collect from the extensions
var Hostnames = map[string]bool{}

// Hostnames mutex used to sync access to hostnames data across the go routines that populate and read the array
var HostnamesMutex sync.Mutex
