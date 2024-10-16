package main

import (
	"fmt"
	"main/api_discovery"
	"main/context"
	"main/grpc"
	"main/ipc/protos"
	"main/log"
	"main/utils"
	"time"
)

func OnRequestInit() string {
	context.Clear()

	method := context.GetMethod()
	route := context.GetParsedRoute()
	if method == "" || route == "" {
		return "{}"
	}

	log.Infof("[RINIT] Got request metadata: %s %s", method, route)

	if !utils.AreEndpointsConfigured() {
		log.Debugf("[RINIT] No endpoints configured! Skipping checks...")
		return "{}"
	}

	endpointData, err := utils.GetEndpointConfig(method, route)
	if err != nil {
		log.Debugf("[RINIT] Method+route in not configured in endpoints! Skipping checks...")
		return "{}"
	}

	ip := context.GetIp()

	if !utils.IsIpAllowed(endpointData.AllowedIPAddresses, ip) {
		message := "Your IP address is not allowed to access this resource!"
		if ip != "" {
			message += fmt.Sprintf(" (Your IP: %s)", ip)
		}
		return fmt.Sprintf(`{"action": "exit", "message": "%s", "response_code": 403}`, message)
	}

	if endpointData.RateLimiting.Enabled {
		if !context.IsIpBypassed() {
			// If request is monitored for rate limiting and the IP is not bypassed,
			// do a sync call via gRPC to see if the request should be aborded or not
			if !grpc.OnRequestInit(method, route, 10*time.Millisecond) {
				return `{"action": "exit", "message": "This request was rate limited by Aikido Security!", "response_code": 429}`
			}
		} else {
			log.Infof("IP \"%s\" is bypassed for rate limiting!", ip)
		}
	}

	return "{}"
}

func OnRequestShutdownReporting(method string, route string, statusCode int, apiSpec *protos.APISpec) {
	if method == "" || route == "" || statusCode == 0 {
		return
	}

	log.Info("[RSHUTDOWN] Got request metadata: ", method, " ", route, " ", statusCode)

	if !utils.ShouldDiscoverRoute(statusCode, route, method) {
		return
	}

	log.Info("[RSHUTDOWN] Got API spec: ", apiSpec)
	grpc.OnRequestShutdown(method, route, statusCode, 5*time.Second, apiSpec)
}

func OnRequestShutdown() string {
	go OnRequestShutdownReporting(context.GetMethod(), context.GetParsedRoute(), context.GetStatusCode(), api_discovery.GetApiInfo())
	context.Clear()
	return ""
}
