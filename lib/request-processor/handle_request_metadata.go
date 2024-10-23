package main

import (
	"fmt"
	"main/api_discovery"
	"main/context"
	"main/grpc"
	"main/ipc/protos"
	"main/log"
	"main/utils"
)

func OnPreRequest() string {
	context.Clear()

	method := context.GetMethod()
	route := context.GetParsedRoute()
	if method == "" || route == "" {
		return ""
	}

	log.Infof("Got request metadata: %s %s", method, route)

	endpointData, err := utils.GetEndpointConfig(method, route)
	if err != nil {
		log.Debugf("Method+route in not configured in endpoints! Skipping checks...")
		return ""
	}

	ip := context.GetIp()

	if !utils.IsIpAllowed(endpointData.AllowedIPAddresses, ip) {
		message := "Your IP address is not allowed to access this resource!"
		if ip != "" {
			message += fmt.Sprintf(" (Your IP: %s)", ip)
		}
		return fmt.Sprintf(`{"action": "exit", "message": "%s", "response_code": 403}`, message)
	}

	return ""
}

func OnRequestShutdownReporting(method string, route string, statusCode int, user string, ip string, apiSpec *protos.APISpec) {
	if method == "" || route == "" || statusCode == 0 {
		return
	}

	log.Info("[RSHUTDOWN] Got request metadata: ", method, " ", route, " ", statusCode)

	if !utils.ShouldDiscoverRoute(statusCode, route, method) {
		return
	}

	log.Info("[RSHUTDOWN] Got API spec: ", apiSpec)
	grpc.OnRequestShutdown(method, route, statusCode, user, ip, apiSpec)
}

func OnPostRequest() string {
	go OnRequestShutdownReporting(context.GetMethod(), context.GetParsedRoute(), context.GetStatusCode(), context.GetUserId(), context.GetIp(), api_discovery.GetApiInfo())
	context.Clear()
	return ""
}
