package ssrf

import (
	"main/helpers"
	"strings"
)

/**
 * Check if the hostname contains a private IP address
 * This function is used to detect obvious SSRF attacks (with a private IP address being used as the hostname)
 *
 * Examples
 * http://192.168.0.1/some/path
 * http://[::1]/some/path
 * http://localhost/some/path
 *
 * This function gets to see "192.168.0.1", "[::1]", and "localhost"
 *
 * We won't flag this-domain-points-to-a-private-ip.com
 * This will be handled by the inspectDNSLookupCalls function
 */
func containsPrivateIPAddress(hostname string) bool {
	if hostname == "localhost" {
		return true
	}

	// Attempt to parse the hostname as an IP address or domain
	url := helpers.TryParseURL("http://" + hostname)
	if url == nil {
		return false
	}

	// IPv6 addresses are enclosed in square brackets
	// e.g. http://[::1]
	host := url.Hostname()
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		ipv6 := host[1 : len(host)-1] // Extract IPv6 address inside brackets
		if isPrivateIP(ipv6) {
			return true
		}
	}

	// Check if the hostname is an IPv4 or an IPv6 private IP
	return isPrivateIP(host)
}
