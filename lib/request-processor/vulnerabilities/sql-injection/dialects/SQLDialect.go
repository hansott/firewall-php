package dialects

type SQLDialect interface {
	// GetDangerousStrings returns a list of dangerous strings specific to the SQL dialect.
	// These are matched without surrounding spaces, so if you add "SELECT" it will match "SELECT" and "SELECTED".
	GetDangerousStrings() []string

	// GetKeywords returns a list of keywords specific to the SQL dialect.
	// These are matched with surrounding spaces, so if you add "SELECT" it will match "SELECT" but not "SELECTED".
	GetKeywords() []string
}
