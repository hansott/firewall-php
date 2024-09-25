package dialects

type SQLDialectSQLite struct{}

// GetDangerousStrings returns an empty list of dangerous strings for SQLite
func (SQLDialectSQLite) GetDangerousStrings() []string {
	return []string{}
}

// GetKeywords returns a list of keywords for SQLite
func (SQLDialectSQLite) GetKeywords() []string {
	return []string{
		"VACUUM",
		"ATTACH",
		"DETACH",
	}
}
