package context

// #include "../../API.h"
import "C"
import (
	"encoding/json"
	"fmt"
)

var TestContext map[string]string

func UnitTestsCallback(context_id int) string {
	switch context_id {
	case C.CONTEXT_REMOTE_ADDRESS:
		return TestContext["remoteAddress"]
	case C.CONTEXT_HTTPS:
		return TestContext["https"]
	case C.CONTEXT_METHOD:
		return TestContext["method"]
	case C.CONTEXT_ROUTE:
		return TestContext["route"]
	case C.CONTEXT_URL:
		return TestContext["url"]
	case C.CONTEXT_QUERY:
		return TestContext["query"]
	case C.CONTEXT_STATUS_CODE:
		return TestContext["statusCode"]
	case C.CONTEXT_HEADERS:
		return TestContext["headers"]
	case C.CONTEXT_HEADER_X_FORWARDED_FOR:
		return TestContext["xForwardedFor"]
	case C.CONTEXT_HEADER_USER_AGENT:
		return TestContext["userAgent"]
	case C.CONTEXT_COOKIES:
		return TestContext["cookies"]
	case C.CONTEXT_BODY:
		return TestContext["body"]
	}
	return ""
}

func LoadForUnitTests(context map[string]string) {
	Context.Callback = UnitTestsCallback
	TestContext = context
}

func GetJsonString(m map[string]interface{}) string {
	jsonStr, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error converting map to JSON:", err)
		return ""
	}

	return string(jsonStr)
}
