package sql_injection

import (
	"main/vulnerabilities/sql-injection/dialects"
	"testing"
)

func TestDetectSQLInjectionSQLite(t *testing.T) {
	isSqlInjection := func(t *testing.T, sql string, input string) {
		if !detectSQLInjection(sql, input, dialects.SQLDialectSQLite{}) {
			t.Errorf("Expected SQL injection for input %q in query %q", input, sql)
		}
	}

	isNotSQLInjection := func(t *testing.T, sql string, input string) {
		if detectSQLInjection(sql, input, dialects.SQLDialectSQLite{}) {
			t.Errorf("Did not expect SQL injection for input %q in query %q", input, sql)
		}
	}

	t.Run("It flags the VACUUM command as SQL injection", func(t *testing.T) {
		isSqlInjection(t, "VACUUM;", "VACUUM")
	})

	t.Run("It flags the ATTACH command as SQL injection", func(t *testing.T) {
		isSqlInjection(t, "ATTACH DATABASE 'test.db' AS test;", "'test.db' AS test")
	})

	t.Run("It ignores postgres dollar signs", func(t *testing.T) {
		isNotSQLInjection(t, "SELECT $$", "$$")
		isNotSQLInjection(t, "SELECT $$text$$", "$$text$$")
		isNotSQLInjection(t, "SELECT $tag$text$tag$", "$tag$text$tag$")
	})
}
