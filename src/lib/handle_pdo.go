package main

import "fmt"

func OnMethodExecutedPdoConstruct(parameters map[string]interface{}) string {
	data_source := GetFromMap[string](parameters, "data_source")
	fmt.Println("[AIKIDO-GO] Got PDO data_source:", *data_source)
	return "{}"
}

func OnMethodExecutedPdoQuery(parameters map[string]interface{}) string {
	query := GetFromMap[string](parameters, "query")
	fmt.Println("[AIKIDO-GO] Got PDO query:", *query)
	return "{}"
}
