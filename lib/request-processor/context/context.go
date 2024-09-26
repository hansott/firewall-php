package context

// #include "../../API.h"
import "C"
import (
	"main/helpers"
	"main/log"
	"main/utils"
)

type CallbackFunction func(int) string

type ContextData struct {
	Callback                  CallbackFunction
	Method                    *string
	Route                     *string
	RouteParsed               *string
	URL                       *string
	StatusCode                *int
	IP                        *string
	IsIpBypassed              *bool
	UserAgent                 *string
	UserId                    *string
	UserName                  *string
	Body                      *string
	BodyParsed                *map[string]string
	Query                     *string
	QueryParsed               *map[string]string
	Cookies                   *string
	CookiesParsed             *map[string]string
	Headers                   *map[string]interface{}
	HeadersParsed             *map[string]string
	OutgoingRequestHostname   *string
	OutgoingRequestPort       *int
	OutgoingRequestResolvedIp *string
	PartialInterceptorResult  *utils.InterceptorResult
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

func GetBody() string {
	return GetFromCache(ContextSetBody, &Context.Body)
}

func GetParsedBody() map[string]string {
	return GetFromCache(ContextSetBody, &Context.BodyParsed)
}

func GetQuery() string {
	return GetFromCache(ContextSetQuery, &Context.Query)
}

func GetParsedQuery() map[string]string {
	return GetFromCache(ContextSetQuery, &Context.QueryParsed)
}

func GetCookies() string {
	return GetFromCache(ContextSetCookies, &Context.Cookies)
}

func GetParsedCookies() map[string]string {
	return GetFromCache(ContextSetCookies, &Context.CookiesParsed)
}

func GetHeaders() map[string]interface{} {
	return GetFromCache(ContextSetHeaders, &Context.Headers)
}

func GetParsedHeaders() map[string]string {
	return GetFromCache(ContextSetHeaders, &Context.HeadersParsed)
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

func GetOutgoingRequestHostname() string {
	return GetFromCache(ContextSetOutgoingRequestHostnameAndPort, &Context.OutgoingRequestHostname)
}

func GetOutgoingRequestPort() int {
	portStr := Context.Callback(C.OUTGOING_REQUEST_PORT)
	port := helpers.ParsePort(portStr)
	if port != 0 {
		return port
	}
	return GetFromCache(ContextSetOutgoingRequestHostnameAndPort, &Context.OutgoingRequestPort)
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
