package sql_injection

import (
	"main/context"
	"main/utils"
	"main/vulnerabilities/sql-injection/dialects"
)

/**
 * This function goes over all the different input types in the context and checks
 * if it's a possible SQL Injection, if so the function returns an InterceptorResult
 */
func CheckContextForSqlInjection(sql string, operation string, dialect dialects.SQLDialect) *utils.InterceptorResult {
	for _, source := range context.SOURCES {
		mapss := source.CacheGet()

		for str, path := range mapss {
			if detectSQLInjection(sql, str, dialect) {
				return &utils.InterceptorResult{
					Operation:     operation,
					Kind:          utils.Sql_injection,
					Source:        source.Name,
					PathToPayload: path,
					Metadata:      map[string]string{},
					Payload:       str,
				}
			}
		}
	}
	return nil

}
