package globals

const (
	Version                  = "1.0.20"
	ConfigFilePath           = "/opt/aikido-" + Version + "/config.json"
	DevConfigFilePath        = "/opt/aikido-" + Version + "/config-dev.json"
	LogFilePath              = "/var/log/aikido-" + Version + "/aikido-agent.log"
	SocketPath               = "/run/aikido-" + Version + ".sock"
	ConfigAPIMethod          = "GET"
	ConfigAPI                = "/api/runtime/config"
	ConfigUpdatedAtMethod    = "GET"
	ConfigUpdatedAtAPI       = "/config"
	EventsAPIMethod          = "POST"
	EventsAPI                = "/api/runtime/events"
	MinHeartbeatIntervalInMS = 120000
)
