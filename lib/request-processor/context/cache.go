package context

// #include "../../ContextCallback.c"
import "C"
import (
	"main/helpers"
	"main/utils"
	"strconv"
)

type ParseFunction func(string) map[string]interface{}

func ContextCacheParsedStrings(context_id int, m *map[string]string, parseFunc ParseFunction) {
	if Context.Query != nil {
		return
	}
	context_data := Context.Callback(context_id)
	parsed := parseFunc(context_data)
	strings := helpers.ExtractStringsFromUserInput(parsed, []helpers.PathPart{})
	*m = strings
}

func ContextCacheBody() {
	ContextCacheParsedStrings(C.CONTEXT_BODY, Context.Body, utils.ParseBody)
}

func ContextCacheQuery() {
	ContextCacheParsedStrings(C.CONTEXT_QUERY, Context.Query, utils.ParseQuery)
}

func ContextCacheCookies() {
	ContextCacheParsedStrings(C.CONTEXT_COOKIES, Context.Cookies, utils.ParseCookies)
}

func ContextCacheStatusCode() {
	if Context.StatusCode != nil {
		return
	}
	status_code_str := Context.Callback(C.CONTEXT_STATUS_CODE)

	status_code, err := strconv.Atoi(status_code_str)
	if err != nil {
		return
	}
	Context.StatusCode = &status_code
}

func ContextCacheRoute() {
	if Context.Route != nil {
		return
	}
	route := Context.Callback(C.CONTEXT_ROUTE)
	Context.Route = &route
}

func ContextCacheMethod() {
	if Context.Method != nil {
		return
	}
	method := Context.Callback(C.CONTEXT_METHOD)
	Context.Method = &method
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
	Context.UserId = userId
}
