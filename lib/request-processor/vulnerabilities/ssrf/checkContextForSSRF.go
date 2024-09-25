package ssrf

import (
	"main/context"
	"main/helpers"
	"main/utils"
)

func CheckContextForSSRF(url string, operation string) *utils.InterceptorResult {
	hostname, port := helpers.GetHostnameAndPortFromURL(url)

	for _, source := range context.SOURCES {
		mapss := source.CacheGet()

		for str, path := range mapss {
			found := findHostnameInUserInput(str, hostname, port)
			if found && containsPrivateIPAddress(hostname) {
				return &utils.InterceptorResult{
					Operation:     operation,
					Kind:          utils.Ssrf,
					Source:        source.Name,
					PathToPayload: path,
					Metadata:      getMetadataForSSRFAttack(hostname, port),
					Payload:       str,
				}
			}
		}
	}
	return nil
}
