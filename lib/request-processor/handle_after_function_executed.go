package main

import (
	. "main/aikido_types"
	"main/utils"
)

var afterFunctionExecutedHandlers = map[string]HandlerFunction{
	"curl_exec": OnAfterFunctionExecutedCurl,
}

func OnAfterFunctionExecuted(data map[string]interface{}) string {
	functionName := utils.MustGetFromMap[string](data, "function_name")
	parameters := utils.MustGetFromMap[map[string]interface{}](data, "parameters")

	utils.KeyMustExist(afterFunctionExecutedHandlers, functionName)

	return afterFunctionExecutedHandlers[functionName](parameters)
}
