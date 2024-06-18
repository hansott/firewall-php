package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/globals"
	"main/log"
	"os"
)

func Init() {
	file, err := os.Open(globals.ConfigFilePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to open config file: %v", err))
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("Failed to read config file: %v", err))
	}

	if err := json.Unmarshal(byteValue, &globals.Config); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal JSON: %v", err))
	}

	log.Infof("Loaded config: %+v", globals.Config)
}
