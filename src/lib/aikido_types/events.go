package aikido_types

type OsInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Hostname struct {
	URL  string `json:"hostname"`
	Port int64  `json:"port,omitempty"`
}

type Route struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type User struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	LastIpAddress string `json:"lastIpAddress"`
	FirstSeenAt   int64  `json:"firstSeenAt"`
	LastSeenAt    int64  `json:"lastSeenAt"`
}

type AgentInfo struct {
	DryMode                   bool              `json:"dryMode"`
	Hostname                  string            `json:"hostname"`
	Version                   string            `json:"version"`
	IPAddress                 string            `json:"ipAddress"`
	OS                        OsInfo            `json:"os"`
	Packages                  map[string]string `json:"packages"`
	PreventPrototypePollution bool              `json:"preventedPrototypePollution"`
	NodeEnv                   string            `json:"nodeEnv"`
}

type Started struct {
	Type  string    `json:"type"`
	Agent AgentInfo `json:"agent"`
	Time  int64     `json:"time"`
}

type Heartbeat struct {
	Type      string            `json:"type"`
	Stats     map[string]string `json:"stats"`
	Hostnames []Hostname        `json:"hostnames"`
	Routes    []Route           `json:"routes"`
	Users     []User            `json:"users"`
	Agent     AgentInfo         `json:"agent"`
	Time      int64             `json:"time"`
}
