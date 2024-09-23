package sql_injection

import (
	"main/vulnerabilities/sql-injection/dialects"
	"testing"
)

func TestDetectSQLInjectionMySQL(t *testing.T) {

	isSqlInjection := func(t *testing.T, sql string, input string) {
		if !detectSQLInjection(sql, input, dialects.SQLDialectMySQL{}) {
			t.Errorf("Expected SQL injection for input %q in query %q", input, sql)
		}
	}

	isNotSQLInjection := func(t *testing.T, sql string, input string) {
		if detectSQLInjection(sql, input, dialects.SQLDialectMySQL{}) {
			t.Errorf("Did not expect SQL injection for input %q in query %q", input, sql)
		}
	}
	t.Run("It flags MySQL bitwise operator as SQL injection", func(t *testing.T) {
		isSqlInjection(t, "SELECT 10 ^ 12", "10 ^ 12")
	})

	t.Run("It ignores postgres dollar signs", func(t *testing.T) {
		isNotSQLInjection(t, "SELECT $$", "$$")
		isNotSQLInjection(t, "SELECT $$text$$", "$$text$$")
		isNotSQLInjection(t, "SELECT $tag$text$tag$", "$tag$text$tag$")
	})

	t.Run("It flags SET GLOBAL as SQL injection", func(t *testing.T) {
		isSqlInjection(t, "SET GLOBAL max_connections = 1000", "GLOBAL max_connections")
		isSqlInjection(t, "SET @@GLOBAL.max_connections = 1000", "@@GLOBAL.max_connections = 1000")
		isSqlInjection(t, "SET @@GLOBAL.max_connections=1000", "@@GLOBAL.max_connections=1000")

		isNotSQLInjection(t, "SELECT * FROM users WHERE id = 'SET GLOBAL max_connections = 1000'", "SET GLOBAL max_connections = 1000")
		isNotSQLInjection(t, "SELECT * FROM users WHERE id = 'SET @@GLOBAL.max_connections = 1000'", "SET @@GLOBAL.max_connections = 1000")
	})

	t.Run("It flags SET SESSION as SQL injection", func(t *testing.T) {
		isSqlInjection(t, "SET SESSION max_connections = 1000", "SESSION max_connections")
		isSqlInjection(t, "SET @@SESSION.max_connections = 1000", "@@SESSION.max_connections = 1000")
		isSqlInjection(t, "SET @@SESSION.max_connections=1000", "@@SESSION.max_connections=1000")

		isNotSQLInjection(t, "SELECT * FROM users WHERE id = 'SET SESSION max_connections = 1000'", "SET SESSION max_connections = 1000")
		isNotSQLInjection(t, "SELECT * FROM users WHERE id = 'SET @@SESSION.max_connections = 1000'", "SET @@SESSION.max_connections = 1000")
	})

	t.Run("It flags SET CHARACTER SET as SQL injection", func(t *testing.T) {
		isSqlInjection(t, "SET CHARACTER SET utf8", "CHARACTER SET utf8")
		isSqlInjection(t, "SET CHARACTER SET=utf8", "CHARACTER SET=utf8")
		isSqlInjection(t, "SET CHARSET utf8", "CHARSET utf8")
		isSqlInjection(t, "SET CHARSET=utf8", "CHARSET=utf8")

		isNotSQLInjection(t, "SELECT * FROM users WHERE id = 'SET CHARACTER SET utf8'", "SET CHARACTER SET utf8")
		isNotSQLInjection(t, "SELECT * FROM users WHERE id = 'SET CHARACTER SET=utf8'", "SET CHARACTER SET=utf8")
		isNotSQLInjection(t, "SELECT * FROM users WHERE id = 'SET CHARSET utf8'", "SET CHARSET utf8")
		isNotSQLInjection(t, "SELECT * FROM users WHERE id = 'SET CHARSET=utf8'", "SET CHARSET=utf8")
	})

}
