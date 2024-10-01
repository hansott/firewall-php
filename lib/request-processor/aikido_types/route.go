package aikido_types

import "main/ipc/protos"

type RouteData struct {
	Method  string
	Path    string
	APISpec protos.APISpec
}

const (
	JSON           = "json"
	FormURLEncoded = "form-urlencoded"
	FormData       = "form-data"
	XML            = "xml"
	Undefined      = ""
)
