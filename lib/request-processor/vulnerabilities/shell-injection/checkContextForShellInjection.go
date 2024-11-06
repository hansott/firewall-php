package shell_injection

import (
	"main/context"
	"main/utils"
)

func CheckContextForShellInjection(command string, operation string) *utils.InterceptorResult {
	for _, source := range context.SOURCES {
		mapss := source.CacheGet()

		for str, path := range mapss {
			if detectShellInjection(command, str) {
				return &utils.InterceptorResult{
					Operation:     operation,
					Kind:          utils.Shell_injection,
					Source:        source.Name,
					PathToPayload: path,
					Metadata: map[string]string{
						"command": command,
					},
					Payload: str,
				}
			}
		}
	}

	return nil
}
