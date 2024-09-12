package main

import (
	"fmt"
	"main/context"
	"main/grpc"
	"main/log"
	"main/utils"
	"time"
)

func OnRequestInit() string {
	method := context.GetMethod()
	route := context.GetRoute()
	if method == "" || route == "" {
		return "{}"
	}

	log.Infof("[RINIT] Got request metadata: %s %s", method, route)

	if !grpc.AreEndpointsConfigured() {
		return "{}"
	}

	route = context.GetParsedRoute()

	endpointData, err := grpc.GetEndpointConfig(method, route)
	if err != nil {
		// This endpoint (method + route) has not configuration -> continue
		return "{}"
	}

	ip := context.GetIp()
	log.Infof("[RINIT] Got IP from request: %s", ip)

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

func OnRequestShutdownReporting(method string, route string, status_code int) {
	if method == "" || route == "" || status_code == 0 {
		return
	}

	log.Info("[RSHUTDOWN] Got request metadata: ", method, " ", route, " ", status_code)

	route = context.GetParsedRoute()

	if !utils.ShouldDiscoverRoute(status_code, route, method) {
		return
	}

	go grpc.OnRequestShutdown(method, route, status_code, 10*time.Millisecond)
}

func OnRequestShutdown() string {
	method := context.GetMethod()
	route := context.GetRoute()
	status_code := context.GetStatusCode()
	go OnRequestShutdownReporting(method, route, status_code)
	return ""
}
