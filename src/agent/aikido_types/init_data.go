package aikido_types

type MachineData struct {
	HostName   string `json:"hostname"`
	DomainName string `json:"domainname"`
	OS         string `json:"os"`
	OSVersion  string `json:"os_version"`
	IPAddress  string `json:"ip_address"`
}

type ConfigData struct {
	LogLevel string `json:"log_level"`
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
	Blocking bool   `json:"blocking"`
}
