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
	return GetFromCache(ContextCacheIp, &Context.IP)
}

func GetMethod() string {
	return GetFromCache(ContextCacheMethod, &Context.Method)
}

func GetRoute() string {
	return GetFromCache(ContextCacheRoute, &Context.Route)
}

func GetUrl() string {
	return GetFromCache(ContextCacheUrl, &Context.URL)
}

func GetStatusCode() int {
	return GetFromCache(ContextCacheStatusCode, &Context.StatusCode)
}

func IsIpBypassed() bool {
	return GetFromCache(ContextCacheIsIpBypassed, &Context.IsIpBypassed)
}

func GetBody() string {
	return GetFromCache(ContextCacheBody, &Context.Body)
}

func GetBodyParsed() map[string]string {
	return GetFromCache(ContextCacheBody, &Context.BodyParsed)
}

func GetQuery() string {
	return GetFromCache(ContextCacheQuery, &Context.Query)
}

func GetQueryParsed() map[string]string {
	return GetFromCache(ContextCacheQuery, &Context.QueryParsed)
}

func GetCookies() string {
	return GetFromCache(ContextCacheCookies, &Context.Cookies)
}

func GetCookiesParsed() map[string]string {
	return GetFromCache(ContextCacheCookies, &Context.CookiesParsed)
}

func GetHeaders() map[string]interface{} {
	return GetFromCache(ContextCacheHeaders, &Context.Headers)
}

func GetHeadersParsed() map[string]string {
	return GetFromCache(ContextCacheHeaders, &Context.HeadersParsed)
}

func GetUserAgent() string {
	return GetFromCache(ContextCacheUserAgent, &Context.UserAgent)
}

func GetUserId() string {
	return GetFromCache(nil, &Context.UserId)
}
