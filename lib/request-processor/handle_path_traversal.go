package main

import (
	"main/attack"
	"main/context"
	"main/grpc"
	path_traversal "main/vulnerabilities/path-traversal"
)

func OnPrePathAccessed() string {
	filename := context.GetFilename()
	filename2 := context.GetFilename2()
	operation := context.GetFunctionName()

	if filename == "" || operation == "" {
		return "{}"
	}

	res := path_traversal.CheckContextForPathTraversal(filename, operation, true)
	if res != nil {
		go grpc.OnAttackDetected(attack.GetAttackDetectedProto(*res))
		return attack.GetAttackDetectedAction(*res)
	}

	if filename2 != "" {
		res = path_traversal.CheckContextForPathTraversal(filename2, operation, true)
		if res != nil {
			go grpc.OnAttackDetected(attack.GetAttackDetectedProto(*res))
			return attack.GetAttackDetectedAction(*res)
		}
	}
	return "{}"
}
