package dialects

type SQLDialectMySQL struct{}

func (SQLDialectMySQL) GetDangerousStrings() []string {
	return []string{}
}

func (SQLDialectMySQL) GetKeywords() []string {
	return []string{
		"GLOBAL",
		"SESSION",
		"PERSIST",
		"PERSIST_ONLY",
		"@@GLOBAL",
		"@@SESSION",
		"CHARACTER SET",
		"CHARSET",
	}
}
