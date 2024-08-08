package main

import (
	"main/grpc"
	"main/log"
	"main/utils"
	"time"
)

func OnRequestInit(data map[string]interface{}) string {
	method := utils.MustGetFromMap[string](data, "method")
	route := utils.MustGetFromMap[string](data, "route")

	if method == "" || route == "" {
		return "{\"status\": \"ok\"}"
	}

	log.Info("[RINIT] Got request metadata: ", method, " ", route)

	route = utils.BuildRouteFromURL(route)

	if grpc.IsRequestMonitoredForRateLimiting(method, route) {
		// If request is monitored for rate limiting, do a sync call via gRPC to see if the request should be aborded or not
		if !grpc.OnRequestInit(method, route, 10*time.Millisecond) {
			return "{\"action\": \"exit\", \"message\": \"This request was rate limited by Aikido Security!\", \"response_code\": 429}"
		}
	}

	return "{\"status\": \"ok\"}"
}

func OnRequestShutdown(data map[string]interface{}) string {
	method := utils.MustGetFromMap[string](data, "method")
	route := utils.MustGetFromMap[string](data, "route")
	status_code := int(utils.MustGetFromMap[float64](data, "status_code"))

	if method == "" || route == "" || status_code == 0 {
		return "{\"status\": \"ok\"}"
	}

	log.Info("[RSHUTDOWN] Got request metadata: ", method, " ", route, " ", status_code)

	route = utils.BuildRouteFromURL(route)

	if !utils.ShouldDiscoverRoute(status_code, route, method) {
		return "{\"status\": \"ok\"}"
	}

	go grpc.OnRequestShutdown(method, route, status_code, 10*time.Millisecond)

	return "{\"status\": \"ok\"}"
}
