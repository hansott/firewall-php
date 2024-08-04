package main

import (
	"main/grpc"
	"main/log"
	"main/utils"
)

func OnRequestMetadata(data map[string]interface{}) string {
	method := utils.MustGetFromMap[string](data, "method")
	route := utils.MustGetFromMap[string](data, "route")

	if method == "" || route == "" {
		return "{\"status\": \"ok\"}"
	}

	log.Info("Got request metadata: ", method, " ", route)

	route = utils.BuildRouteFromURL(route)
	go grpc.OnRequest(method, route)

	return "{\"status\": \"ok\"}"
}
