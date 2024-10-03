package helpers

import (
	"context"
	"net"
	"time"
)

/*
This function tries to resolve the hostname to a private IP adress, if possible.
It does this by calling DNS resolution from the OS (getaddrinfo for Linux).
*/
func ResolveHostname(hostname string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	resolvedIps, _ := net.DefaultResolver.LookupHost(ctx, hostname)
	return resolvedIps
}

func TryGetPrivateIp(resolvedIps []string) string {
	for _, resolvedIp := range resolvedIps {
		if isPrivateIP(resolvedIp) {
			return resolvedIp
		}
	}
	return ""
}
