package globals

import (
	. "main/aikido_types"
	"sync"
)

var LocalConfig LocalConfigData
var CloudConfig CloudConfigData
var ConfigMutex sync.Mutex

var Machine MachineData

var Hostnames = map[string]bool{}
var HostnamesMutex sync.Mutex
