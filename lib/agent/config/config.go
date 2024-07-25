package config

import (
	"encoding/json"
	"fmt"
	"main/cloud"
	"main/globals"
	"main/log"
)

var quit chan struct{}

func setConfigFromJson(jsonString []byte) bool {
	if err := json.Unmarshal(jsonString, &globals.LocalConfig); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal JSON: %v", err))
	}

	if globals.LocalConfig.LogLevel != "" {
		if err := log.SetLogLevel(globals.LocalConfig.LogLevel); err != nil {
			panic(fmt.Sprintf("Error setting log level: %s", err))
		}
	}

	log.Infof("Loaded local config: %+v", globals.LocalConfig)

	if globals.LocalConfig.Token == "" {
		log.Infof("No token set! Aikido agent will not load!")
		return false
	}

	cloud.Init()
	return true
}

func Init(initJson string) bool {
	return setConfigFromJson([]byte(initJson))
}

func Uninit() {
	close(quit)
}
