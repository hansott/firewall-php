package helpers

import "net"

func TryResolveToPrivateIp(hostname string) string {
	resolvedIps, _ := net.LookupHost(hostname)
	for _, resolvedIp := range resolvedIps {
		if isPrivateIP(resolvedIp) {
			return resolvedIp
		}
	}
	return ""
}
