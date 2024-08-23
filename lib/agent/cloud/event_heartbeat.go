package cloud

import (
	. "main/aikido_types"
	"main/globals"
	"main/log"
	"main/utils"
)

func GetHostnamesAndClear() []Hostname {
	globals.HostnamesMutex.Lock()
	defer globals.HostnamesMutex.Unlock()

	var hostnames []Hostname
	for domain := range globals.Hostnames {
		for port := range globals.Hostnames[domain] {
			hostnames = append(hostnames, Hostname{URL: domain, Port: int64(port)})
		}
	}

	globals.Hostnames = make(map[string]map[int]bool)
	return hostnames
}

func GetRoutesAndClear() []Route {
	globals.RoutesMutex.Lock()
	defer globals.RoutesMutex.Unlock()

	var routes []Route
	for method, routeMap := range globals.Routes {
		for route, hits := range routeMap {
			routes = append(routes, Route{Path: route, Method: method, Hits: int64(hits)})
		}
	}

	globals.Routes = make(map[string]map[string]int)
	return routes
}

func GetUsersAndClear() []User {
	globals.UsersMutex.Lock()
	defer globals.UsersMutex.Unlock()

	var users []User
	for _, user := range globals.Users {
		users = append(users, user)
	}

	globals.Users = make(map[string]User)
	return users
}

func GetStatsAndClear() Stats {
	globals.StatsData.StatsMutex.Lock()
	defer globals.StatsData.StatsMutex.Unlock()

	stats := Stats{
		Sinks:     make(map[string]MonitoredSinkStats),
		StartedAt: globals.StatsData.StartedAt,
		EndedAt:   utils.GetTime(),
		Requests: Requests{
			Total:   globals.StatsData.Requests,
			Aborted: globals.StatsData.RequestsAborted,
			AttacksDetected: AttacksDetected{
				Total:   globals.StatsData.Attacks,
				Blocked: globals.StatsData.AttacksBlocked,
			},
		},
	}

	globals.StatsData.StartedAt = utils.GetTime()
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
		Time:      utils.GetTime(),
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
