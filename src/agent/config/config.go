package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/globals"
	"main/log"
	"os"
	"strings"
	"time"
)

var quit chan struct{}

func loadEnvConfigIfExistsStr(environmentVariableName string, originalValue string) string {
	environmentVariableValue := os.Getenv(environmentVariableName)

	if environmentVariableValue == "" {
		return originalValue
	}
	return environmentVariableValue
}

func loadEnvConfigIfExistsBool(environmentVariableName string, originalValue bool) bool {
	return strings.ToLower(os.Getenv(environmentVariableName)) == "true"
}

// Reloads the local config from /opt/aikido once every minute, in order to provide fast
// reload of critical info like the token used for cloud comms or the blocking flags.
// Also loads the same from ENV variables (if it exists). The ENV config takes precedence over the json config.
// This allows for fast local fixes if something goes wrong and needs to be enabled/disabled.
func loadLocalConfig() {
	globals.ConfigMutex.Lock()
	defer globals.ConfigMutex.Unlock()

	file, err := os.Open(globals.ConfigFilePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to open config file: %v", err))
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

	globals.LocalConfig.LogLevel = loadEnvConfigIfExistsStr("AIKIDO_LOG_LEVEL", globals.LocalConfig.LogLevel)
	globals.LocalConfig.Endpoint = loadEnvConfigIfExistsStr("AIKIDO_ENDPOINT", globals.LocalConfig.Endpoint)
	globals.LocalConfig.Token = loadEnvConfigIfExistsStr("AIKIDO_TOKEN", globals.LocalConfig.Token)
	globals.LocalConfig.Blocking = loadEnvConfigIfExistsBool("AIKIDO_BLOCKING", globals.LocalConfig.Blocking)

	log.Infof("Loaded local config: %+v", globals.LocalConfig)
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
