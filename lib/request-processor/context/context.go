package context

// #include "../../API.h"
import "C"
import (
	"main/helpers"
	"main/log"
	"main/utils"
	"net/url"
)

type CallbackFunction func(int) string

type ContextData struct {
	Callback                       CallbackFunction
	Method                         *string
	Route                          *string
	RouteParsed                    *string
	URL                            *string
	StatusCode                     *int
	IP                             *string
	IsIpBypassed                   *bool
	UserAgent                      *string
	UserId                         *string
	UserName                       *string
	BodyRaw                        *string
	BodyParsed                     *map[string]interface{}
	BodyParsedWithPathToPayload    *map[string]string
	QueryParsed                    *map[string]interface{}
	QueryParsedWithPathToPayload   *map[string]string
	CookiesParsed                  *map[string]interface{}
	CookiesParsedWithPathToPayload *map[string]string
	HeadersParsed                  *map[string]interface{}
	HeadersParsedWithPathToPayload *map[string]string
	OutgoingRequestHostname        *string
	OutgoingRequestPort            *int
	OutgoingRequestResolvedIp      *string
	PartialInterceptorResult       *utils.InterceptorResult
}

var Context ContextData

func Init(callback CallbackFunction) bool {
	if callback == nil {
		callback = Context.Callback
	}
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

func GetParsedRoute() string {
	return GetFromCache(ContextSetParsedRoute, &Context.RouteParsed)
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

func GetBodyRaw() string {
	return GetFromCache(ContextSetBody, &Context.BodyRaw)
}

func GetBodyParsed() map[string]interface{} {
	return GetFromCache(ContextSetBody, &Context.BodyParsed)
}

func GetQueryParsed() map[string]interface{} {
	return GetFromCache(ContextSetQuery, &Context.QueryParsed)
}

func GetCookiesParsed() map[string]interface{} {
	return GetFromCache(ContextSetCookies, &Context.CookiesParsed)
}

func GetHeadersParsed() map[string]interface{} {
	return GetFromCache(ContextSetHeaders, &Context.HeadersParsed)
}

func GetBodyParsedWithPathToPayload() map[string]string {
	return GetFromCache(ContextSetBody, &Context.BodyParsedWithPathToPayload)
}

func GetQueryParsedWithPathToPayload() map[string]string {
	return GetFromCache(ContextSetQuery, &Context.QueryParsedWithPathToPayload)
}

func GetCookiesParsedWithPathToPayload() map[string]string {
	return GetFromCache(ContextSetCookies, &Context.CookiesParsedWithPathToPayload)
}

func GetHeadersParsedWithPathToPayload() map[string]string {
	return GetFromCache(ContextSetHeaders, &Context.HeadersParsedWithPathToPayload)
}

func GetUserAgent() string {
	return GetFromCache(ContextSetUserAgent, &Context.UserAgent)
}

func GetUserId() string {
	return GetFromCache(ContextSetUserId, &Context.UserId)
}

func GetUserName() string {
	return GetFromCache(ContextSetUserName, &Context.UserName)
}

func GetHostNameAndPort(urlCallbackOption int) (string, int) {
	urlStr := Context.Callback(urlCallbackOption)
	urlParsed, err := url.Parse(urlStr)
	if err != nil {
		return "", 0
	}
	hostname := urlParsed.Hostname()
	portFromURL := helpers.GetPortFromURL(urlParsed)

	portStr := Context.Callback(C.OUTGOING_REQUEST_PORT)
	port := helpers.ParsePort(portStr)
	if port == 0 {
		port = portFromURL
	}
	return hostname, port
}

func GetOutgoingRequestHostnameAndPort() (string, int) {
	return GetHostNameAndPort(C.OUTGOING_REQUEST_URL)
}

func GetOutgoingRequestEffectiveHostnameAndPort() (string, int) {
	return GetHostNameAndPort(C.OUTGOING_REQUEST_EFFECTIVE_URL)
}

func GetOutgoingRequestResolvedIp() string {
	return Context.Callback(C.OUTGOING_REQUEST_RESOLVED_IP)
}

func GetPartialInterceptorResult() *utils.InterceptorResult {
	return Context.PartialInterceptorResult
}

func GetFunctionName() string {
	return Context.Callback(C.FUNCTION_NAME)
}

func GetCmd() string {
	return Context.Callback(C.CMD)
}

func GetFilename() string {
	return Context.Callback(C.FILENAME)
}

func GetFilename2() string {
	return Context.Callback(C.FILENAME2)
}

func GetSqlQuery() string {
	return Context.Callback(C.SQL_QUERY)
}

func GetSqlDialect() string {
	return Context.Callback(C.SQL_DIALECT)
}

func GetModule() string {
	return Context.Callback(C.MODULE)
}
