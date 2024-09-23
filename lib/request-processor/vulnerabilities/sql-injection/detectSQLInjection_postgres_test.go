package sql_injection

import (
	"main/vulnerabilities/sql-injection/dialects"
	"testing"
)

func TestDetectSQLInjectionPostgres(t *testing.T) {
	isSqlInjection := func(t *testing.T, sql string, input string) {
		if !detectSQLInjection(sql, input, dialects.SQLDialectPostgres{}) {
			t.Errorf("Expected SQL injection for input %q in query %q", input, sql)
		}
	}

	isNotSQLInjection := func(t *testing.T, sql string, input string) {
		if detectSQLInjection(sql, input, dialects.SQLDialectPostgres{}) {
			t.Errorf("Did not expect SQL injection for input %q in query %q", input, sql)
		}
	}

	t.Run("It flags postgres bitwise operator as SQL injection", func(t *testing.T) {
		isSqlInjection(t, "SELECT 10 # 12", "10 # 12")
	})

	t.Run("It flags postgres type cast operator as SQL injection", func(t *testing.T) {
		isSqlInjection(t, "SELECT abc::date", "abc::date")
	})

	t.Run("It flags double dollar sign as SQL injection", func(t *testing.T) {
		isSqlInjection(t, "SELECT $$", "$$")
		isSqlInjection(t, "SELECT $$text$$", "$$text$$")
		isSqlInjection(t, "SELECT $tag$text$tag$", "$tag$text$tag$")

		isNotSQLInjection(t, "SELECT '$$text$$'", "$$text$$")
	})

	t.Run("It flags CLIENT_ENCODING as SQL injection", func(t *testing.T) {
		isSqlInjection(t, "SET CLIENT_ENCODING TO 'UTF8'", "CLIENT_ENCODING TO 'UTF8'")
		isSqlInjection(t, "SET CLIENT_ENCODING = 'UTF8'", "CLIENT_ENCODING = 'UTF8'")
		isSqlInjection(t, "SET CLIENT_ENCODING='UTF8'", "CLIENT_ENCODING='UTF8'")

		isNotSQLInjection(t, `SELECT * FROM users WHERE id = 'SET CLIENT_ENCODING = "UTF8"'`, `SET CLIENT_ENCODING = "UTF8"`)
		isNotSQLInjection(t, `SELECT * FROM users WHERE id = 'SET CLIENT_ENCODING TO "UTF8"'`, `SET CLIENT_ENCODING TO "UTF8"`)
	})
}
