package grpc

import (
	"main/globals"
)

func storeDomain(domain string, port int) {
	if port == 0 {
		return
	}

	globals.HostnamesMutex.Lock()
	defer globals.HostnamesMutex.Unlock()

	if _, ok := globals.Hostnames[domain]; !ok {
		globals.Hostnames[domain] = make(map[int]bool)
	}

	globals.Hostnames[domain][int(port)] = true
}
