package aikido_types

type EventFunctionExecutedFn func(map[string]interface{}) string

type FunctionExecutedHandlersFn func(map[string]interface{}) string

type MethodExecutedHandlersFn func(map[string]interface{}) string

type Method struct {
	ClassName  string
	MethodName string
}
