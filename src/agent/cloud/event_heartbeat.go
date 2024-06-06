package cloud

import (
	. "main/aikido_types"
	"main/globals"
	"main/log"
)

func GetHostnames() []Hostname {
	hostnames := make([]Hostname, len(globals.Hostnames))
	i := 0
	for domain := range globals.Hostnames {
		hostnames[i] = Hostname{URL: domain}
		i += 1
	}
	return hostnames
}

func SendHeartbeatEvent() {
	heartbeatEvent := Heartbeat{
		Type:      "heartbeat",
		Agent:     GetAgentInfo(),
		Time:      GetTime(),
		Stats:     make(map[string]string, 0),
		Hostnames: GetHostnames(),
		Routes:    make([]Route, 0),
		Users:     make([]User, 0),
	}

	_, err := SendEvent(globals.EventsAPI, globals.EventsAPIMethod, heartbeatEvent)
	if err != nil {
		log.Debug("Error in sending heartbeat event: ", err)
	}
}
