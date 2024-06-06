package cloud

import (
	. "main/aikido_types"
	. "main/globals"
	"time"
)

func GetAgentInfo() AgentInfo {
	return AgentInfo{
		DryMode:   !Config.Blocking,
		Hostname:  Machine.HostName,
		Version:   Version,
		IPAddress: Machine.IPAddress,
		OS: OsInfo{
			Name:    Machine.OS,
			Version: Machine.OSVersion,
		},
		Packages: make(map[string]string, 0),
		NodeEnv:  "",
	}
}

func GetTime() int64 {
	return time.Now().UnixMilli()
}
