package main

import "fmt"

var shellCommands = map[string]bool{}

func OnFunctionExecutedShell(parameters map[string]interface{}) string {
	cmd := GetFromMap[string](parameters, "cmd")
	if cmd == nil {
		return "{}"
	}
	shellCommands[*cmd] = false
	fmt.Println("[AIKIDO-GO] Got shell command:", *cmd)
	return "{}"
}
