package globals

const (
	Version                  = "2.0.7"
	ConfigFilePath           = "/opt/aikido/config.json"
	DevConfigFilePath        = "/opt/aikido/config-dev.json"
	LogFilePath              = "/var/log/aikido/aikido_agent.log"
	SocketPath               = "/run/aikido.sock"
	ConfigAPIMethod          = "GET"
	ConfigAPI                = "/api/runtime/config"
	ConfigUpdatedAtMethod    = "GET"
	ConfigUpdatedAtAPI       = "/config"
	EventsAPIMethod          = "POST"
	EventsAPI                = "/api/runtime/events"
	MinHeartbeatIntervalInMS = 120000
)
