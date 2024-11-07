package sql_injection

import (
	zen_internals "main/vulnerabilities/zen-internals"
)

func detectSQLInjection(query string, userInput string, dialect int) int {
	if len(userInput) <= 1 {
		// We ignore single characters since they are only able to crash the SQL Server,
		// And don't pose a big threat.
		return 0
	}

	if len(userInput) > len(query) {
		// We ignore cases where the user input is longer than the query.
		// Because the user input can't be part of the query.
		return 0
	}

	if !queryContainsUserInput(query, userInput) {
		// If the user input is not part of the query, return false (No need to check)
		return 0
	}

	if userInputOccurrencesSafelyEncapsulated(query, userInput) {
		return 0
	}

	// Executing our final check with zen_internals
	return zen_internals.DetectSQLInjection(query, userInput, dialect)

}
