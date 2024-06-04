package main

import (
	"main/log"
	"main/utils"
)

var outgoingHostnames = map[string]bool{}

func OnFunctionExecutedCurl(parameters map[string]interface{}) string {
	url := utils.GetFromMap[string](parameters, "url")
	if url == nil {
		return "{}"
	}
	domain := utils.GetDomain(*url)
	outgoingHostnames[domain] = false
	log.Info("Got domain: ", domain)
	return "{}"
}
