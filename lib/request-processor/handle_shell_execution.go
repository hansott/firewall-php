package main

import (
	"main/attack"
	"main/grpc"
	"main/log"
	"main/utils"
	shell_injection "main/vulnerabilities/shell-injection"
)

func OnPreFunctionExecutedShell(parameters map[string]interface{}) string {
	cmd := utils.GetFromMap[string](parameters, "cmd")
	operation := utils.GetFromMap[string](parameters, "operation")
	if cmd == nil {
		return "{}"
	}

	log.Info("Got shell command: ", *cmd)
	res := shell_injection.CheckContextForShellInjection(*cmd, *operation)
	if res != nil {
		go grpc.OnAttackDetected(*res)
		return attack.GetAttackDetectedAction(*res)
	}
	return "{}"
}
