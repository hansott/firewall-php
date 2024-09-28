package helpers

func TryGetResolvedPrivateIp(resolvedIps []string) string {
	for _, resolvedIp := range resolvedIps {
		if isPrivateIP(resolvedIp) {
			return resolvedIp
		}
	}
	return ""
}
