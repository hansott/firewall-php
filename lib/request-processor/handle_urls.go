package main

import (
	"main/attack"
	"main/context"
	"main/grpc"
	"main/log"
	"main/utils"
	ssrf "main/vulnerabilities/ssrf"
)

func OnPreFunctionExecutedCurl() string {
	url := context.GetOutgoingRequestUrl()
	operation := context.GetFunctionName()

	if url == "" {
		return "{}"
	}

	res := ssrf.CheckContextForSSRF(url, operation)
	if res != nil {
		go grpc.OnAttackDetected(*res)
		return attack.GetAttackDetectedAction(*res)
	}

	domain := utils.GetDomain(url)
	log.Info("[BEFORE] Got domain: ", domain)
	//TODO: check if domain is blacklisted
	return ""
}

func OnAfterFunctionExecutedCurl() string {
	url := context.GetOutgoingRequestUrl()
	port := context.GetOutgoingRequestPort()
	if url == "" {
		return ""
	}
	domain := utils.GetDomain(url)
	log.Info("[AFTER] Got domain: ", domain, " port: ", port)

	go grpc.OnDomain(domain, port)

	return ""
}
