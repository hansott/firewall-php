package main

//#include "../ContextCallback.c"
import "C"
import (
	"encoding/json"
	"fmt"
	. "main/aikido_types"
	"main/config"
	"main/context"
	"main/globals"
	"main/grpc"
	"main/log"
	"main/utils"
	"unsafe"
)

var eventHandlers = map[string]HandlerFunction{
	"before_function_executed": OnBeforeFunctionExecuted,
	"after_function_executed":  OnAfterFunctionExecuted,
	"method_executed":          OnMethodExecuted,
	"request_init":             OnRequestInit,
	"request_shutdown":         OnRequestShutdown,
	"user_event":               OnUserEvent,
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

	log.Debugf("Aikido Request Processor v%s started in \"%s\" mode!", globals.Version, globals.EnvironmentConfig.SAPI)
	log.Debugf("Init data: %s", initJson)

	if globals.EnvironmentConfig.SAPI != "cli" {
		grpc.Init()
	}
	return true
}

var CContextCallback C.ContextCallback

func GoContextCallback(contextId int) string {
	if CContextCallback == nil {
		return ""
	}

	contextData := C.call(CContextCallback, C.int(contextId))
	if contextData == nil {
		return ""
	}

	goContextData := C.GoString(contextData)

	/*
		In order to pass dynamic strings from the PHP extension (C++), we need a dynamically allocated buffer, that is allocated by the C++ extension.
		This buffer needs to be freed by the RequestProcessor (Go) once it has finished copying the data.
	*/
	C.free(unsafe.Pointer(contextData))
	return goContextData
}

//export RequestProcessorContextInit
func RequestProcessorContextInit(contextCallback C.ContextCallback) (initOk bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn("Recovered from panic:", r)
			initOk = false
		}
	}()

	CContextCallback = contextCallback
	return context.Init(GoContextCallback)
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

	dataMap := map[string]interface{}{}
	data := utils.GetFromMap[map[string]interface{}](event, "data")
	if data != nil {
		dataMap = *data
	}

	utils.KeyMustExist(eventHandlers, eventName)

	goString := eventHandlers[eventName](dataMap)
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
	if globals.EnvironmentConfig.SAPI != "cli" {
		grpc.Uninit()
	}
	log.Debugf("Aikido Request Processor v%s stopped!", globals.Version)
}

func main() {}
