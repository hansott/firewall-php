package globals

const (
	Version               = "0.2.0"
	ConfigFilePath        = "/opt/aikido/config.json"
	LogFilePath           = "/var/log/aikido.log"
	SocketPath            = "/var/aikido.sock"
	ConfigAPIMethod       = "GET"
	ConfigAPI             = "/api/runtime/config"
	ConfigUpdatedAtMethod = "GET"
	ConfigUpdatedAtAPI    = "/config"
	EventsAPIMethod       = "POST"
	EventsAPI             = "/api/runtime/events"
)
