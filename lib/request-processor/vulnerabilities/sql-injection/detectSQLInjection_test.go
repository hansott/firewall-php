package sql_injection

import (
	"bufio"
	"fmt"
	"main/vulnerabilities/sql-injection/dialects"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var BAD_SQL_COMMANDS = []string{
	"Roses are red insErt are blue",
	"Roses are red cREATE are blue",
	"Roses are red drop are blue",
	"Roses are red updatE are blue",
	"Roses are red SELECT are blue",
	"Roses are red dataBASE are blue",
	"Roses are red alter are blue",
	"Roses are red grant are blue",
	"Roses are red savepoint are blue",
	"Roses are red commit are blue",
	"Roses are red or blue",
	"Roses are red and lovely",
	"This is a group_concat_test",
	// Test some special characters
	"I'm writting you",
	"Termin;ate",
	"Roses <> violets",
	"Roses < Violets",
	"Roses > Violets",
	"Roses != Violets",
}

var GOOD_SQL_COMMANDS = []string{
	"Roses are red rollbacks are blue",
	"Roses are red truncates are blue",
	"Roses are reddelete are blue",
	"Roses are red WHEREis blue",
	"Roses are red ORis isAND",
	// Check for some general statements
	`abcdefghijklmnop@hotmail.com`,
	// Test some special characters
	"steve@yahoo.com",
	// Test SQL Function (that should not be blocked)
	"I was benchmark ing",
	"We were delay ed",
	"I will waitfor you",
	// Allow single characters
	"#",
	"'",
}

var IS_NOT_INJECTION = [][]string{
	{`'UNION 123' UNION "UNION 123"`, "UNION 123"}, // String encapsulation
	{`'union'  is not "UNION"`, "UNION!"},          // String not present in SQL
	{`"UNION;"`, "UNION;"},                         // String encapsulation
	{"SELECT * FROM table", "*"},
	{`"COPY/*"`, "COPY/*"},                   // String encapsulated but dangerous chars
	{`'union'  is not "UNION--"`, "UNION--"}, // String encapsulated but dangerous chars
}

var IS_INJECTION = [][]string{
	{`'union'  is not UNION`, "UNION"}, // String not always encapsulated
	{`UNTER;`, "UNTER;"},               // String not encapsulated and dangerous char (;)
}

func TestDetectSQLInjection(t *testing.T) {
	isSqlInjection := func(t *testing.T, sql string, input string) {
		if !detectSQLInjection(sql, input, dialects.SQLDialectSQLite{}) {
			t.Errorf("[SQLite] Expected SQL injection for input %q in query %q", input, sql)
		}
		if !detectSQLInjection(sql, input, dialects.SQLDialectPostgres{}) {
			t.Errorf("[Postgres] Expected SQL injection for input %q in query %q", input, sql)
		}
		if !detectSQLInjection(sql, input, dialects.SQLDialectMySQL{}) {
			t.Errorf("[MySQL] Expected SQL injection for input %q in query %q", input, sql)
		}

	}

	isNotSqlInjection := func(t *testing.T, sql string, input string) {
		if detectSQLInjection(sql, input, dialects.SQLDialectSQLite{}) {
			t.Errorf("Did not expect SQL injection for input %q in query %q", input, sql)
		}
		if detectSQLInjection(sql, input, dialects.SQLDialectPostgres{}) {
			t.Errorf("Did not expect SQL injection for input %q in query %q", input, sql)
		}
		if detectSQLInjection(sql, input, dialects.SQLDialectMySQL{}) {
			t.Errorf("Did not expect SQL injection for input %q in query %q", input, sql)
		}
	}

	t.Run("Test the detectSQLInjection() function", func(t *testing.T) {
		for _, sql := range BAD_SQL_COMMANDS {
			isSqlInjection(t, sql, sql)
		}

		for _, sql := range GOOD_SQL_COMMANDS {
			isNotSqlInjection(t, sql, sql)
		}
	})

	t.Run("Test detectSQLInjection() function", func(t *testing.T) {
		for _, test := range IS_NOT_INJECTION {
			isNotSqlInjection(t, test[0], test[1])
		}

		for _, test := range IS_INJECTION {
			isSqlInjection(t, test[0], test[1])
		}
	})

	t.Run("It allows escape sequences", func(t *testing.T) {
		isSqlInjection(t, "SELECT * FROM users WHERE id = 'users\\'", "users\\")
		isSqlInjection(t, "SELECT * FROM users WHERE id = 'users\\\\'", "users\\\\")

		isNotSqlInjection(t, "SELECT * FROM users WHERE id = '\nusers'", "\nusers")
		isNotSqlInjection(t, "SELECT * FROM users WHERE id = '\rusers'", "\rusers")
		isNotSqlInjection(t, "SELECT * FROM users WHERE id = '\tusers'", "\tusers")
	})

	t.Run("user input inside IN (...)", func(t *testing.T) {
		isSqlInjection(t, "SELECT * FROM users WHERE id IN ('123')", "'123'")
		isNotSqlInjection(t, "SELECT * FROM users WHERE id IN (123)", "123")
		isNotSqlInjection(t, "SELECT * FROM users WHERE id IN (123, 456)", "123")
		isNotSqlInjection(t, "SELECT * FROM users WHERE id IN (123, 456)", "456")
		isNotSqlInjection(t, "SELECT * FROM users WHERE id IN ('123')", "123")
		isNotSqlInjection(t, "SELECT * FROM users WHERE id IN (13,14,15)", "13,14,15")
		isNotSqlInjection(t, "SELECT * FROM users WHERE id IN (13, 14, 154)", "13, 14, 154")
		isSqlInjection(t, "SELECT * FROM users WHERE id IN (13, 14, 154) OR (1=1)", "13, 14, 154) OR (1=1")
	})

	t.Run("It checks whether the string is safely escaped", func(t *testing.T) {
		isSqlInjection(t, `SELECT * FROM comments WHERE comment = 'I'm writting you'`, "I'm writting you")
		isSqlInjection(t, `SELECT * FROM comments WHERE comment = "I"m writting you"`, `I"m writting you`)
		isSqlInjection(t, "SELECT * FROM `comm`ents`", "`comm`ents")

		isNotSqlInjection(t, `SELECT * FROM comments WHERE comment = "I'm writting you"`, "I'm writting you")
		isNotSqlInjection(t, `SELECT * FROM comments WHERE comment = 'I"m writting you'`, `I"m writting you`)
		isNotSqlInjection(t, "SELECT * FROM comments WHERE comment = \"I`m writting you\"", "I`m writting you")
		isNotSqlInjection(t, "SELECT * FROM `comm'ents`", "comm'ents")

	})

	t.Run("it does not flag queries starting with SELECT and having select in user input", func(t *testing.T) {
		isNotSqlInjection(t, "SELECT * FROM users WHERE id = 1", "SELECT")
	})

	t.Run("It does not flag escaped # as SQL injection", func(t *testing.T) {
		isNotSqlInjection(t, "SELECT * FROM hashtags WHERE name = '#hashtag'", "#hashtag")
	})

	t.Run("Comment is same as user input", func(t *testing.T) {
		isSqlInjection(t, "SELECT * FROM hashtags WHERE name = '-- Query by name' -- Query by name", "-- Query by name")
	})

	t.Run("input occurs in comment", func(t *testing.T) {
		isNotSqlInjection(t, "SELECT * FROM hashtags WHERE name = 'name' -- Query by name", "name")
	})

	t.Run("User input is multiline", func(t *testing.T) {
		isSqlInjection(t, `SELECT * FROM users WHERE id = 'a'
OR 1=1#'`, `a'
OR 1=1#`)
		isNotSqlInjection(t, `SELECT * FROM users WHERE id = 'a
b
c';`, `a
b
c`)
	})

	t.Run("user input is longer than query", func(t *testing.T) {
		isNotSqlInjection(t, `SELECT * FROM users`, `SELECT * FROM users WHERE id = 'a'`)
	})

	t.Run("It flags multiline queries correctly", func(t *testing.T) {
		isSqlInjection(
			t,
			`
			  SELECT * FROM `+"`users`"+`
			  WHERE id = 123
			`,
			"users`",
		)
		isSqlInjection(t, `
			SELECT *
			FROM users
			WHERE id = '1' OR 1=1
				AND is_escaped = '1'' OR 1=1'
		`,
			"1' OR 1=1",
		)
		isSqlInjection(t, `
			SELECT *
			FROM users
			WHERE id = '1' OR 1=1
				AND is_escaped = "1' OR 1=1"
		`,
			"1' OR 1=1",
		)
		isNotSqlInjection(t, `
			SELECT * FROM `+"`users`"+`
			WHERE id = 123
		`,
			"123",
		)
		isNotSqlInjection(t, `
			SELECT * FROM `+"`us``ers`"+`
			WHERE id = 123
		`,
			"users",
		)
		isNotSqlInjection(t, `
			SELECT * FROM users
			WHERE id = 123
		`,
			"123",
		)

		isNotSqlInjection(t, `
			SELECT * FROM users
			WHERE id = '123'
		`,
			"123",
		)
		isNotSqlInjection(t, `
			SELECT *
			FROM users
			WHERE is_escaped = "1' OR 1=1"
		`,
			"1' OR 1=1",
		)
	})

	for _, dangerous := range SQL_DANGEROUS_IN_STRING {
		t.Run("It flags dangerous string "+dangerous+" as SQL injection", func(t *testing.T) {
			input := dangerous + "a"
			isSqlInjection(t, "SELECT * FROM users WHERE "+input, input)
		})
	}

	t.Run("It flags function calls as SQL injections", func(t *testing.T) {
		isSqlInjection(t, "foobar()", "foobar()")
		isSqlInjection(t, "foobar(1234567)", "foobar(1234567)")
		isSqlInjection(t, "foobar       ()", "foobar       ()")
		isSqlInjection(t, ".foobar()", ".foobar()")
		isSqlInjection(t, "20+foobar()", "20+foobar()")
		isSqlInjection(t, "20-foobar(", "20-foobar(")
		isSqlInjection(t, "20<foobar()", "20<foobar()")
		isSqlInjection(t, "20*foobar  ()", "20*foobar  ()")
		isSqlInjection(t, "!foobar()", "!foobar()")
		isSqlInjection(t, "=foobar()", "=foobar()")
		isSqlInjection(t, "1foobar()", "1foobar()")
		isSqlInjection(t, "1foo_bar()", "1foo_bar()")
		isSqlInjection(t, "1foo-bar()", "1foo-bar()")
		isSqlInjection(t, "#foobar()", "#foobar()")

		isNotSqlInjection(t, "foobar)", "foobar)")
		isNotSqlInjection(t, "foobar      )", "foobar      )")
		isNotSqlInjection(t, "€foobar()", "€foobar()")
	})

	t.Run("It flags lowercased input as SQL injection", func(t *testing.T) {
		isSqlInjection(t,
			`
      SELECT id,
               email,
               password_hash,
               registered_at,
               is_confirmed,
               first_name,
               last_name
        FROM users WHERE email_lowercase = '' or 1=1 -- a',
    `,
			"' OR 1=1 -- a",
		)
	})

	// Taken from https://github.com/payloadbox/sql-injection-payload-list/tree/master
	files := []string{
		filepath.Join("payloads", "Auth_Bypass.txt"),
		filepath.Join("payloads", "postgres.txt"),
		filepath.Join("payloads", "mysql.txt"),
		filepath.Join("payloads", "mssql_and_db2.txt"),
	}

	escapeLikeDatabase := func(str string, char string) string {
		// Escape the specified character and backslash

		// Replace occurrences of char and \ with \char and \\
		replacer := strings.NewReplacer(
			char, "\\"+char, // escape char
			"\\", "\\\\", // double escape backslash
		)
		return char + replacer.Replace(str) + char
	}

	testSqlPayload := func(sql string) {
		t.Run(fmt.Sprintf("It flags %q as SQL injection", sql), func(t *testing.T) {
			isSqlInjection(t, sql, sql)
		})

		t.Run(fmt.Sprintf("It flags %q as SQL injection (in query)", sql), func(t *testing.T) {
			isSqlInjection(t, "SELECT * FROM users WHERE id = "+sql, sql)
		})

		t.Run(fmt.Sprintf("It does not flag %q as SQL injection (escaped with single quotes using backslash)", sql), func(t *testing.T) {
			escaped := escapeLikeDatabase(sql, `'`)
			fmt.Println("Before: " + sql + " After: " + escaped)
			isNotSqlInjection(t, "SELECT * FROM users WHERE id = "+escaped, sql)
		})

		t.Run(fmt.Sprintf("It does not flag %q as SQL injection (escaped with double quotes using backslash)", sql), func(t *testing.T) {
			escaped := escapeLikeDatabase(sql, `"`)
			isNotSqlInjection(t, "SELECT * FROM users WHERE id = "+escaped, sql)
		})

	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", file, err)
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			sql := scanner.Text()
			testSqlPayload(sql)
		}
	}

}
