package main

import "main/log"

var shellCommands = map[string]bool{}

func OnFunctionExecutedShell(parameters map[string]interface{}) string {
	cmd := GetFromMap[string](parameters, "cmd")
	if cmd == nil {
		return "{}"
	}
	shellCommands[*cmd] = false
	log.Info("Got shell command: ", *cmd)
	return "{}"
}
