package main

import (
	"main/attack"
	"main/context"
	"main/grpc"
	"main/log"
	ssrf "main/vulnerabilities/ssrf"
)

/*
	Defends agains:

- basic SSRF (local IP address used as hostname)
- direct SSRF attacks (hostname that resolves directly to a local IP address - does not go through redirects)
- direct IMDS SSRF attacks (hostname is an IMDS IP)

All these checks first verify if the hostname was provided via user input.
Protects both curl and fopen wrapper functions (file_get_contents, etc...).
*/
func OnPreOutgoingRequest() string {
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

/*
	Downgrades a potential SSRF attack to a blind SSRF attack.
	Defends agains:

- re-direct SSRF attacks (redirects lead to a hostname that resolves to a local IP address)
- re-direct IMDS SSRF attacks (redirects lead to a hostname that resolves to an IMDS IP address)

All these checks first verify if the hostname was provided via user input.
Protects curl.
*/
func OnPostOutgoingRequest() string {
	hostname, port := context.GetOutgoingRequestHostnameAndPort()
	effectiveHostname, effectivePort := context.GetOutgoingRequestEffectiveHostnameAndPort()
	resolvedIp := context.GetOutgoingRequestResolvedIp()
	if hostname == "" {
		return ""
	}

	log.Info("[AFTER] Got domain: ", hostname, " port: ", port)

	go grpc.OnDomain(hostname, port)

	res := ssrf.CheckResolvedIpForSSRF(resolvedIp)
	if effectiveHostname != hostname {
		// After the request was made, if effective hostname is different that the initially requested one (redirects, ...)
		// -> check for SSRF
		// -> report it to the agent
		if res == nil {
			res = ssrf.CheckEffectiveHostnameForSSRF(effectiveHostname)
		}
		go grpc.OnDomain(effectiveHostname, effectivePort)
	}

	if res != nil {
		go grpc.OnAttackDetected(*res)
		return attack.GetAttackDetectedAction(*res)
	}
	return ""
}
