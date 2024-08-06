package main

import (
	"main/grpc"
	"main/log"
	"main/utils"
	"time"
)

func OnRequestMetadata(data map[string]interface{}) string {
	method := utils.MustGetFromMap[string](data, "method")
	route := utils.MustGetFromMap[string](data, "route")

	if method == "" || route == "" {
		return "{\"status\": \"ok\"}"
	}

	log.Info("Got request metadata: ", method, " ", route)

	route = utils.BuildRouteFromURL(route)

	if grpc.IsRequestMonitoredForRateLimiting(method, route) {
		// If request is monitored for rate limiting, do a sync call via gRPC to see if the request should be aborded or not
		if !grpc.OnRequest(method, route, 10*time.Millisecond) {
			return "{\"action\": \"exit\", \"message\": \"This request was rate limited by Aikido Security!\", \"response_code\": 429}"
		}
	} else {
		// Otherwise, do an async call via gRPC
		go grpc.OnRequest(method, route, 1*time.Second)
	}

	return "{\"status\": \"ok\"}"
}
