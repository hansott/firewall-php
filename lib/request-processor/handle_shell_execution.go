package main

import (
	"main/log"
	"main/utils"
	shell_injection "main/vulnerabilities/shell-injection"
)

func OnFunctionExecutedShell(parameters map[string]interface{}) string {
	cmd := utils.GetFromMap[string](parameters, "cmd")
	context := utils.ParseContext(utils.MustGetFromMap[map[string]interface{}](parameters, "context"))
	operation := utils.GetFromMap[string](parameters, "operation")
	if cmd == nil {
		return "{}"
	}

	log.Info("Got shell command: ", *cmd)
	res := shell_injection.CheckContextForShellInjection(*cmd, *operation, context)
	if res != nil {
		return `{"action": "throw", "message": "Shell injection detected", "code": -1}`
	}
	return "{}"
}
