package aikido_types

import (
	"regexp"

	"inet.af/netaddr"
)

type EnvironmentConfigData struct {
	SocketPath                string `json:"socket_path"`                  // '/run/aikido-{version}/aikido-{datetime}-{randint}.sock'
	SAPI                      string `json:"sapi"`                         // '{php-sapi}'
	TrustProxy                bool   `json:"trust_proxy"`                  // default: true
	LocalhostAllowedByDefault bool   `json:"localhost_allowed_by_default"` // default: true
	CollectApiSchema          bool   `json:"collect_api_schema"`           // default: true
}

type AikidoConfigData struct {
	Token                     string `json:"token"`                        // default: ''
	LogLevel                  string `json:"log_level"`                    // default: 'WARN'
	Blocking                  bool   `json:"blocking"`                     // default: false
	TrustProxy                bool   `json:"trust_proxy"`                  // default: true
	LocalhostAllowedByDefault bool   `json:"localhost_allowed_by_default"` // default: true
	CollectApiSchema          bool   `json:"collect_api_schema"`           // default: true
	DiskLogs                  bool   `json:"disk_logs"`                    // default: false
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

type IpBlockList struct {
	Description string
	IpSet       netaddr.IPSet
}

type CloudConfigData struct {
	ConfigUpdatedAt   int64
	Endpoints         map[EndpointKey]EndpointData
	BlockedUserIds    map[string]bool
	BypassedIps       map[string]bool
	BlockedIps        map[string]IpBlockList
	BlockedUserAgents *regexp.Regexp
	Block             int
}
