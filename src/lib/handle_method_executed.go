package main

type Method struct {
	ClassName  string
	MethodName string
}

type methodExecutedHandlersFn func(map[string]interface{}) string

var methodExecutedHandlers = map[Method]methodExecutedHandlersFn{
	Method{ClassName: "pdo", MethodName: "__construct"}: OnMethodExecutedPdoConstruct,
	Method{ClassName: "pdo", MethodName: "query"}:       OnMethodExecutedPdoQuery,
}

func OnMethodExecuted(data map[string]interface{}) string {
	className := MustGetFromMap[string](data, "class_name")
	methodName := MustGetFromMap[string](data, "method_name")
	parameters := MustGetFromMap[map[string]interface{}](data, "parameters")

	methodKey := Method{ClassName: className, MethodName: methodName}

	CheckIfKeyExists(methodExecutedHandlers, methodKey)

	return methodExecutedHandlers[methodKey](parameters)
}
