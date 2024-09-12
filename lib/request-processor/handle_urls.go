package main

import (
	"main/context"
	"main/grpc"
	"main/log"
	"main/utils"
)

func OnPreFunctionExecutedCurl() string {
	url := context.GetOutgoingRequestUrl()
	if url == "" {
		return "{}"
	}
	domain := utils.GetDomain(url)
	log.Info("[BEFORE] Got domain: ", domain)
	//TODO: check if domain is blacklisted
	return ""
}

func OnAfterFunctionExecutedCurl() string {
	url := context.GetOutgoingRequestUrl()
	port := context.GetOutgoingRequestPort()
	if url == "" || port == 0 {
		return ""
	}
	domain := utils.GetDomain(url)
	log.Info("[AFTER] Got domain: ", domain, " port: ", port)

	go grpc.OnDomain(domain, port)

	return ""
}
