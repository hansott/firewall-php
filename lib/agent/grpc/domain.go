package grpc

import (
	"main/globals"
	"main/ipc/protos"
)

func storeDomain(req *protos.Domain) {
	globals.HostnamesMutex.Lock()
	defer globals.HostnamesMutex.Unlock()

	globals.Hostnames[req.GetDomain()] = true
}
