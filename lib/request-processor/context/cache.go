package context

// #include "../../ContextCallback.c"
import "C"
import (
	"main/helpers"
	"main/log"
	"main/utils"
	"strconv"
)

type ParseFunction func(string) map[string]interface{}

func ContextCacheMap(contextId int, rawDataPtr **string, parsedPtr **map[string]interface{}, stringsPtr **map[string]string, parseFunc ParseFunction) {
	if *stringsPtr != nil {
		return
	}
	contextData := Context.Callback(contextId)
	if rawDataPtr != nil {
		*rawDataPtr = &contextData
	}
	parsed := parseFunc(contextData)
	if parsedPtr != nil {
		*parsedPtr = &parsed
	}
	strings := helpers.ExtractStringsFromUserInput(parsed, []helpers.PathPart{})
	*stringsPtr = &strings
}

func ContextCacheString(context_id int, m **string) {
	if *m != nil {
		return
	}
	temp := Context.Callback(context_id)
	*m = &temp
}

func ContextCacheBody() {
	ContextCacheMap(C.CONTEXT_BODY, &Context.Body, nil, &Context.BodyParsed, utils.ParseBody)
}

func ContextCacheQuery() {
	ContextCacheMap(C.CONTEXT_QUERY, &Context.Query, nil, &Context.QueryParsed, utils.ParseQuery)
}

func ContextCacheCookies() {
	ContextCacheMap(C.CONTEXT_COOKIES, &Context.Cookies, nil, &Context.CookiesParsed, utils.ParseCookies)
}

func ContextCacheHeaders() {
	ContextCacheMap(C.CONTEXT_HEADERS, nil, &Context.Headers, &Context.HeadersParsed, utils.ParseHeaders)
}

func ContextCacheStatusCode() {
	if Context.StatusCode != nil {
		return
	}
	status_code_str := Context.Callback(C.CONTEXT_STATUS_CODE)
	status_code, err := strconv.Atoi(status_code_str)
	if err != nil {
		log.Warnf("Error parsing status code %v: %v", status_code_str, err)
		return
	}
	Context.StatusCode = &status_code
}

func ContextCacheRoute() {
	ContextCacheString(C.CONTEXT_ROUTE, &Context.Route)
}

func ContextCacheMethod() {
	ContextCacheString(C.CONTEXT_METHOD, &Context.Method)
}

func ContextCacheUrl() {
	ContextCacheString(C.CONTEXT_URL, &Context.URL)
}

func ContextCacheUserAgent() {
	ContextCacheString(C.CONTEXT_HEADER_USER_AGENT, &Context.UserAgent)
}

func ContextCacheIp() {
	if Context.IP != nil {
		return
	}
	remoteAddress := Context.Callback(C.CONTEXT_REMOTE_ADDRESS)
	xForwardedFor := Context.Callback(C.CONTEXT_HEADER_X_FORWARDED_FOR)
	ip := utils.GetIpFromRequest(remoteAddress, xForwardedFor)
	Context.IP = &ip
}

func ContextCacheIsIpBypassed() {
	ContextCacheIp()
	if Context.IsIpBypassed != nil {
		return
	}

	if Context.IP == nil {
		return
	}
	isIpBypassed := utils.IsIpBypassed(*Context.IP)
	Context.IsIpBypassed = &isIpBypassed
}

func ContextCacheUserId(userId string) {
	Context.UserId = &userId
}
