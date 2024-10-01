package context

type Source struct {
	Name     string
	CacheGet func() map[string]string
}

var SOURCES = []Source{
	{"body", GetBodyParsedWithPathToPayload},
	{"query", GetQueryParsedWithPathToPayload},
	{"headers", GetHeadersParsedWithPathToPayload},
	{"cookies", GetCookiesParsedWithPathToPayload},
}
