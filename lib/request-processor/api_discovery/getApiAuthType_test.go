package api_discovery

import (
	"main/context"
	"main/ipc/protos"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test for detecting authorization header
func TestDetectAuthorizationHeader(t *testing.T) {
	assert := assert.New(t)

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{
			"authorization": "Bearer token",
		}),
	})
	assert.Equal([]*protos.APIAuthType{
		{Type: "http", Scheme: "bearer"},
	}, GetApiAuthType())
	context.UnloadForUnitTests()

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{
			"authorization": "Basic base64",
		}),
	})
	assert.Equal([]*protos.APIAuthType{
		{Type: "http", Scheme: "basic"},
	}, GetApiAuthType())
	context.UnloadForUnitTests()

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{
			"authorization": "custom",
		}),
	})
	assert.Equal([]*protos.APIAuthType{
		{Type: "apiKey", In: "header", Name: "Authorization"},
	}, GetApiAuthType())
	context.UnloadForUnitTests()
}

// Test for detecting API keys
func TestDetectApiKeys(t *testing.T) {
	assert := assert.New(t)

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{
			"x-api-key": "token",
		}),
	})
	assert.Equal([]*protos.APIAuthType{
		{Type: "apiKey", In: ("header"), Name: ("x-api-key")},
	}, GetApiAuthType())
	context.UnloadForUnitTests()

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{
			"api-key": "token",
		}),
	})
	assert.Equal([]*protos.APIAuthType{
		{Type: "apiKey", In: ("header"), Name: ("api-key")},
	}, GetApiAuthType())
	context.UnloadForUnitTests()
}

// Test for detecting auth cookies
func TestDetectAuthCookies(t *testing.T) {
	assert := assert.New(t)

	context.LoadForUnitTests(map[string]string{
		"cookies": "api-key=token",
	})
	assert.Equal([]*protos.APIAuthType{
		{Type: "apiKey", In: ("cookie"), Name: ("api-key")},
	}, GetApiAuthType())
	context.UnloadForUnitTests()

	context.LoadForUnitTests(map[string]string{
		"cookies": "session=test",
	})
	assert.Equal([]*protos.APIAuthType{
		{Type: "apiKey", In: ("cookie"), Name: ("session")},
	}, GetApiAuthType())
	context.UnloadForUnitTests()
}

// Test for no authentication
func TestNoAuth(t *testing.T) {
	assert := assert.New(t)

	context.LoadForUnitTests(map[string]string{})
	assert.Empty(GetApiAuthType())
	context.UnloadForUnitTests()

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{}),
	})
	assert.Empty(GetApiAuthType())
	context.UnloadForUnitTests()

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{
			"authorization": "",
		}),
	})
	assert.Empty(GetApiAuthType())
	context.UnloadForUnitTests()
}
