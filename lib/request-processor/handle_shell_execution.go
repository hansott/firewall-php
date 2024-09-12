package main

import (
	"main/attack"
	"main/context"
	"main/grpc"
	"main/log"
	shell_injection "main/vulnerabilities/shell-injection"
)

func OnPreShellExecuted() string {
	cmd := context.GetCmd()
	operation := context.GetFunctionName()
	if cmd == "" {
		return ""
	}

	log.Info("Got shell command: ", cmd)
	res := shell_injection.CheckContextForShellInjection(cmd, operation)
	if res != nil {
		go grpc.OnAttackDetected(*res)
		return attack.GetAttackDetectedAction(*res)
	}
	return ""
}
