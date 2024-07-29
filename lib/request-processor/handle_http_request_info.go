package main

import (
	"main/grpc"
	"main/log"
	"main/utils"
)

func OnHttpRequestInfo(data map[string]interface{}) string {
	method := utils.MustGetFromMap[string](data, "method")
	route := utils.MustGetFromMap[string](data, "route")

	if method == "" || route == "" {
		log.Error("Missing required fields in OnHttpRequestInfo")
		return "{\"status\": \"ok\"}"
	}

	log.Info("Got HTTP request: ", method, " ", route)
	go grpc.OnReceiveHttpRequestInfo(method, route)

	return "{\"status\": \"ok\"}"
}
