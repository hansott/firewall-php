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
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	resolvedIps, err := net.DefaultResolver.LookupHost(ctx, hostname)
	if err != nil {
		// If timeout is reached or the OS lookup fail, return an emtpy list of resolved IPs
		return []string{}
	}
	return resolvedIps
}

func FindPrivateIp(resolvedIps []string) string {
	for _, resolvedIp := range resolvedIps {
		if isPrivateIP(resolvedIp) {
			return resolvedIp
		}
	}
	return ""
}
