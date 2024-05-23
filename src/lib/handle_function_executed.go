package main

type functionExecutedHandlersFn func(map[string]interface{}) string

var functionExecutedHandlers = map[string]functionExecutedHandlersFn{
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
	functionName := MustGetFromMap[string](data, "function_name")
	parameters := MustGetFromMap[map[string]interface{}](data, "parameters")

	ExitIfKeyDoesNotExistInMap(functionExecutedHandlers, functionName)

	return functionExecutedHandlers[functionName](parameters)
}
