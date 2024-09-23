package main

import (
	"main/attack"
	"main/context"
	"main/grpc"
	"main/log"
	"main/utils"
	sql_injection "main/vulnerabilities/sql-injection"
)

func OnPreExecutedPdoQuery() string {
	query := context.GetSqlQuery()
	dialect := context.GetSqlDialect()
	operation := context.GetFunctionName()
	if query == "" || dialect == "" {
		return ""
	}
	dialect_class := utils.GetSqlDialectFromString(dialect)
	log.Info("Got PDO query: ", query, " dialect: ", dialect)
	res := sql_injection.CheckContextForSqlInjection(query, operation, dialect_class)
	if res != nil {
		go grpc.OnAttackDetected(*res)
		return attack.GetAttackDetectedAction(*res)
	}
	return ""
}
