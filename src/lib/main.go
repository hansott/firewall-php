package main

import "C"
import (
	"encoding/json"
	"fmt"
	. "main/aikido_types"
	"main/globals"
	"main/log"
	"main/utils"
)

var eventHandlers = map[string]HandlerFunction{
	"function_executed": OnFunctionExecuted,
	"method_executed":   OnMethodExecuted,
}

//export OnEvent
func OnEvent(eventJson string) (outputJson string) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn("Recovered from panic:", r)
			outputJson = "{}"
		}
	}()

	log.Debug("OnEvent: ", eventJson)

	var event map[string]interface{}
	err := json.Unmarshal([]byte(eventJson), &event)
	if err != nil {
		panic(fmt.Sprintf("Error parsing JSON: %s", err))
	}

	eventName := utils.MustGetFromMap[string](event, "event")
	data := utils.MustGetFromMap[map[string]interface{}](event, "data")

	utils.CheckIfKeyExists(eventHandlers, eventName)

	return eventHandlers[eventName](data)
}

//export Init
func Init(initJson string) (initOk bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn("Recovered from panic:", r)
			initOk = false
		}
	}()

	err := json.Unmarshal([]byte(initJson), &globals.InitData)
	if err != nil {
		panic(fmt.Sprintf("Error parsing JSON: %s", err))
	}

	if err := log.SetLogLevel(globals.InitData.LogLevel); err != nil {
		panic(fmt.Sprintf("Error setting log level: %s", err))
	}

	log.Debug("Init: ", initJson)

	return true
}

//export Uninit
func Uninit() {
	log.Debug("Uninit: {}")
}

func main() {}
