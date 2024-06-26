package aikido_types

type MachineData struct {
	HostName   string `json:"hostname"`
	DomainName string `json:"domainname"`
	OS         string `json:"os"`
	OSVersion  string `json:"os_version"`
	IPAddress  string `json:"ip_address"`
}

type LocalConfigData struct {
	LogLevel string `json:"log_level"`
	Endpoint string `json:"endpoint"`
	Blocking bool   `json:"blocking"`
}

type CloudConfigData struct {
	Success               bool     `json:"success"`
	ServiceID             int      `json:"serviceId"`
	ConfigUpdatedAt       int64    `json:"configUpdatedAt"`
	HeartbeatIntervalInMS int64    `json:"heartbeatIntervalInMS"`
	Endpoints             []string `json:"endpoints"`
	BlockedUserIDs        []string `json:"blockedUserIds"`
	AllowedIPAddresses    []string `json:"allowedIPAddresses"`
}
