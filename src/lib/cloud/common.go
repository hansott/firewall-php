package cloud

import (
	. "main/aikido_types"
	. "main/globals"
	"time"
)

func GetAgentInfo() AgentInfo {
	return AgentInfo{
		DryMode:   !InitData.Aikido.Blocking,
		Hostname:  InitData.Machine.HostName,
		Version:   InitData.Aikido.Version,
		IPAddress: InitData.Machine.IPAddress,
		OS: OsInfo{
			Name:    InitData.Machine.OS,
			Version: InitData.Machine.OSVersion,
		},
		Packages: make(map[string]string, 0),
		NodeEnv:  "",
	}
}

func GetTime() int64 {
	return time.Now().UnixMilli()
}
