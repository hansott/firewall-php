package grpc

import (
	"main/globals"
	"main/ipc/protos"
)

func incrementRequests() {
	globals.StatsData.StatsMutex.Lock()
	defer globals.StatsData.StatsMutex.Unlock()

	globals.StatsData.Requests += 1
}

func storeRoute(req *protos.RequestMetadata) {
	globals.RoutesMutex.Lock()
	defer globals.RoutesMutex.Unlock()

	if _, ok := globals.Routes[req.GetMethod()]; !ok {
		globals.Routes[req.GetMethod()] = make(map[string]int)
	}
	if _, ok := globals.Routes[req.GetMethod()][req.GetRoute()]; !ok {
		globals.Routes[req.GetMethod()][req.GetRoute()] = 0
	}
	globals.Routes[req.GetMethod()][req.GetRoute()]++
}
