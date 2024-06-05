package main

import (
	"main/globals"
	"main/log"
	"main/utils"
)

func OnFunctionExecutedCurl(parameters map[string]interface{}) string {
	url := utils.GetFromMap[string](parameters, "url")
	if url == nil {
		return "{}"
	}
	domain := utils.GetDomain(*url)
	globals.OutgoingHostnames[domain] = false
	log.Info("Got domain: ", domain)
	return "{}"
}
