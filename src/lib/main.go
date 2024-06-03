package main

import "C"
import (
	"encoding/json"
	"fmt"
	"main/cloud"
	"main/log"
	"time"
)

type eventFunctionExecutedFn func(map[string]interface{}) string

var eventHandlers = map[string]eventFunctionExecutedFn{
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

	eventName := MustGetFromMap[string](event, "event")
	data := MustGetFromMap[map[string]interface{}](event, "data")

	CheckIfKeyExists(eventHandlers, eventName)

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

	var initData map[string]interface{}
	err := json.Unmarshal([]byte(initJson), &initData)
	if err != nil {
		panic(fmt.Sprintf("Error parsing JSON: %s", err))
	}

	log_level := MustGetFromMap[string](initData, "log_level")

	if err := log.SetLogLevel(log_level); err != nil {
		panic(fmt.Sprintf("Error setting log level: %s", err))
	}

	cloud.StartConfigThread("http://example.com")

	log.Debug("Init: ", initJson)

	time.Sleep(100 * time.Second)
	return true
}

//export Uninit
func Uninit() {
	log.Debug("Uninit: {}")
	cloud.StopConfigThread()
}

func main() {}
