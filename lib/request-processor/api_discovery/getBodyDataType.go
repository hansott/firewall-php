package api_discovery

import (
	. "main/aikido_types"
	"strings"
)

// getBodyDataType tries to determine the type of the body data based on the content type header.
func getBodyDataType(headers map[string]interface{}) string {
	if headers == nil {
		return Undefined
	}

	contentType, exists := headers["content_type"].(string)
	if !exists {
		return Undefined
	}

	// Check if contentType has multiple values (comma separated or otherwise)
	// and use the first one.
	if strings.Contains(contentType, ",") {
		contentType = strings.Split(contentType, ",")[0]
	}

	if IsJsonContentType(contentType) {
		return JSON
	}

	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		return FormURLEncoded
	}

	if strings.HasPrefix(contentType, "multipart/form-data") {
		return FormData
	}

	if strings.Contains(contentType, "xml") {
		return XML
	}

	return Undefined
}
