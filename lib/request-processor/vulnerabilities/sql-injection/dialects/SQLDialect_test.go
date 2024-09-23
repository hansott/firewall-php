package dialects

import "testing"

var dialects = []SQLDialect{&SQLDialectMySQL{}, &SQLDialectPostgres{}, &SQLDialectSQLite{}}

func TestGetDangerousStrings(t *testing.T) {
	for _, dialect := range dialects {
		keywords := dialect.GetDangerousStrings()
		// check if the keywords len is equal to len of the set of keywords
		dangerousStrings := map[string]bool{}
		for _, keyword := range keywords {
			dangerousStrings[keyword] = true
		}
		if len(dangerousStrings) != len(keywords) {
			t.Errorf("Duplicate dangerous strings found in %T", dialect)
		}

	}
}

func TestGetKeywords(t *testing.T) {
	for _, dialect := range dialects {
		keywords := dialect.GetKeywords()
		for _, keyword := range keywords {
			if keyword == "" {
				t.Errorf("Empty keyword found in %T", dialect)
			}
		}
	}
}
