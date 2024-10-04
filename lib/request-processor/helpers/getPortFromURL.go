package helpers

import (
	"net/url"
	"strconv"
)

func GetPortFromURL(u *url.URL) int {
	// If a port is explicitly specified and it's a valid integer
	if u.Port() != "" {
		port, err := strconv.Atoi(u.Port())
		if err == nil {
			return port
		}
	}

	// Default ports based on protocol
	switch u.Scheme {
	case "https":
		return 443
	case "http":
		return 80
	}
	return 0
}
