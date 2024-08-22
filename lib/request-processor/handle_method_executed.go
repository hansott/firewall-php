package main

import (
	. "main/aikido_types"
	"main/utils"
)

var methodExecutedHandlers = map[Method]HandlerFunction{
	{ClassName: "pdo", MethodName: "__construct"}: OnMethodExecutedPdoConstruct,
	{ClassName: "pdo", MethodName: "query"}:       OnMethodExecutedPdoQuery,
}

func OnMethodExecuted(data map[string]interface{}) string {
	className := utils.MustGetFromMap[string](data, "class_name")
	methodName := utils.MustGetFromMap[string](data, "method_name")
	parameters := utils.MustGetFromMap[map[string]interface{}](data, "parameters")

	methodKey := Method{ClassName: className, MethodName: methodName}

	utils.KeyMustExist(methodExecutedHandlers, methodKey)

	return methodExecutedHandlers[methodKey](parameters)
}
