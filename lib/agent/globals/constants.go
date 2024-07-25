package globals

const (
	Version                  = "1.0.25"
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
