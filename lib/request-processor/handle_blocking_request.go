package main

import (
	"encoding/json"
	"main/context"
	"main/grpc"
	"main/log"
	"main/utils"
	"time"
)

func OnGetBlockingStatus() string {
	userId := context.GetUserId()
	if utils.IsUserBlocked(userId) {
		return `{"action": "store", "type": "blocked", "trigger": "user"}`
	}

	method := context.GetMethod()
	route := context.GetParsedRoute()
	if method == "" || route == "" {
		return ""
	}

	endpointData, err := utils.GetEndpointConfig(method, route)
	if err != nil {
		log.Debugf("[RINIT] Method+route in not configured in endpoints! Skipping checks...")
		return ""
	}

	ip := context.GetIp()

	if endpointData.RateLimiting.Enabled {
		if !context.IsIpBypassed() {
			// If request is monitored for rate limiting and the IP is not bypassed,
			// do a sync call via gRPC to see if the request should be blocked or not
			rateLimitingStatus := grpc.GetRateLimitingStatus(method, route, userId, ip, 10*time.Millisecond)
			if rateLimitingStatus != nil && rateLimitingStatus.Block {
				action := map[string]interface{}{
					"action":  "store",
					"type":    "ratelimited",
					"trigger": rateLimitingStatus.Trigger,
				}
				if rateLimitingStatus.Trigger == "ip" {
					action["ip"] = ip
				}
				actionJson, err := json.Marshal(action)
				if err == nil {
					return string(actionJson)
				}
			}
		} else {
			log.Infof("IP \"%s\" is bypassed for rate limiting!", ip)
		}
	}

	return ""
}
