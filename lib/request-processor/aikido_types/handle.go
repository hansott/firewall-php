package aikido_types

type HandlerFunction func() string

type Method struct {
	ClassName  string
	MethodName string
}
