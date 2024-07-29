package cloud

import (
	. "main/aikido_types"
	"main/globals"
	"main/log"
)

func GetHostnamesAndClear() []Hostname {
	globals.HostnamesMutex.Lock()
	defer globals.HostnamesMutex.Unlock()

	hostnames := make([]Hostname, len(globals.Hostnames))
	i := 0
	for domain := range globals.Hostnames {
		hostnames[i] = Hostname{URL: domain}
		i += 1
	}

	globals.Hostnames = map[string]bool{}
	return hostnames
}

func GetRoutesAndClear() []Route {
	globals.RoutesMutex.Lock()
	defer globals.RoutesMutex.Unlock()

	routes := make([]Route, 0)
	for method, routeMap := range globals.Routes {
		for route, hits := range routeMap {
			routes = append(routes, Route{Path: route, Method: method, Hits: int64(hits)})
		}
	}

	globals.Routes = map[string]map[string]int{}
	return routes
}

func SendHeartbeatEvent() {
	heartbeatEvent := Heartbeat{
		Type:      "heartbeat",
		Agent:     GetAgentInfo(),
		Time:      GetTime(),
		Stats:     make(map[string]string, 0),
		Hostnames: GetHostnamesAndClear(),
		Routes:    GetRoutesAndClear(),
		Users:     make([]User, 0),
	}

	response, err := SendEvent(globals.EventsAPI, globals.EventsAPIMethod, heartbeatEvent)
	if err != nil {
		log.Debug("Error in sending heartbeat event: ", err)
	}
	UpdateCloudConfig(response)
}
