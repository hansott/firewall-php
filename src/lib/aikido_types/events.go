package aikido_types

type OsInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
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
