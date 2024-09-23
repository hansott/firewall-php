package sql_injection

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"main/vulnerabilities/sql-injection/dialects"
	"testing"
)

func TestUserInputContainsSQLSyntax(t *testing.T) {

	t.Run("it flags dialect specific keywords", func(t *testing.T) {
		if !userInputContainsSQLSyntax("@@GLOBAL", dialects.SQLDialectMySQL{}) {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("it does not flag common SQL keywords", func(t *testing.T) {
		if userInputContainsSQLSyntax("SELECT", dialects.SQLDialectMySQL{}) {
			t.Errorf("Expected false, got true")
		}
	})

	t.Run("it does not flag common SQL operators", func(t *testing.T) {
		files := []string{
			filepath.Join("payloads", "Auth_Bypass.txt"),
			filepath.Join("payloads", "postgres.txt"),
			filepath.Join("payloads", "mysql.txt"),
			filepath.Join("payloads", "mssql_and_db2.txt"),
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

				if !userInputContainsSQLSyntax(sql, dialects.SQLDialectMySQL{}) {
					t.Errorf("Expected true, got false for [%s] from file %s", sql, file)
				}
				if !userInputContainsSQLSyntax(sql, dialects.SQLDialectPostgres{}) {
					t.Errorf("Expected true, got false for [%s] from file %s", sql, file)
				}
			}
		}
	})

}
