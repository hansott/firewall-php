package ssrf

import "main/helpers"

type ResolvedIpStatus struct {
	ip        string
	isPrivate bool
	isIMDS    bool
}

func getResolvedIpStatusForHostname(hostname string) *ResolvedIpStatus {
	resolvedIps := helpers.ResolveHostname(hostname)
	imdsIP := TryGetIMDSIp(hostname, resolvedIps)
	if imdsIP != "" {
		return &ResolvedIpStatus{ip: imdsIP, isPrivate: false, isIMDS: true}
	}

	privateIp := helpers.TryGetPrivateIp(resolvedIps)
	if privateIp != "" {
		return &ResolvedIpStatus{ip: privateIp, isPrivate: true, isIMDS: false}
	}
	return nil
}
