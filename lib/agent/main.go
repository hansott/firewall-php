package main

import (
	"C"
	"main/config"
	"main/globals"
	"main/grpc"
	"main/log"
	"main/machine"
)
import (
	"main/cloud"
	"main/rate_limiting"
)

//export AgentInit
func AgentInit(initJson string) (initOk bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn("Recovered from panic:", r)
			initOk = false
		}
	}()

	log.Init()
	machine.Init()
	if !config.Init(initJson) || !grpc.Init() {
		return false
	}

	cloud.Init()
	rate_limiting.Init()

	log.Infof("Aikido Agent v%s loaded!", globals.Version)
	return true
}

//export AgentUninit
func AgentUninit() {
	grpc.Uninit()

	log.Infof("Aikido Agent v%s unloaded!", globals.Version)
	log.Uninit()
}

func main() {}
