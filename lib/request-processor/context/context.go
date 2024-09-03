package context

// #include "../../ContextCallback.c"
import "C"

type CallbackFunction func(int) string

type ContextData struct {
	Callback     CallbackFunction
	Method       *string
	Route        *string
	StatusCode   *int
	IP           *string
	IsIpBypassed *bool
	UserId       string
	Body         *map[string]string
	Query        *map[string]string
	Cookies      *map[string]string
	Headers      *map[string]string
}

var Context ContextData

func Init(callback CallbackFunction) bool {
	Context = ContextData{
		Callback: callback,
	}
	return true
}

type StoreToCacheFn func()

func GetFromCache[T any](fn StoreToCacheFn, s *T) T {
	fn()
	if s == nil {
		var t T
		return t
	}
	return *s
}

func GetIp() string {
	return GetFromCache(ContextCacheIp, Context.IP)
}

func GetMethod() string {
	return GetFromCache(ContextCacheMethod, Context.Method)
}

func GetRoute() string {
	return GetFromCache(ContextCacheRoute, Context.Route)
}

func GetStatusCode() int {
	return GetFromCache(ContextCacheStatusCode, Context.StatusCode)
}

func IsIpBypassed() bool {
	return GetFromCache(ContextCacheIsIpBypassed, Context.IsIpBypassed)
}

func GetBody() map[string]string {
	return GetFromCache(ContextCacheBody, Context.Body)
}

func GetQuery() map[string]string {
	return GetFromCache(ContextCacheQuery, Context.Query)
}

func GetCookies() map[string]string {
	return GetFromCache(ContextCacheCookies, Context.Cookies)
}

func ContextGetUserId() string {
	return Context.UserId
}
