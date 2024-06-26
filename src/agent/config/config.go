package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/globals"
	"main/log"
	"os"
	"time"
)

var quit chan struct{}

func loadEnvConfigIfExistsStr(environmentVariableName string, originalValue string) string {
	environmentVariableValue := os.Getenv(environmentVariableName)

	if environmentVariableValue == "" || environmentVariableValue == "UNSET" {
		return originalValue
	}
	return environmentVariableValue
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

	if err := json.Unmarshal(byteValue, &globals.LocalConfig); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal JSON: %v", err))
	}

	if err := log.SetLogLevel(globals.LocalConfig.LogLevel); err != nil {
		panic(fmt.Sprintf("Error setting log level: %s", err))
	}

	log.Infof("Loaded local config: %+v", globals.LocalConfig)
	if globals.Token == "" {
		log.Infof("Waiting for token to initiate cloud communication...")
	}
}

func Init() {
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
