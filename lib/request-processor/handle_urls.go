package main

import (
	"main/attack"
	"main/context"
	"main/grpc"
	"main/log"
	ssrf "main/vulnerabilities/ssrf"
)

func OnPreFunctionExecutedCurl() string {
	hostname := context.GetOutgoingRequestHostname()
	port := context.GetOutgoingRequestPort()
	operation := context.GetFunctionName()

	res := ssrf.CheckContextForSSRF(hostname, port, operation)
	if res != nil {
		go grpc.OnAttackDetected(*res)
		return attack.GetAttackDetectedAction(*res)
	}

	log.Info("[BEFORE] Got domain: ", hostname)
	//TODO: check if domain is blacklisted
	return ""
}

func OnAfterFunctionExecutedCurl() string {
	hostname := context.GetOutgoingRequestHostname()
	port := context.GetOutgoingRequestPort()
	resolvedIp := context.GetOutgoingRequestResolvedIp()
	if hostname == "" {
		return ""
	}

	log.Info("[AFTER] Got domain: ", hostname, " port: ", port)

	go grpc.OnDomain(hostname, port)

	res := ssrf.CheckResolvedIpForSSRF(resolvedIp)
	if res != nil {
		go grpc.OnAttackDetected(*res)
		return attack.GetAttackDetectedAction(*res)
	}
	return ""
}
