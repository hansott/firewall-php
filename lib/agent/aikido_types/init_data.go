package aikido_types

type MachineData struct {
	HostName   string `json:"hostname"`
	DomainName string `json:"domainname"`
	OS         string `json:"os"`
	OSVersion  string `json:"os_version"`
	IPAddress  string `json:"ip_address"`
}

type EnvironmentConfigData struct {
	Token          string `json:"token,omitempty"`
	LogLevel       string `json:"log_level,omitempty"`
	Endpoint       string `json:"endpoint,omitempty"`
	ConfigEndpoint string `json:"config_endpoint,omitempty"`
	Blocking       bool   `json:"blocking,omitempty"`
}

type RateLimiting struct {
	Enabled        bool `json:"enabled"`
	MaxRequests    int  `json:"maxRequests"`
	WindowSizeInMS int  `json:"windowSizeInMS"`
}

type Endpoint struct {
	Method             string       `json:"method"`
	Route              string       `json:"route"`
	ForceProtectionOff bool         `json:"forceProtectionOff"`
	Graphql            interface{}  `json:"graphql"`
	AllowedIPAddresses []string     `json:"allowedIPAddresses"`
	RateLimiting       RateLimiting `json:"rateLimiting"`
}

type CloudConfigData struct {
	Success               bool       `json:"success"`
	ServiceId             int        `json:"serviceId"`
	ConfigUpdatedAt       int64      `json:"configUpdatedAt"`
	HeartbeatIntervalInMS int        `json:"heartbeatIntervalInMS"`
	Endpoints             []Endpoint `json:"endpoints"`
	BlockedUserIds        []string   `json:"blockedUserIds"`
	BypassedIps           []string   `json:"allowedIPAddresses"`
	ReceivedAnyStats      bool       `json:"receivedAnyStats"`
}

type CloudConfigUpdatedAt struct {
	ServiceId       int   `json:"serviceId"`
	ConfigUpdatedAt int64 `json:"configUpdatedAt"`
}
