package globals

const (
	Version                  = "1.15.0"
	ConfigFilePath           = "/opt/aikido/config.json"
	DevConfigFilePath        = "/opt/aikido/config-dev.json"
	LogFilePath              = "/opt/aikido/aikido.log"
	SocketPath               = "/run/aikido.sock"
	ConfigAPIMethod          = "GET"
	ConfigAPI                = "/api/runtime/config"
	ConfigUpdatedAtMethod    = "GET"
	ConfigUpdatedAtAPI       = "/config"
	EventsAPIMethod          = "POST"
	EventsAPI                = "/api/runtime/events"
	MinHeartbeatIntervalInMS = 120000
)
