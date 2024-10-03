package ssrf

import (
	"main/context"
	"main/helpers"
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
					return &interceptorResult
				}

				resolvedPrivateIp := helpers.TryResolveToPrivateIp(hostname)
				if resolvedPrivateIp != "" {
					interceptorResult.Metadata["resolvedIp"] = resolvedPrivateIp
					return &interceptorResult
				}

				context.ContextSetPartialInterceptorResult(interceptorResult)
			}
		}
	}
	return nil
}

func CheckEffectiveHostnameForSSRF(effectiveHostname string) *utils.InterceptorResult {
	interceptorResult := context.GetPartialInterceptorResult()
	if interceptorResult == nil {
		// The initially requested hostname was not found in the user input -> no SSRF
		return nil
	}

	resolvedPrivateIp := helpers.TryResolveToPrivateIp(effectiveHostname)
	if resolvedPrivateIp == "" {
		// Effective hostname did not resolve to private IP -> no SSRF
		return nil
	}

	// After redirects the effective hostname resolved to a private IP -> SSRF
	interceptorResult.Metadata["effectiveHostname"] = effectiveHostname
	interceptorResult.Metadata["resolvedIp"] = resolvedPrivateIp
	return interceptorResult
}

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
	return interceptorResult
}
