package config

import (
	"encoding/json"
	"fmt"
	"main/globals"
	"main/log"
)

func Init(initJson string) {
	globals.CloudConfig.Block = -1

	err := json.Unmarshal([]byte(initJson), &globals.EnvironmentConfig)
	if err != nil {
		panic(fmt.Sprintf("Error parsing JSON: %s", err))
	}

	if globals.EnvironmentConfig.SocketPath == "" {
		panic("Socket path not set!")
	}

	log.Init()
}

func Uninit() {
	log.Uninit()
}
