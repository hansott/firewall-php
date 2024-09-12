package main

import (
	"main/context"
	"main/grpc"
	"main/log"
	"main/utils"
)

func OnUserEvent() string {
	id := context.GetUserId()
	username := context.GetUserName()

	ip := context.GetIp()

	log.Infof("[UEVENT] Got user event: %s %s %s", id, username, ip)

	if id == "" || username == "" || ip == "" {
		return ""
	}

	go grpc.OnUserEvent(id, username, ip)

	if utils.IsUserBlocked(id) {
		return `{"action": "exit", "message": "You are blocked by Aikido Firewall!", "response_code": 403}`
	}

	return ""
}
