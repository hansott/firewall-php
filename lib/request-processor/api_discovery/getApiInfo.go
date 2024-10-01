package api_discovery

import (
	. "main/aikido_types"
	"main/context"
	"main/globals"
)

func GetApiInfo() *APISpec {
	if !globals.EnvironmentConfig.CollectApiSchema {
		return nil
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
			return nil
		}

		bodySchema := GetDataSchema(body, 0)

		bodyInfo = &APIBodyInfo{
			Type:   bodyType,
			Schema: bodySchema,
		}
	}

	// Check query data
	if query != nil && isObject(query) && len(query) > 0 {
		queryInfo = GetDataSchema(query, 0)
	}

	// Get Auth Info
	authInfo := GetApiAuthType()

	if bodyInfo == nil && queryInfo == nil && authInfo == nil {
		return nil
	}

	return &APISpec{
		Body:  bodyInfo,
		Query: queryInfo,
		Auth:  authInfo,
	}
}

func isObject(data interface{}) bool {
	// Helper function to determine if the data is an object (map in Go)
	_, ok := data.(map[string]interface{})
	return ok
}
