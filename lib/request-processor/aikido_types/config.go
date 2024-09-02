package aikido_types

type EnvironmentConfigData struct {
	LogLevel                  string `json:"log_level"`
	SocketPath                string `json:"socket_path"`
	SAPI                      string `json:"sapi"`
	TrustProxy                bool   `json:"trust_proxy"`
	LocalhostAllowedByDefault bool   `json:"localhost_allowed_by_default"`
}

type RateLimiting struct {
	Enabled        bool
	MaxRequests    int
	WindowSizeInMS int
}

type EndpointData struct {
	ForceProtectionOff bool
	RateLimiting       RateLimiting
	AllowedIPAddresses map[string]bool
}

type EndpointKey struct {
	Method string
	Route  string
}

type CloudConfigData struct {
	Endpoints      map[EndpointKey]EndpointData
	BlockedUserIds map[string]bool
	BypassedIps    map[string]bool
	Block          int
}
