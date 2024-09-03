package path_traversal

import (
	"main/context"
	"main/utils"
)

type Source struct {
	Name     string
	CacheGet func() map[string]string
}

//var SOURCES = []string{"body", "query", "headers", "cookies"}

var SOURCES = []Source{
	{"body", context.GetBody},
	{"query", context.GetQuery},
	{"headers", context.GetHeaders},
	{"cookies", context.GetCookies},
}

func CheckContextForPathTraversal(filename string, operation string, checkPathStart bool) *utils.InterceptorResult {
	for _, source := range SOURCES {
		mapss := source.CacheGet()
		for str, path := range mapss {
			if detectPathTraversal(filename, str, checkPathStart) {
				return &utils.InterceptorResult{
					Operation:     operation,
					Kind:          utils.Path_traversal,
					Source:        source.Name,
					PathToPayload: path,
					Metadata: map[string]string{
						"filename": filename,
					},
					Payload: str,
				}
			}
		}

	}
	return nil
}
