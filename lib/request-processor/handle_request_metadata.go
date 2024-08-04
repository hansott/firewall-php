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

	if grpc.IsRequestConfiguredForRateLimiting(method, route) {
		if !grpc.OnRequest(method, route, 5*time.Millisecond) {
			return `{"action": "throw", "message": "Request was rate limited by Aikido Security", "code": -1}`
		}
	} else {
		go grpc.OnRequest(method, route, 2*time.Second)
	}

	return "{\"status\": \"ok\"}"
}
