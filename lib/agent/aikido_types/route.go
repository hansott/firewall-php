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
	// Type of this property, e.g., "string", "number", "object", "array", "null"
	Type interface{} `json:"type"`
	// Whether this property is optional
	Optional *bool `json:"optional,omitempty"`
	// Map of properties for an object containing the DataSchema for each property
	Properties map[string]*DataSchema `json:"properties,omitempty"`
	// Data schema for the items of an array
	Items *DataSchema `json:"items,omitempty"`
}

type APISpec struct {
	Body  *APIBodyInfo
	Query *DataSchema
	Auth  []APIAuthType
}

type APIBodyInfo struct {
	Type   BodyDataType
	Schema *DataSchema
}

// BodyDataType represents the type of the body data.
type BodyDataType string

const (
	JSON           BodyDataType = "json"
	FormURLEncoded BodyDataType = "form-urlencoded"
	FormData       BodyDataType = "form-data"
	XML            BodyDataType = "xml"
	Undefined      BodyDataType = ""
)
