package main

import "fmt"

var outgoingHostnames = map[string]bool{}

func OnFunctionExecutedCurl(parameters map[string]interface{}) string {
	url := GetFromMap[string](parameters, "url")
	if url == nil {
		return "{}"
	}
	domain := GetDomain(*url)
	outgoingHostnames[domain] = false
	fmt.Println("[AIKIDO-GO] Got domain:", domain)
	return "{}"
}
