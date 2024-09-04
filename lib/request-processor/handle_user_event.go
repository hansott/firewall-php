package main

import (
	"main/context"
	"main/grpc"
	"main/log"
	"main/utils"
)

func OnUserEvent(data map[string]interface{}) string {
	id := utils.MustGetFromMap[string](data, "id")
	username := utils.MustGetFromMap[string](data, "username")

	ip := context.GetIp()

	log.Infof("[UEVENT] Got user event: %s %s %s", id, username, ip)

	if id == "" || username == "" || ip == "" {
		return "{}"
	}

	go grpc.OnUserEvent(id, username, ip)

	if utils.IsUserBlocked(id) {
		return `{"action": "exit", "message": "You are blocked by Aikido Firewall!", "response_code": 403}`
	}

	context.ContextSetUserId(id)
	return "{}"
}
