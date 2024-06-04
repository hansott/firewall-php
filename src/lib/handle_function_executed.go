package main

import (
	. "main/aikido_types"
	"main/utils"
)

var functionExecutedHandlers = map[string]FunctionExecutedHandlersFn{
	"curl_init":   OnFunctionExecutedCurl,
	"curl_setopt": OnFunctionExecutedCurl,

	"exec":       OnFunctionExecutedShell,
	"shell_exec": OnFunctionExecutedShell,
	"system":     OnFunctionExecutedShell,
	"passthru":   OnFunctionExecutedShell,
	"popen":      OnFunctionExecutedShell,
	"proc_open":  OnFunctionExecutedShell,
}

func OnFunctionExecuted(data map[string]interface{}) string {
	functionName := utils.MustGetFromMap[string](data, "function_name")
	parameters := utils.MustGetFromMap[map[string]interface{}](data, "parameters")

	utils.CheckIfKeyExists(functionExecutedHandlers, functionName)

	return functionExecutedHandlers[functionName](parameters)
}
