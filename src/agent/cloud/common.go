package cloud

import (
	"encoding/json"
	. "main/aikido_types"
	"main/globals"
	. "main/globals"
	"main/log"
	"time"
)

func GetAgentInfo() AgentInfo {
	return AgentInfo{
		DryMode:   !LocalConfig.Blocking,
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

func ApplyCloudConfig() {
	log.Infof("Applying new cloud config: %v", globals.CloudConfig)
	HeartBeatTicker.Reset(time.Duration(globals.CloudConfig.HeartbeatIntervalInMS) * time.Millisecond)
}

func UpdateCloudConfig(response []byte) bool {
	globals.ConfigMutex.Lock()
	defer globals.ConfigMutex.Unlock()

	tempCloudConfig := CloudConfigData{}
	err := json.Unmarshal(response, &tempCloudConfig)
	if err != nil {
		return false
	}
	if tempCloudConfig.ConfigUpdatedAt == globals.CloudConfig.ConfigUpdatedAt {
		return true
	}
	globals.CloudConfig = tempCloudConfig
	ApplyCloudConfig()
	return true
}
