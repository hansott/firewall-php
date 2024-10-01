package api_discovery

import (
	. "main/aikido_types"
	"main/context"
	"main/globals"
)

/**
 * Get body data type and schema from context.
 * Returns nil if the body is not an object or if the body type could not be determined.
 */
func GetApiInfo() (*APISpec, error) {
	if !globals.EnvironmentConfig.CollectApiSchema {
		return nil, nil
	}

	var bodyInfo *APIBodyInfo
	var queryInfo *DataSchema

	body := context.GetBodyParsed()
	headers := context.GetHeadersParsed()
	query := context.GetQueryParsed()

	// Check body data
	if body != nil && isObject(body) && len(body) > 0 {
		bodyType := getBodyDataType(headers)
		if bodyType == Undefined {
			return nil, nil
		}

		bodySchema := GetDataSchema(body, 0)

		bodyInfo = &APIBodyInfo{
			Type:   bodyType,
			Schema: bodySchema,
		}
	}

	// Check query data
	if query != nil && isObject(query) && len(query) > 0 {
		querySchema := GetDataSchema(query, 0)
		queryInfo = querySchema
	}

	// Get Auth Info
	authInfo := GetApiAuthType()

	if bodyInfo == nil && queryInfo == nil && authInfo == nil {
		return nil, nil
	}

	return &APISpec{
		Body:  bodyInfo,
		Query: queryInfo,
		Auth:  authInfo,
	}, nil
}

func isObject(data interface{}) bool {
	// Helper function to determine if the data is an object (map in Go)
	_, ok := data.(map[string]interface{})
	return ok
}
