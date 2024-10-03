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
func TryResolveToPrivateIp(hostname string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	resolvedIps, err := net.DefaultResolver.LookupHost(ctx, hostname)
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
