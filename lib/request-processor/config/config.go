package config

import (
	"encoding/json"
	"fmt"
	"main/globals"
	"main/log"
)

func ReloadAikidoConfig(initJson string) {
	err := json.Unmarshal([]byte(initJson), &globals.AikidoConfig)
	if err != nil {
		panic(fmt.Sprintf("Error parsing JSON to AikidoConfig: %s", err))
	}

	if err := log.SetLogLevel(globals.AikidoConfig.LogLevel); err != nil {
		panic(fmt.Sprintf("Error setting log level: %s", err))
	}
}

func Init(initJson string) {
	globals.CloudConfig.Block = -1

	err := json.Unmarshal([]byte(initJson), &globals.EnvironmentConfig)
	if err != nil {
		panic(fmt.Sprintf("Error parsing JSON to EnvironmentConfig: %s", err))
	}

	if globals.EnvironmentConfig.SocketPath == "" {
		panic("Socket path not set!")
	}

	ReloadAikidoConfig(initJson)
	log.Init()
}

func Uninit() {
	log.Uninit()
}
