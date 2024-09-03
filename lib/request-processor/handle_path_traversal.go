package main

import (
	"main/grpc"
	"main/utils"
	path_traversal "main/vulnerabilities/path-traversal"
)

func OnPathAccessed(parameters map[string]interface{}) string {
	filename := utils.GetFromMap[string](parameters, "filename")
	filename2 := utils.GetFromMap[string](parameters, "filename2")
	operation := utils.GetFromMap[string](parameters, "operation")

	if filename == nil || operation == nil {
		return "{}"
	}

	res := path_traversal.CheckContextForPathTraversal(*filename, *operation, true)
	if res != nil {
		go grpc.OnAttackDetected(*res)
		return `{"action": "throw", "message": "Path traversal detected", "code": -1}`
	}

	if filename2 != nil && *filename2 != "" {
		res = path_traversal.CheckContextForPathTraversal(*filename2, *operation, true)
		if res != nil {
			go grpc.OnAttackDetected(*res)
			return `{"action": "throw", "message": "Path traversal detected", "code": -1}`
		}
	}
	return "{}"
}
