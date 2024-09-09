package main

import (
	. "main/aikido_types"
	"main/utils"
)

var PreFunctionExecutedHandlers = map[string]HandlerFunction{
	"curl_exec":      OnPreFunctionExecutedCurl,
	"shell_executed": OnPreShellExecuted, // exec, shell_exec, system, passthru, popen, proc_open
	"path_accessed":  OnPrePathAccessed,  // basename, chgrp, chmod, chown, clearstatcache, copy, dirname ...
}

func OnBeforeFunctionExecuted(data map[string]interface{}) string {
	functionName := utils.MustGetFromMap[string](data, "function_name")
	parameters := utils.MustGetFromMap[map[string]interface{}](data, "parameters")

	utils.KeyMustExist(PreFunctionExecutedHandlers, functionName)

	return PreFunctionExecutedHandlers[functionName](parameters)
}
