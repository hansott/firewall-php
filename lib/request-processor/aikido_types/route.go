package aikido_types

type RouteData struct {
	Method  string
	Path    string
	APISpec APISpec
}

// APIAuthType represents the authentication type for the API.
type APIAuthType struct {
	Type         string
	Scheme       *string
	In           *string
	Name         *string
	BearerFormat *string
}

// DataSchema defines the structure of a schema
type DataSchema struct {
	Type       string                 `json:"type"`                 // Type of this property, e.g., "string", "number", "object", "array", "null"
	Properties map[string]*DataSchema `json:"properties,omitempty"` // Map of properties for an object containing the DataSchema for each property

	// Data schema for the items of an array
	Items *DataSchema `json:"items,omitempty"`
}

type APISpec struct {
	Body  *APIBodyInfo
	Query *DataSchema
	Auth  []APIAuthType
}

type APIBodyInfo struct {
	Type   string
	Schema *DataSchema
}

const (
	JSON           = "json"
	FormURLEncoded = "form-urlencoded"
	FormData       = "form-data"
	XML            = "xml"
	Undefined      = ""
)
