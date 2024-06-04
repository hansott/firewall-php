package main

import (
	"main/log"
	"main/utils"
	"strings"
)

func OnMethodExecutedPdoConstruct(parameters map[string]interface{}) string {
	//data_source := MustGetFromMap[string](parameters, "data_source")
	return "{}"
}

func OnMethodExecutedPdoQuery(parameters map[string]interface{}) string {
	query := utils.MustGetFromMap[string](parameters, "query")
	log.Info("Got PDO query: ", query)
	if strings.Contains(query, " OR 1=1") {
		return `{"action": "throw", "message": "Sql injection detected", "code": -1}`
	}
	return "{}"
}
