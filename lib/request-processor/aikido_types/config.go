package aikido_types

type InitConfigData struct {
	LogLevel string `json:"log_level"`
	SAPI     string `json:"sapi"`
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
	Endpoints                   map[EndpointKey]EndpointData
	BlockedUserIds              map[string]bool
	IpsExcludedFromRateLimiting map[string]bool
}
