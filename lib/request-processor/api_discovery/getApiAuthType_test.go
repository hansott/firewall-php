package api_discovery

import (
	. "main/aikido_types"
	"main/context"
	"main/utils"
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
	assert.Equal([]APIAuthType{
		{Type: "http", Scheme: utils.StringPointer("bearer")},
	}, GetApiAuthType())

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{
			"authorization": "Basic base64",
		}),
	})
	assert.Equal([]APIAuthType{
		{Type: "http", Scheme: utils.StringPointer("basic")},
	}, GetApiAuthType())

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{
			"authorization": "Basic base64",
		}),
	})

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{
			"authorization": "custom",
		}),
	})
	assert.Equal([]APIAuthType{
		{Type: "apiKey", In: utils.StringPointer("header"), Name: utils.StringPointer("Authorization")},
	}, GetApiAuthType())
}

// Test for detecting API keys
func TestDetectApiKeys(t *testing.T) {
	assert := assert.New(t)

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{
			"x-api-key": "token",
		}),
	})
	assert.Equal([]APIAuthType{
		{Type: "apiKey", In: utils.StringPointer("header"), Name: utils.StringPointer("x-api-key")},
	}, GetApiAuthType())

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{
			"api-key": "token",
		}),
	})
	assert.Equal([]APIAuthType{
		{Type: "apiKey", In: utils.StringPointer("header"), Name: utils.StringPointer("api-key")},
	}, GetApiAuthType())
}

// Test for detecting auth cookies
func TestDetectAuthCookies(t *testing.T) {
	assert := assert.New(t)

	context.LoadForUnitTests(map[string]string{
		"cookies": "api-key=token",
	})
	assert.Equal([]APIAuthType{
		{Type: "apiKey", In: utils.StringPointer("cookie"), Name: utils.StringPointer("api-key")},
	}, GetApiAuthType())

	context.LoadForUnitTests(map[string]string{
		"cookies": "session=test",
	})
	assert.Equal([]APIAuthType{
		{Type: "apiKey", In: utils.StringPointer("cookie"), Name: utils.StringPointer("session")},
	}, GetApiAuthType())
}

// Test for no authentication
func TestNoAuth(t *testing.T) {
	assert := assert.New(t)

	context.LoadForUnitTests(map[string]string{})
	assert.Empty(GetApiAuthType())

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{}),
	})
	assert.Empty(GetApiAuthType())

	context.LoadForUnitTests(map[string]string{
		"headers": context.GetJsonString(map[string]interface{}{
			"authorization": "",
		}),
	})
	assert.Empty(GetApiAuthType())
}
