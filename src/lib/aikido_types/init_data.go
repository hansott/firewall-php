package aikido_types

type MachineInitData struct {
	HostName    string `json:"hostname"`
	DomaineName string `json:"domainname"`
	OS          string `json:"os"`
	OSVersion   string `json:"os_version"`
	IPAddress   string `json:"ip_address"`
}

type AikidoInitData struct {
	Version  string `json:"version"`
	LogLevel string `json:"log_level"`
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
	Blocking bool   `json:"blocking"`
}

type InitDataType struct {
	Machine MachineInitData `json:"machine"`
	Aikido  AikidoInitData  `json:"aikido"`
}
