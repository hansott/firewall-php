package ssrf

import "net"

func TryResolveHostnameToPrivateIp(hostname string) string {
	resolvedIps, err := net.LookupHost(hostname)
	if err != nil {
		return ""
	}
	for _, resolvedIp := range resolvedIps {
		if isPrivateIP(resolvedIp) {
			return resolvedIp
		}
	}
	return ""
}
