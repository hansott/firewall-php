package context

// #include "../../ContextCallback.c"
import "C"
import "main/log"

type CallbackFunction func(int) string

type ContextData struct {
	Callback      CallbackFunction
	Method        *string
	Route         *string
	URL           *string
	StatusCode    *int
	IP            *string
	IsIpBypassed  *bool
	UserAgent     *string
	UserId        *string
	Body          *string
	BodyParsed    *map[string]string
	Query         *string
	QueryParsed   *map[string]string
	Cookies       *string
	CookiesParsed *map[string]string
	Headers       *map[string]interface{}
	HeadersParsed *map[string]string
}

var Context ContextData

func Init(callback CallbackFunction) bool {
	Context = ContextData{
		Callback: callback,
	}
	return true
}

func GetFromCache[T any](fetchDataFn func(), s **T) T {
	if fetchDataFn != nil {
		fetchDataFn()
	}
	if *s == nil {
		var t T
		log.Warnf("Error getting from cache. Returning default value %v...", t)
		return t
	}
	return **s
}

func GetIp() string {
	return GetFromCache(ContextSetIp, &Context.IP)
}

func GetMethod() string {
	return GetFromCache(ContextSetMethod, &Context.Method)
}

func GetRoute() string {
	return GetFromCache(ContextSetRoute, &Context.Route)
}

func GetUrl() string {
	return GetFromCache(ContextSetUrl, &Context.URL)
}

func GetStatusCode() int {
	return GetFromCache(ContextSetStatusCode, &Context.StatusCode)
}

func IsIpBypassed() bool {
	return GetFromCache(ContextSetIsIpBypassed, &Context.IsIpBypassed)
}

func GetBody() string {
	return GetFromCache(ContextSetBody, &Context.Body)
}

func GetBodyParsed() map[string]string {
	return GetFromCache(ContextSetBody, &Context.BodyParsed)
}

func GetQuery() string {
	return GetFromCache(ContextSetQuery, &Context.Query)
}

func GetQueryParsed() map[string]string {
	return GetFromCache(ContextSetQuery, &Context.QueryParsed)
}

func GetCookies() string {
	return GetFromCache(ContextSetCookies, &Context.Cookies)
}

func GetCookiesParsed() map[string]string {
	return GetFromCache(ContextSetCookies, &Context.CookiesParsed)
}

func GetHeaders() map[string]interface{} {
	return GetFromCache(ContextSetHeaders, &Context.Headers)
}

func GetHeadersParsed() map[string]string {
	return GetFromCache(ContextSetHeaders, &Context.HeadersParsed)
}

func GetUserAgent() string {
	return GetFromCache(ContextSetUserAgent, &Context.UserAgent)
}

func GetUserId() string {
	return GetFromCache(nil, &Context.UserId)
}
