package main

import (
	"main/attack"
	"main/context"
	"main/grpc"
	"main/log"
	ssrf "main/vulnerabilities/ssrf"
)

func OnPreFunctionExecutedCurl() string {
	hostname, port := context.GetOutgoingRequestHostnameAndPort()
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

func OnPostFunctionExecutedCurl() string {
	hostname, port := context.GetOutgoingRequestHostnameAndPort()
	effectiveHostname, effectivePort := context.GetOutgoingRequestEffectiveHostnameAndPort()
	resolvedIp := context.GetOutgoingRequestResolvedIp()
	if hostname == "" {
		return ""
	}

	log.Info("[AFTER] Got domain: ", hostname, " port: ", port)

	go grpc.OnDomain(hostname, port)
	if effectiveHostname != hostname {
		// After the request was made, if effective hostname is different that the initially requested one (redirects, ...)
		// than report it also to the agent
		go grpc.OnDomain(effectiveHostname, effectivePort)
	}

	res := ssrf.CheckResolvedIpForSSRF(resolvedIp)
	if res != nil {
		go grpc.OnAttackDetected(*res)
		return attack.GetAttackDetectedAction(*res)
	}
	return ""
}
