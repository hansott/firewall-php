package path_traversal

import (
	"main/context"
	"main/utils"
)

func CheckContextForPathTraversal(filename string, operation string, checkPathStart bool) *utils.InterceptorResult {
	for _, source := range context.SOURCES {
		mapss := source.CacheGet()
		pathString := pathToString(filename)

		for str, path := range mapss {
			str = pathToString(str)
			if detectPathTraversal(pathString, str, checkPathStart) {
				return &utils.InterceptorResult{
					Operation:     operation,
					Kind:          utils.Path_traversal,
					Source:        source.Name,
					PathToPayload: path,
					Metadata: map[string]string{
						"filename": pathString,
					},
					Payload: str,
				}
			}
		}

	}
	return nil
}

func pathToString(path string) string {
	// if starts with file:// remove it
	if len(path) > 7 && path[:7] == "file://" {
		path = path[7:]
	}
	return path
}
