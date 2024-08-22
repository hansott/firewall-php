package main

import (
	"main/grpc"
	"main/log"
	"main/utils"
)

func OnBeforeFunctionExecutedCurl(parameters map[string]interface{}) string {
	url := utils.GetFromMap[string](parameters, "url")
	if url == nil {
		return "{}"
	}
	domain := utils.GetDomain(*url)
	log.Info("[BEFORE] Got domain: ", domain)
	//TODO: check if domain is blacklisted
	return "{}"
}

func OnAfterFunctionExecutedCurl(parameters map[string]interface{}) string {
	url := utils.GetFromMap[string](parameters, "url")
	port := utils.GetFromMap[float64](parameters, "port")
	if url == nil || port == nil {
		return "{}"
	}
	domain := utils.GetDomain(*url)
	log.Info("[AFTER] Got domain: ", domain, " port: ", int(*port))

	go grpc.OnDomain(domain, int(*port))

	return "{}"
}
