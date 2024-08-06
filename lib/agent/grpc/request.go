package grpc

import (
	. "main/aikido_types"
	"main/globals"
	"main/ipc/protos"
	"main/log"
)

func storeStats() {
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

func updateRateLimitingStatus(req *protos.RequestMetadata) {
	globals.RateLimitingMutex.Lock()
	defer globals.RateLimitingMutex.Unlock()

	rateLimitingData, exists := globals.RateLimitingMap[RateLimitingKey{Method: req.GetMethod(), Route: req.GetRoute()}]
	if !exists {
		return
	}

	rateLimitingData.Status.TotalNumberOfRequests += 1
	rateLimitingData.Status.NumberOfRequestsPerWindow.IncrementLast()
}

func getRequestStatus(req *protos.RequestMetadata) *protos.RequestStatus {
	globals.RateLimitingMutex.Lock()
	defer globals.RateLimitingMutex.Unlock()

	forwardToServer := true
	rateLimitingData, exists := globals.RateLimitingMap[RateLimitingKey{Method: req.GetMethod(), Route: req.GetRoute()}]
	if exists && rateLimitingData.Status.TotalNumberOfRequests >= rateLimitingData.Config.MaxRequests {
		log.Infof("Rate limited request for (%v) - status (%v)", req, rateLimitingData)
		forwardToServer = false
	}

	return &protos.RequestStatus{
		ForwardToServer: forwardToServer,
	}
}

func getCloudConfig() *protos.CloudConfig {
	globals.CloudConfigMutex.Lock()
	defer globals.CloudConfigMutex.Unlock()

	cloudConfig := &protos.CloudConfig{
		BlockedUserIds:     globals.CloudConfig.BlockedUserIds,
		AllowedIPAddresses: globals.CloudConfig.AllowedIPAddresses,
	}

	for _, endpoint := range globals.CloudConfig.Endpoints {
		cloudConfig.Endpoints = append(cloudConfig.Endpoints, &protos.Endpoint{
			Method:             endpoint.Method,
			Route:              endpoint.Route,
			ForceProtectionOff: endpoint.ForceProtectionOff,
			RateLimiting: &protos.RateLimiting{
				Enabled: endpoint.RateLimiting.Enabled,
			},
		})
	}

	return cloudConfig
}
