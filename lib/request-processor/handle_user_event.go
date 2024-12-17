package main

import (
	"main/context"
	"main/grpc"
	"main/log"
)

func OnUserEvent() string {
	go grpc.OnMiddlewareInstalled()

	id := context.GetUserId()
	username := context.GetUserName()
	ip := context.GetIp()

	log.Infof("[UEVENT] Got user event: %s %s %s", id, username, ip)

	if id == "" || username == "" || ip == "" {
		return ""
	}

	go grpc.OnUserEvent(id, username, ip)
	return ""
}
