package sql_injection

import (
	"fmt"
	"main/utils"
	"main/vulnerabilities/sql-injection/dialects"
	"reflect"
	"regexp"
	"strings"
)

var cachedRegexes = make(map[string]*regexp.Regexp)

func buildRegex(dialect dialects.SQLDialect) *regexp.Regexp {

	matchSqlKeywords := "(?:[^a-z]|^)(" + strings.Join( // Lookbehind : if the keywords are preceded by one or more letters, it should not match
		append(SQL_KEYWORDS, dialect.GetKeywords()...), "|", // Look for SQL Keywords
	) + ")(?:[^a-z]|$)" // Lookahead : if the keywords are followed by one or more letters, it should not match

	// Build matchSqlOperators pattern

	sqlOperations := []string{}
	for _, operator := range SQL_OPERATORS {
		sqlOperations = append(sqlOperations, escapeStringRegexp(operator))
	}
	matchSqlOperators := "(" + strings.Join(
		sqlOperations, "|",
	) + ")"

	matchSqlFunctions := "(?:([\\s|.|" + strings.Join( // Lookbehind : A sql function should be preceded by spaces, dots,
		sqlOperations, "|", // Or sql operators
	) + "]|^)+)" +
		"([a-z0-9_-]+)" + // The name of a sql function can include letters, numbers, "_" and "-"
		"(?:[\\s]*\\()" // Lookahead : A sql function should be followed by a "(" , spaces are allowed.

	sqlDangerousSings := []string{}
	for _, dangerousString := range append(SQL_DANGEROUS_IN_STRING, dialect.GetDangerousStrings()...) {
		sqlDangerousSings = append(sqlDangerousSings, escapeStringRegexp(dangerousString))
	}

	matchDangerousStrings := strings.Join(
		sqlDangerousSings, "|",
	)

	// Combine all patterns
	pattern := fmt.Sprintf(
		"%s|%s|%s|%s",
		matchSqlKeywords,
		matchSqlOperators,
		matchSqlFunctions,
		matchDangerousStrings,
	)

	// Compile and return the regex
	return regexp.MustCompile("(?im)" + pattern)
}

/**
 * This function is the first check in order to determine if a SQL injection is happening,
 * If the user input contains the necessary characters or words for a SQL injection, this
 * function returns true.
 */
func userInputContainsSQLSyntax(userInput string, dialect dialects.SQLDialect) bool {

	// e.g. SELECT * FROM table WHERE column = 'value' LIMIT 1
	// If a query parameter is ?LIMIT=1 it would be blocked
	// If the body contains "LIMIT" or "SELECT" it would be blocked
	// These are common SQL keywords and appear in almost any SQL query

	find := utils.ArrayContains(COMMON_SQL_KEYWORDS, userInput)
	if find {
		return false
	}

	dialectName := reflect.TypeOf(dialect).Name()
	regex := cachedRegexes[dialectName]
	if regex == nil {
		regex = buildRegex(dialect)
		cachedRegexes[dialectName] = regex
	}

	return regex.MatchString(userInput)
}
