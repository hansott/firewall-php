package sql_injection

import "testing"

func TestUserInputOccurrencesSafelyEncapsulated(t *testing.T) {

	t.Run("Test the userInputOccurrencesSafelyEncapsulated() function", func(t *testing.T) {
		if !userInputOccurrencesSafelyEncapsulated(` Hello Hello 'UNION'and also "UNION" `, "UNION") {
			t.Errorf("Expected true, got false")
		}
		if !userInputOccurrencesSafelyEncapsulated(`"UNION"`, "UNION") {
			t.Errorf("Expected true, got false")
		}
		if !userInputOccurrencesSafelyEncapsulated("`UNION`", "UNION") {
			t.Errorf("Expected true, got false")
		}
		if userInputOccurrencesSafelyEncapsulated("`U`NION`", "U`NION") {
			t.Errorf("Expected false, got true")
		}
		if !userInputOccurrencesSafelyEncapsulated(` 'UNION' `, "UNION") {
			t.Errorf("Expected true, got false")
		}
		if !userInputOccurrencesSafelyEncapsulated(`"UNION"'UNION'`, "UNION") {
			t.Errorf("Expected true, got false")
		}
		if userInputOccurrencesSafelyEncapsulated(`'UNION'"UNION"UNION`, "UNION") {
			t.Errorf("Expected false, got true")
		}
		if userInputOccurrencesSafelyEncapsulated(`'UNION'UNION"UNION"`, "UNION") {
			t.Errorf("Expected false, got true")
		}
		if userInputOccurrencesSafelyEncapsulated("UNION", "UNION") {
			t.Errorf("Expected false, got true")
		}
		if !userInputOccurrencesSafelyEncapsulated(`"UN'ION"`, "UN'ION") {
			t.Errorf("Expected true, got false")
		}
		if !userInputOccurrencesSafelyEncapsulated(`'UN"ION'`, `UN"ION`) {
			t.Errorf("Expected true, got false")
		}
		if !userInputOccurrencesSafelyEncapsulated(`SELECT * FROM cats WHERE id = 'UN"ION' AND id = "UN'ION"`, `UN"ION`) {
			t.Errorf("Expected true, got false")
		}
		if userInputOccurrencesSafelyEncapsulated(`SELECT * FROM cats WHERE id = 'UN'ION' AND id = "UN'ION"`, `UN'ION`) {
			t.Errorf("Expected false, got true")
		}
		if userInputOccurrencesSafelyEncapsulated(`SELECT * FROM cats WHERE id = 'UNION\'`, "UNION\\") {
			t.Errorf("Expected false, got true")
		}
		if userInputOccurrencesSafelyEncapsulated(`SELECT * FROM cats WHERE id = 'UNION\\'`, "UNION\\\\") {
			t.Errorf("Expected false, got true")
		}
		if userInputOccurrencesSafelyEncapsulated(`SELECT * FROM cats WHERE id = 'UNION\\\'`, "UNION\\\\\\") {
			t.Errorf("Expected false, got true")
		}
		if !userInputOccurrencesSafelyEncapsulated(`SELECT * FROM cats WHERE id = 'UNION\n'`, `UNION\n`) {
			t.Errorf("Expected true, got false")
		}
		if userInputOccurrencesSafelyEncapsulated(`SELECT * FROM users WHERE id = '\'hello'`, "'hello'") {
			t.Errorf("Expected false, got true")
		}
		if userInputOccurrencesSafelyEncapsulated(`SELECT * FROM users WHERE id = "\"hello"`, `"hello"`) {
			t.Errorf("Expected false, got true")
		}
	})

	t.Run("surrounded with single quotes", func(t *testing.T) {
		if !userInputOccurrencesSafelyEncapsulated(`SELECT * FROM users WHERE id = '\'hello\''`, "'hello'") {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("surrounded with double quotes", func(t *testing.T) {
		if !userInputOccurrencesSafelyEncapsulated("SELECT * FROM users WHERE id = \"\\\"hello\\\"\"", `"hello"`) {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("starts with single quote", func(t *testing.T) {
		if !userInputOccurrencesSafelyEncapsulated(`SELECT * FROM users WHERE id = '\' or true--'`, "' or true--") {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("starts with double quote", func(t *testing.T) {
		if !userInputOccurrencesSafelyEncapsulated(`SELECT * FROM users WHERE id = "\" or true--"`, `" or true--`) {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("starts with single quote without SQL syntax", func(t *testing.T) {
		if !userInputOccurrencesSafelyEncapsulated(`SELECT * FROM users WHERE id = '\' hello world'`, "' hello world") {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("starts with double quote without SQL syntax", func(t *testing.T) {
		if !userInputOccurrencesSafelyEncapsulated(`SELECT * FROM users WHERE id = "\" hello world"`, `" hello world`) {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("starts with single quote (multiple occurrences)", func(t *testing.T) {
		if !userInputOccurrencesSafelyEncapsulated(`SELECT * FROM users WHERE id = '\'hello' AND id = '\'hello'`, "'hello") {
			t.Errorf("Expected true, got false")
		}
		if userInputOccurrencesSafelyEncapsulated(`SELECT * FROM users WHERE id = 'hello' AND id = '\'hello'`, "'hello") {
			t.Errorf("Expected false, got true")
		}
	})

	t.Run("starts with double quote (multiple occurrences)", func(t *testing.T) {
		if !userInputOccurrencesSafelyEncapsulated(`SELECT * FROM users WHERE id = "\"hello" AND id = "\"hello"`, `"hello`) {
			t.Errorf("Expected true, got false")
		}
		if userInputOccurrencesSafelyEncapsulated(`SELECT * FROM users WHERE id = "hello" AND id = "\"hello"`, `"hello`) {
			t.Errorf("Expected false, got true")
		}
	})

	t.Run("single quotes escaped with single quotes", func(t *testing.T) {
		if userInputOccurrencesSafelyEncapsulated(`SELECT * FROM users WHERE id = '''&'''`, "'&'") {
			t.Errorf("Expected false, got true")
		}
	})

}
