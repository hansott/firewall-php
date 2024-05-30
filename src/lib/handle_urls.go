package main

import "main/log"

var outgoingHostnames = map[string]bool{}

func OnFunctionExecutedCurl(parameters map[string]interface{}) string {
	url := GetFromMap[string](parameters, "url")
	if url == nil {
		return "{}"
	}
	domain := GetDomain(*url)
	outgoingHostnames[domain] = false
	log.Info("Got domain: ", domain)
	return "{}"
}
