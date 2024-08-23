package main

import "C"
import (
	"encoding/json"
	"fmt"
	. "main/aikido_types"
	"main/config"
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

	config.Init(initJson)

	log.Debugf("Aikido Request Processor v%s started in \"%s\" mode!", globals.Version, globals.InitData.SAPI)
	log.Debugf("Init data: %s", initJson)

	if globals.InitData.SAPI != "cli" {
		grpc.Init()
	}
	return true
}

//export RequestProcessorOnEvent
func RequestProcessorOnEvent(eventJson string) (outputJson *C.char) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn("Recovered from panic:", r)
			outputJson = C.CString("{}")
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

	utils.KeyMustExist(eventHandlers, eventName)

	goString := eventHandlers[eventName](data)
	cString := C.CString(goString)
	return cString
}

/*
	Returns -1 if the config was not yet pulled from Agent.
	Otherwise, if blocking was set from cloud, it returns that value.
	Otherwise, it returns the environment value.
*/
//export RequestProcessorGetBlockingMode
func RequestProcessorGetBlockingMode() int {
	return utils.GetBlockingMode()
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
