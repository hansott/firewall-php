package main

import (
	"main/grpc"
	"main/log"
	"main/utils"
)

func OnUserEvent(data map[string]interface{}) string {
	id := utils.MustGetFromMap[string](data, "id")
	username := utils.MustGetFromMap[string](data, "username")
	ip := utils.MustGetFromMap[string](data, "ip")
	log.Info("[UEVENT] Got user event: ", id, " ", username, " ", ip)

	if id == "" || username == "" || ip == "" {
		return "{\"status\": \"ok\"}"
	}

	go grpc.OnUserEvent(id, username, ip)

	return "{\"status\": \"ok\"}"
}
