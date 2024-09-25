package dialects

type SQLDialectPostgres struct{}

// GetDangerousStrings returns a list of dangerous strings for Postgres
func (SQLDialectPostgres) GetDangerousStrings() []string {
	return []string{
		// https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-DOLLAR-QUOTING
		"$",
	}
}

// GetKeywords returns a list of keywords for Postgres
func (SQLDialectPostgres) GetKeywords() []string {

	return []string{
		// https://www.postgresql.org/docs/current/sql-set.html
		"CLIENT_ENCODING",
	}
}
