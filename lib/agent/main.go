package main

import (
	"C"
	"main/config"
	"main/globals"
	"main/grpc"
	"main/log"
	"main/machine"
)

//export AgentInit
func AgentInit(initJson string) (initOk bool) {
	log.SetLogLevel("DEBUG")

	log.Init()
	log.Infof("Aikido Agent v%s loaded!", globals.Version)

	config.Init(initJson)
	machine.Init()
	go grpc.Init()
	return true
}

//export AgentUninit
func AgentUninit() {
	grpc.Uninit()

	log.Infof("Aikido Agent v%s unloaded!", globals.Version)
	log.Uninit()
}

func main() {}
