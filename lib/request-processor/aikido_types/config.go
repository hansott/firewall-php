package aikido_types

type EnvironmentConfigData struct {
	SocketPath                string `json:"socket_path"`                  // '/run/aikido-{version}/aikido-{datetime}-{randint}.sock'
	LogLevel                  string `json:"log_level"`                    // default: 'INFO'
	SAPI                      string `json:"sapi"`                         // '{php-sapi}'
	TrustProxy                bool   `json:"trust_proxy"`                  // default: true
	LocalhostAllowedByDefault bool   `json:"localhost_allowed_by_default"` // default: true
	CollectApiSchema          bool   `json:"collect_api_schema"`           // default: false
}

type AikidoConfigData struct {
	Token                     string `json:"token"`                        // default: ''
	Blocking                  bool   `json:"blocking"`                     // default: false
	TrustProxy                bool   `json:"trust_proxy"`                  // default: true
	LocalhostAllowedByDefault bool   `json:"localhost_allowed_by_default"` // default: true
	CollectApiSchema          bool   `json:"collect_api_schema"`           // default: false
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
