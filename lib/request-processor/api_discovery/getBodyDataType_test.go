package api_discovery

import (
	"strings"
	"testing"
)

func determineType(contentType string) string {
	contentType = strings.ToLower(contentType)
	switch {
	case strings.Contains(contentType, "json"):
		return "json"
	case contentType == "application/x-www-form-urlencoded":
		return "form-urlencoded"
	case contentType == "multipart/form-data":
		return "form-data"
	case strings.HasPrefix(contentType, "text/xml"):
		return "xml"
	default:
		return ""
	}
}

func TestGetBodyDataType(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]interface{}
		expected string
	}{
		{"JSON content-type", map[string]interface{}{"content-type": "application/json"}, "json"},
		{"API JSON content-type", map[string]interface{}{"content-type": "application/vnd.api+json"}, "json"},
		{"CSP report content-type", map[string]interface{}{"content-type": "application/csp-report"}, "json"},
		{"X JSON content-type", map[string]interface{}{"content-type": "application/x-json"}, "json"},
		{"Form-urlencoded content-type", map[string]interface{}{"content-type": "application/x-www-form-urlencoded"}, "form-urlencoded"},
		{"Multipart form-data content-type", map[string]interface{}{"content-type": "multipart/form-data"}, "form-data"},
		{"XML content-type", map[string]interface{}{"content-type": "text/xml"}, "xml"},
		{"HTML content-type", map[string]interface{}{"content-type": "text/html"}, ""},
		{"Multiple content-types", map[string]interface{}{"content-type": "application/json"}, "json"},
		{"Nonexistent content-type", map[string]interface{}{"x-test": "abc"}, ""},
		{"Null input", nil, ""},
		{"Empty headers", map[string]interface{}{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBodyDataType(tt.headers)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
