package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/cloud"
	"main/globals"
	"main/log"
	"os"
	"time"
)

var quit chan struct{}

func setConfigFromJson(jsonString []byte) {
	previousToken := globals.LocalConfig.Token

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
		log.Infof("Waiting for token to initiate cloud communication...")
	} else if globals.LocalConfig.Token != previousToken {
		cloud.Init()
	}

}

// Reloads the local config from /opt/aikido once every minute, in order to provide fast
// reload of critical info used for cloud comms or the blocking flags.
// This allows for fast local fixes if something goes wrong and needs to be enabled/disabled.
func loadLocalConfig() {
	globals.ConfigMutex.Lock()
	defer globals.ConfigMutex.Unlock()

	file, err := os.Open(globals.ConfigFilePath)
	if err != nil {
		file, err = os.Open(globals.DevConfigFilePath)
		if err != nil {
			panic(fmt.Sprintf("Failed to open config file: %v", err))
		}
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("Failed to read config file: %v", err))
	}
	setConfigFromJson(byteValue)
}

func Init(initJson string) {
	setConfigFromJson([]byte(initJson))

	loadLocalConfig()
	ticker := time.NewTicker(1 * time.Minute)
	quit = make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				loadLocalConfig()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func Uninit() {
	close(quit)
}
