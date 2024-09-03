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

func ContextCacheParsedStrings(context_id int, m_str **string, m **map[string]string, parseFunc ParseFunction) {
	if *m != nil {
		return
	}
	context_data := Context.Callback(context_id)
	if m_str != nil {
		*m_str = &context_data
	}
	parsed := parseFunc(context_data)
	strings := helpers.ExtractStringsFromUserInput(parsed, []helpers.PathPart{})
	*m = &strings
}

func ContextCacheString(context_id int, m **string) {
	if *m != nil {
		return
	}
	temp := Context.Callback(context_id)
	*m = &temp
}

func ContextCacheBody() {
	ContextCacheParsedStrings(C.CONTEXT_BODY, &Context.Body, &Context.BodyParsed, utils.ParseBody)
}

func ContextCacheQuery() {
	ContextCacheParsedStrings(C.CONTEXT_QUERY, nil, &Context.QueryParsed, utils.ParseQuery)
}

func ContextCacheCookies() {
	ContextCacheParsedStrings(C.CONTEXT_COOKIES, nil, &Context.CookiesParsed, utils.ParseCookies)
}

func ContextCacheHeaders() {
	ContextCacheParsedStrings(C.CONTEXT_HEADERS, nil, &Context.HeadersParsed, utils.ParseHeaders)
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
