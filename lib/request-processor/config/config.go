package config

import (
	"encoding/json"
	"fmt"
	"main/globals"
	"main/log"
)

func Init(initJson string) {
	globals.CloudConfig.Block = -1

	err := json.Unmarshal([]byte(initJson), &globals.InitData)
	if err != nil {
		panic(fmt.Sprintf("Error parsing JSON: %s", err))
	}

	if globals.InitData.SAPI != "cli" {
		log.Init()
	} else {
		if err := log.SetLogLevel(globals.InitData.LogLevel); err != nil {
			panic(fmt.Sprintf("Error setting log level: %s", err))
		}
	}
}
