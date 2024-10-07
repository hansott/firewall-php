package main

import (
	"main/attack"
	"main/context"
	"main/grpc"
	"main/log"
	sql_injection "main/vulnerabilities/sql-injection"
)

func OnPreSqlQueryExecuted() string {
	query := context.GetSqlQuery()
	dialect := context.GetSqlDialect()
	operation := context.GetFunctionName()
	if query == "" || dialect == "" {
		return ""
	}
	log.Info("Got PDO query: ", query, " dialect: ", dialect)
	res := sql_injection.CheckContextForSqlInjection(query, operation, dialect)
	if res != nil {
		go grpc.OnAttackDetected(attack.GetAttackDetectedProto(*res))
		return attack.GetAttackDetectedAction(*res)
	}
	return ""
}
