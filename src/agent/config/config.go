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
