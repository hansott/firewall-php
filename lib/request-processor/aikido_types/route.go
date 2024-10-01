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
	Items      *DataSchema            `json:"items,omitempty"`      // Data schema for the items of an array
}

type APIBodyInfo struct {
	Type   string      `json:"type"`
	Schema *DataSchema `json:"schema,omitempty"`
}

type APISpec struct {
	Body  *APIBodyInfo  `json:"body,omitempty"`
	Query *DataSchema   `json:"query,omitempty"`
	Auth  []APIAuthType `json:"auth,omitempty"`
}

const (
	JSON           = "json"
	FormURLEncoded = "form-urlencoded"
	FormData       = "form-data"
	XML            = "xml"
	Undefined      = ""
)
