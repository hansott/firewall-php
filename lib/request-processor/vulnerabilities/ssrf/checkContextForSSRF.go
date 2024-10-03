package ssrf

import (
	"main/context"
	"main/utils"
)

/* This is called before a request is made to check for SSRF */
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
					interceptorResult.Metadata["isPrivateIp"] = "true"
					return &interceptorResult
				}

				resolvedIpStatus := getResolvedIpStatusForHostname(hostname)
				if resolvedIpStatus != nil {
					interceptorResult.Metadata["resolvedIp"] = resolvedIpStatus.ip
					if resolvedIpStatus.isIMDS {
						interceptorResult.Metadata["isIMDSIp"] = "true"
					}
					if resolvedIpStatus.isPrivate {
						interceptorResult.Metadata["isPrivateIp"] = "true"
					}
					return &interceptorResult
				}

				context.ContextSetPartialInterceptorResult(interceptorResult)
			}
		}
	}
	return nil
}

/* This is called after the request is made to check for SSRF in the effectiveHostname - hostname optained after redirects from the PHP library that made the request (curl) */
func CheckEffectiveHostnameForSSRF(effectiveHostname string) *utils.InterceptorResult {
	interceptorResult := context.GetPartialInterceptorResult()
	if interceptorResult == nil {
		// The initially requested hostname was not found in the user input -> no SSRF
		return nil
	}

	interceptorResult.Metadata["effectiveHostname"] = effectiveHostname
	resolvedIpStatus := getResolvedIpStatusForHostname(effectiveHostname)
	if resolvedIpStatus != nil {
		interceptorResult.Metadata["resolvedIp"] = resolvedIpStatus.ip
		if resolvedIpStatus.isIMDS {
			interceptorResult.Metadata["isIMDSIp"] = "true"
		}
		if resolvedIpStatus.isPrivate {
			interceptorResult.Metadata["isPrivateIp"] = "true"
		}
	}

	return nil
}

/* This is called after the request is made to check for SSRF in the resolvedIP - IP optained from the PHP library that made the request (curl) */
func CheckResolvedIpForSSRF(resolvedIp string) *utils.InterceptorResult {
	interceptorResult := context.GetPartialInterceptorResult()
	if interceptorResult == nil {
		// The initially requested hostname was not found in the user input -> no SSRF
		return nil
	}

	if !isPrivateIP(resolvedIp) {
		// The resolved IP address is not private -> no SSRF
		return nil
	}

	interceptorResult.Metadata["resolvedIp"] = resolvedIp
	interceptorResult.Metadata["isPrivateIp"] = "true"
	return interceptorResult
}
