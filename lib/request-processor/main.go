package main

//#include "../API.h"
import "C"
import (
	. "main/aikido_types"
	"main/config"
	"main/context"
	"main/globals"
	"main/grpc"
	"main/log"
	"main/utils"
	"unsafe"
)

var eventHandlers = map[int]HandlerFunction{
	C.EVENT_PRE_REQUEST:            OnRequestInit,
	C.EVENT_POST_REQUEST:           OnRequestShutdown,
	C.EVENT_PRE_USER:               OnUserEvent,
	C.EVENT_PRE_OUTGOING_REQUEST:   OnPreFunctionExecutedCurl,
	C.EVENT_POST_OUTGOING_REQUEST:  OnAfterFunctionExecutedCurl,
	C.EVENT_PRE_SHELL_EXECUTED:     OnPreShellExecuted,
	C.EVENT_PRE_PATH_ACCESSED:      OnPrePathAccessed,
	C.EVENT_PRE_EXECUTED_PDO_QUERY: OnPreExecutedPdoQuery,
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
func RequestProcessorOnEvent(eventId int) (outputJson *C.char) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn("Recovered from panic:", r)
			outputJson = nil
		}
	}()

	log.Debug("OnEvent: ", eventId)

	goString := eventHandlers[eventId]()
	if goString == "" {
		return nil
	}
	return C.CString(goString)
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
