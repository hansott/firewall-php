package main

import (
	"main/log"
	"main/utils"
)

var shellCommands = map[string]bool{}

func OnFunctionExecutedShell(parameters map[string]interface{}) string {
	cmd := utils.GetFromMap[string](parameters, "cmd")
	if cmd == nil {
		return "{}"
	}
	shellCommands[*cmd] = false
	log.Info("Got shell command: ", *cmd)
	return "{}"
}
