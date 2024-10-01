package api_discovery

import (
	. "main/aikido_types"
	"main/context"
	"main/utils"
	"strings"
)

// Common API key header and cookie names
var commonApiKeyHeaderNames = []string{
	"x-api-key",
	"api-key",
	"apikey",
	"x-token",
	"token",
}

var commonAuthCookieNames = append([]string{
	"auth",
	"session",
	"jwt",
	"token",
	"sid",
	"connect.sid",
	"auth_token",
	"access_token",
	"refresh_token",
}, commonApiKeyHeaderNames...)

// GetApiAuthType returns the authentication type of the API request.
// Returns nil if the authentication type could not be determined.
func GetApiAuthType() []APIAuthType {
	var result []APIAuthType

	headers := context.GetHeadersParsed()

	// Check the Authorization header
	authHeader, authHeaderExists := headers["authorization"].(string)
	if authHeaderExists {
		authHeaderType := getAuthorizationHeaderType(authHeader)
		if authHeaderType != nil {
			result = append(result, *authHeaderType)
		}
	}

	result = append(result, findApiKeys()...)
	return result
}

// getAuthorizationHeaderType returns the authentication type from the Authorization header.
func getAuthorizationHeaderType(authHeader string) *APIAuthType {
	if len(authHeader) == 0 {
		return nil
	}
	if strings.Contains(authHeader, " ") {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 {
			authType := parts[0]
			if isHTTPAuthScheme(authType) {
				scheme := strings.ToLower(authType)
				return &APIAuthType{
					Type:   "http",
					Scheme: &scheme,
				}
			}
		}
	}

	// Default to apiKey if the auth type is not recognized
	name := "Authorization"
	return &APIAuthType{
		Type: "apiKey",
		In:   utils.StringPointer("header"),
		Name: &name,
	}
}

// findApiKeys searches for API keys in headers and cookies.
func findApiKeys() []APIAuthType {
	var result []APIAuthType

	headers := context.GetHeadersParsed()
	cookies := context.GetCookiesParsed()
	for header_index, header := range commonApiKeyHeaderNames {
		if value, exists := headers[header]; exists && value != "" {
			result = append(result, APIAuthType{
				Type: "apiKey",
				In:   utils.StringPointer("header"),
				Name: &commonApiKeyHeaderNames[header_index],
			})
		}
	}

	if len(cookies) > 0 {
		for cookieName := range cookies {
			lowerCookieName := strings.ToLower(cookieName)
			if contains(commonAuthCookieNames, lowerCookieName) {
				cookieNameCopy := cookieName
				result = append(result, APIAuthType{
					Type: "apiKey",
					In:   utils.StringPointer("cookie"),
					Name: &cookieNameCopy,
				})
			}
		}
	}

	return result
}

// contains checks if a string exists in a slice.
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// isHTTPAuthScheme checks if the given string is a valid HTTP authentication scheme.
// You will need to implement this function similar to the TypeScript helper function.
func isHTTPAuthScheme(scheme string) bool {
	// You can add proper logic here to check the scheme, e.g., "basic", "bearer", etc.
	// For example:
	allowedSchemes := []string{"basic", "bearer"}
	return contains(allowedSchemes, strings.ToLower(scheme))
}
