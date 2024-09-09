package main

import (
	"main/attack"
	"main/grpc"
	"main/utils"
	path_traversal "main/vulnerabilities/path-traversal"
)

func OnPrePathAccessed(parameters map[string]interface{}) string {
	filename := utils.GetFromMap[string](parameters, "filename")
	filename2 := utils.GetFromMap[string](parameters, "filename2")
	operation := utils.GetFromMap[string](parameters, "operation")

	if filename == nil || operation == nil {
		return "{}"
	}

	res := path_traversal.CheckContextForPathTraversal(*filename, *operation, true)
	if res != nil {
		go grpc.OnAttackDetected(*res)
		return attack.GetAttackDetectedAction(*res)
	}

	if filename2 != nil && *filename2 != "" {
		res = path_traversal.CheckContextForPathTraversal(*filename2, *operation, true)
		if res != nil {
			go grpc.OnAttackDetected(*res)
			return attack.GetAttackDetectedAction(*res)
		}
	}
	return "{}"
}
