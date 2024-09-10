package context

type Source struct {
	Name     string
	CacheGet func() map[string]string
}

var SOURCES = []Source{
	{"body", GetParsedBody},
	{"query", GetParsedQuery},
	{"headers", GetParsedHeaders},
	{"cookies", GetParsedCookies},
}
