package sql_injection

import (
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("SQL_KEYWORDS are not empty", func(t *testing.T) {
		for _, keyword := range SQL_KEYWORDS {
			if keyword == "" {
				t.Errorf("SQL_KEYWORDS should not be empty")
			}
		}
	})

	t.Run("QL_KEYWORDS are uppercase", func(t *testing.T) {
		for _, keyword := range SQL_KEYWORDS {
			if keyword != strings.ToUpper(keyword) {
				t.Errorf("SQL_KEYWORDS should be uppercase")
			}
		}
	})

	t.Run("COMMON_SQL_KEYWORDS are not empty", func(t *testing.T) {
		for _, keyword := range COMMON_SQL_KEYWORDS {
			if keyword == "" {
				t.Errorf("COMMON_SQL_KEYWORDS should not be empty")
			}
		}
	})

	t.Run("COMMON_SQL_KEYWORDS are uppercase", func(t *testing.T) {
		for _, keyword := range COMMON_SQL_KEYWORDS {
			if keyword != strings.ToUpper(keyword) {
				t.Errorf("COMMON_SQL_KEYWORDS should be uppercase")
			}
		}
	})

	t.Run("SQL_OPERATORS are not empty", func(t *testing.T) {
		for _, operator := range SQL_OPERATORS {
			if operator == "" {
				t.Errorf("SQL_OPERATORS should not be empty")
			}
		}
	})

	t.Run("SQL_STRING_CHARS are single chars", func(t *testing.T) {
		for _, char := range SQL_STRING_CHARS {
			if len(char) != 1 {
				t.Errorf("SQL_STRING_CHARS should be single chars")
			}
		}
	})

	t.Run("SQL_DANGEROUS_IN_STRING are not empty", func(t *testing.T) {
		for _, char := range SQL_DANGEROUS_IN_STRING {
			if char == "" {
				t.Errorf("SQL_DANGEROUS_IN_STRING should not be empty")
			}
		}
	})

	t.Run("SQL_ESCAPE_SEQUENCES are not empty", func(t *testing.T) {
		for _, sequence := range SQL_ESCAPE_SEQUENCES {
			if sequence == "" {
				t.Errorf("SQL_ESCAPE_SEQUENCES should not be empty")
			}
		}
	})

}
