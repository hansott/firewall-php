package main

import (
	"main/grpc"
	"main/log"
	"main/utils"
)

func OnFunctionExecutedCurl(parameters map[string]interface{}) string {
	url := utils.GetFromMap[string](parameters, "url")
	if url == nil {
		return "{}"
	}
	domain := utils.GetDomain(*url)
	log.Info("Got domain: ", domain)
	go grpc.OnReceiveDomain(domain)
	return "{}"
}
