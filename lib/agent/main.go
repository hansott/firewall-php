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

	if !config.Init(initJson) {
		return false
	}
	log.Init()
	machine.Init()
	if !grpc.Init() {
		return false
	}

	cloud.Init()
	rate_limiting.Init()

	log.Infof("Aikido Agent v%s loaded!", globals.Version)
	return true
}

//export AgentUninit
func AgentUninit() {
	rate_limiting.Uninit()
	cloud.Uninit()
	grpc.Uninit()
	config.Uninit()

	log.Infof("Aikido Agent v%s unloaded!", globals.Version)
	log.Uninit()
}

func main() {}
