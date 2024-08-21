package grpc

import (
	"main/globals"
	"main/ipc/protos"
)

func storeDomain(req *protos.Domain) {
	if req.GetPort() == 0 {
		return
	}

	globals.HostnamesMutex.Lock()
	defer globals.HostnamesMutex.Unlock()

	if _, ok := globals.Hostnames[req.GetDomain()]; !ok {
		globals.Hostnames[req.GetDomain()] = make(map[int]bool)
	}

	globals.Hostnames[req.GetDomain()][int(req.GetPort())] = true
}
