package main

import "C"
import (
	"encoding/json"
	"fmt"
	. "main/aikido_types"
	"main/globals"
	"main/grpc"
	"main/log"
	"main/utils"
)

var eventHandlers = map[string]HandlerFunction{
	"before_function_executed": OnBeforeFunctionExecuted,
	"after_function_executed":  OnAfterFunctionExecuted,
	"method_executed":          OnMethodExecuted,
	"request_init":             OnRequestInit,
	"request_shutdown":         OnRequestShutdown,
}

//export RequestProcessorInit
func RequestProcessorInit(initJson string) (initOk bool) {
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

	if globals.InitData.SAPI != "cli" {
		log.Init()
	} else {
		if err := log.SetLogLevel(globals.InitData.LogLevel); err != nil {
			panic(fmt.Sprintf("Error setting log level: %s", err))
		}
	}

	log.Debugf("Aikido Request Processor v%s started in \"%s\" mode!", globals.Version, globals.InitData.SAPI)
	log.Debugf("Init data: %s", initJson)

	if globals.InitData.SAPI != "cli" {
		grpc.Init()
	}
	return true
}

//export RequestProcessorOnEvent
func RequestProcessorOnEvent(eventJson string) (outputJson string) {
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

//export RequestProcessorUninit
func RequestProcessorUninit() {
	log.Debug("Uninit: {}")
	if globals.InitData.SAPI != "cli" {
		grpc.Uninit()
	}
	log.Debugf("Aikido Request Processor v%s stopped!", globals.Version)
}

func main() {}
