package ssrf

import (
	"main/context"
	"main/utils"
)

func CheckContextForSSRF(hostname string, port int, operation string) *utils.InterceptorResult {
	for _, source := range context.SOURCES {
		mapss := source.CacheGet()

		for str, path := range mapss {
			if findHostnameInUserInput(str, hostname, port) {
				interceptorResult := utils.InterceptorResult{
					Operation:     operation,
					Kind:          utils.Ssrf,
					Source:        source.Name,
					PathToPayload: path,
					Metadata:      getMetadataForSSRFAttack(hostname, port),
					Payload:       str,
				}

				if containsPrivateIPAddress(hostname) {
					return &interceptorResult
				} else {
					context.ContextSetPartialInterceptorResult(interceptorResult)
				}
			}
		}
	}
	return nil
}

func CheckResolvedIpForSSRF(resolvedIp string) *utils.InterceptorResult {
	interceptorResult := context.GetPartialInterceptorResult()
	if interceptorResult != nil && isPrivateIP(resolvedIp) {
		interceptorResult.Metadata["resolvedIp"] = resolvedIp
		return interceptorResult
	}
	return nil
}
