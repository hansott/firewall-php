package sql_injection

import "testing"

func TestQueryContainsUserInput(t *testing.T) {
	t.Run("it checks if query contains user input", func(t *testing.T) {
		if !queryContainsUserInput("SELECT * FROM 'Jonas';", "Jonas") {
			t.Errorf("queryContainsUserInput should return true")
		}
		if !queryContainsUserInput("Hi I'm MJoNaSs", "jonas") {
			t.Errorf("queryContainsUserInput should return true")
		}
		if !queryContainsUserInput("Hiya, 123^&*( is a real string", "123^&*(") {
			t.Errorf("queryContainsUserInput should return true")
		}
		if queryContainsUserInput("Roses are red", "violet") {
			t.Errorf("queryContainsUserInput should return false")
		}
	})
}
