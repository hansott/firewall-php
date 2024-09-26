package helpers

import (
	"strconv"
	"strings"
)

func ParsePort(portStr string) int {
	if portStr == "" {
		return 0
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0
	}
	return port
}

func FixURL(url string) string {
	if !strings.HasPrefix(url, "https://") && strings.HasPrefix(url, "https:/") {
		return strings.Replace(url, "https:/", "https://", 1)
	}
	if !strings.HasPrefix(url, "http://") && strings.HasPrefix(url, "http:/") {
		return strings.Replace(url, "http:/", "http://", 1)
	}
	return url
}

func GetHostnameAndPortFromURL(url string) (string, int) {
	parsedURL := TryParseURL(FixURL(url))
	if parsedURL == nil {
		return "", 0
	}
	return parsedURL.Hostname(), GetPortFromURL(parsedURL)
}
