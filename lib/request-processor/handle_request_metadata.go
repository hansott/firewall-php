package main

import (
	"fmt"
	"main/grpc"
	"main/log"
	"main/utils"
	"time"
)

func OnRequestInit(data map[string]interface{}) string {
	method := utils.MustGetFromMap[string](data, "method")
	route := utils.MustGetFromMap[string](data, "route")
	ip := utils.MustGetFromMap[string](data, "remoteAddress")

	if method == "" || route == "" {
		return "{\"status\": \"ok\"}"
	}

	log.Infof("[RINIT] Got request metadata: %s %s (%s)", method, route, ip)

	route = utils.BuildRouteFromURL(route)

	endpointData, err := grpc.GetEndpointConfig(method, route)
	if err != nil {
		// This endpoint (method + route) has not configuration -> continue
		return "{\"status\": \"ok\"}"
	}

	if !utils.IsIpAllowed(endpointData.AllowedIPAddresses, ip) {
		message := "Your IP address is not allowed to access this resource!"
		if ip != "" {
			message += fmt.Sprintf(" (Your IP: %s)", ip)
		}
		return fmt.Sprintf(`{"action": "exit", "message": %s, "response_code": 403}`, message)
	}

	if endpointData.RateLimiting.Enabled {
		if !utils.IsIpExcludedFromRateLimiting(ip) {
			// If request is monitored for rate limiting and the IP is not excluded from rate limiting,
			// do a sync call via gRPC to see if the request should be aborded or not
			if !grpc.OnRequestInit(method, route, 10*time.Millisecond) {
				return "{\"action\": \"exit\", \"message\": \"This request was rate limited by Aikido Security!\", \"response_code\": 429}"
			}
		} else {
			log.Infof("IP \"%s\" is excluded from rate limiting!", ip)
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
