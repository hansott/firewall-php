package main

import (
	"encoding/json"
	"main/context"
	"main/grpc"
	"main/log"
	"main/utils"
	"time"
)

func GetStoreAction(actionType, trigger, ip string) string {
	actionMap := map[string]interface{}{
		"action":  "store",
		"type":    actionType,
		"trigger": trigger,
	}
	if trigger == "ip" {
		actionMap["ip"] = ip
	}
	actionJson, err := json.Marshal(actionMap)
	if err != nil {
		return ""
	}
	return string(actionJson)
}

func OnGetBlockingStatus() string {
	go grpc.OnMiddlewareInstalled()

	userId := context.GetUserId()
	if utils.IsUserBlocked(userId) {
		log.Infof("User \"%s\" is blocked!", userId)
		return GetStoreAction("blocked", "user", "")
	}

	ip := context.GetIp()
	if ipBlockingStatus := utils.GetIpBlockingStatus(ip); ipBlockingStatus != "" {
		log.Infof("IP \"%s\" is blocked by %s!", ip, ipBlockingStatus)
		return GetStoreAction("blocked", ipBlockingStatus, ip)
	}

	method := context.GetMethod()
	route := context.GetParsedRoute()
	if method == "" || route == "" {
		return ""
	}

	endpointData, err := utils.GetEndpointConfig(method, route)
	if err != nil {
		log.Debugf("Method+route is not configured in endpoints! Skipping checks...")
		return ""
	}

	if endpointData.RateLimiting.Enabled {
		if !context.IsIpBypassed() {
			// If request is monitored for rate limiting and the IP is not bypassed,
			// do a sync call via gRPC to see if the request should be blocked or not
			rateLimitingStatus := grpc.GetRateLimitingStatus(method, route, userId, ip, 10*time.Millisecond)
			if rateLimitingStatus != nil && rateLimitingStatus.Block {
				log.Infof("Request made from IP \"%s\" is ratelimited by \"%s\"!", ip, rateLimitingStatus.Trigger)
				return GetStoreAction("ratelimited", rateLimitingStatus.Trigger, ip)
			}
		} else {
			log.Infof("IP \"%s\" is bypassed for rate limiting!", ip)
		}
	}

	if !utils.IsIpAllowed(endpointData.AllowedIPAddresses, ip) {
		log.Infof("IP \"%s\" is not allowd to access this endpoint!", ip)
		return GetStoreAction("blocked", "ip", ip)
	}
	return ""
}
