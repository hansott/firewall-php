package config

import (
	"encoding/json"
	"fmt"
	"main/globals"
	"main/log"
)

func setConfigFromJson(jsonString []byte) bool {
	if err := json.Unmarshal(jsonString, &globals.EnvironmentConfig); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal JSON: %v", err))
	}

	if globals.EnvironmentConfig.LogLevel != "" {
		if err := log.SetLogLevel(globals.EnvironmentConfig.LogLevel); err != nil {
			panic(fmt.Sprintf("Error setting log level: %s", err))
		}
	}

	log.Infof("Loaded local config: %+v", globals.EnvironmentConfig)

	if globals.EnvironmentConfig.SocketPath == "" {
		log.Errorf("No socket path set! Aikido agent will not load!")
		return false
	}

	if globals.EnvironmentConfig.Token == "" {
		log.Infof("No token set! Aikido agent will not load!")
		return false
	}

	return true
}

func Init(initJson string) bool {
	return setConfigFromJson([]byte(initJson))
}

func Uninit() {

}
