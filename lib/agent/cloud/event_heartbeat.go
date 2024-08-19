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

func GetUsersAndClear() []User {
	return make([]User, 0)
}

func GetStatsAndClear() Stats {
	globals.StatsData.StatsMutex.Lock()
	defer globals.StatsData.StatsMutex.Unlock()

	stats := Stats{
		Sinks:     make(map[string]MonitoredSinkStats),
		StartedAt: globals.StatsData.StartedAt,
		EndedAt:   GetTime(),
		Requests: Requests{
			Total:   globals.StatsData.Requests,
			Aborted: globals.StatsData.RequestsAborted,
			AttacksDetected: AttacksDetected{
				Total:   globals.StatsData.Attacks,
				Blocked: globals.StatsData.AttacksBlocked,
			},
		},
	}

	globals.StatsData.StartedAt = GetTime()
	globals.StatsData.Requests = 0
	globals.StatsData.RequestsAborted = 0
	globals.StatsData.Attacks = 0
	globals.StatsData.AttacksBlocked = 0

	return stats
}

func SendHeartbeatEvent() {
	heartbeatEvent := Heartbeat{
		Type:      "heartbeat",
		Agent:     GetAgentInfo(),
		Time:      GetTime(),
		Stats:     GetStatsAndClear(),
		Hostnames: GetHostnamesAndClear(),
		Routes:    GetRoutesAndClear(),
		Users:     GetUsersAndClear(),
	}

	response, err := SendCloudRequest(globals.EnvironmentConfig.Endpoint, globals.EventsAPI, globals.EventsAPIMethod, heartbeatEvent)
	if err != nil {
		log.Warn("Error in sending heartbeat event: ", err)
		return
	}
	StoreCloudConfig(response)
}
