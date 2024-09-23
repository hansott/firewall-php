package sql_injection

import "strings"

func queryContainsUserInput(query string, userInput string) bool {
	lowerQuery := strings.ToLower(query)
	lowerUserInput := strings.ToLower(userInput)

	return strings.Contains(lowerQuery, lowerUserInput)

}
