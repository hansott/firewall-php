package main

import "main/log"

func OnMethodExecutedPdoConstruct(parameters map[string]interface{}) string {
	//data_source := MustGetFromMap[string](parameters, "data_source")
	return "{}"
}

func OnMethodExecutedPdoQuery(parameters map[string]interface{}) string {
	query := MustGetFromMap[string](parameters, "query")
	log.Info("Got PDO query: ", query)
	//if strings.Contains(query, "users") {
	//	return `{"action": "throw", "message": "test", "code": -1}`
	//}
	return "{}"
}
