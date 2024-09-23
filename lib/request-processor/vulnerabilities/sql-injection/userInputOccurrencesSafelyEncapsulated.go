package sql_injection

import (
	"main/helpers"
	"main/utils"
	"regexp"
	"strings"
)

func escapeStringRegexp(str string) string {
	re := regexp.MustCompile(`[|\\{}()[\]^$+*?.]`)
	str = re.ReplaceAllStringFunc(str, func(match string) string {
		return "\\" + match
	})
	str = strings.ReplaceAll(str, "-", "\\x2d")
	return str
}

func getEscapeSequencesRegex() *regexp.Regexp {
	var patterns []string
	for _, seq := range SQL_ESCAPE_SEQUENCES {
		patterns = append(patterns, escapeStringRegexp(seq))
	}
	pattern := strings.Join(patterns, "|")
	return regexp.MustCompile(pattern)
}

var escapeSequencesRegex = getEscapeSequencesRegex()

func userInputOccurrencesSafelyEncapsulated(query string, userInput string) bool {

	segmentsInBetween := helpers.GetCurrentAndNextSegments(strings.Split(strings.ToLower(query), strings.ToLower(userInput)))

	for _, segment := range segmentsInBetween {
		currentSegment, nextSegment := segment["currentSegment"], segment["nextSegment"]
		input := userInput
		charBeforeUserInput := ""
		if len(currentSegment) > 0 {
			charBeforeUserInput = currentSegment[len(currentSegment)-1:]
		}

		charAfterUserInput := ""
		if len(nextSegment) > 0 {
			charAfterUserInput = nextSegment[:1]
		}

		quoteChar := ""
		find := utils.ArrayContains(SQL_STRING_CHARS, charBeforeUserInput)
		if find {
			quoteChar = charBeforeUserInput
		}
		// Special case for when the user input starts with a single quote
		// If the user input is `'value`
		// And the single quote is properly escaped with a backslash we split the following
		// `SELECT * FROM table WHERE column = '\'value'`
		// Into [`SELECT * FROM table WHERE column = '\`, `'`]
		// The char before the user input will be `\` and the char after the user input will be `'`

		for _, char := range []string{"\"", "'"} {
			if (quoteChar == "") &&
				strings.HasPrefix(input, char) &&
				len(currentSegment) > 1 &&
				currentSegment[len(currentSegment)-2:] == char+"\\" && charAfterUserInput == char {
				quoteChar = char
				charBeforeUserInput = currentSegment[len(currentSegment)-2 : len(currentSegment)-1]
				input = input[1:]
				break
			}
		}

		if quoteChar == "" {
			return false
		}

		if charBeforeUserInput != charAfterUserInput {
			return false
		}

		if strings.Contains(input, charBeforeUserInput) {
			return false
		}

		withoutEscapeSequences := escapeSequencesRegex.ReplaceAllString(input, "")
		if strings.Contains(withoutEscapeSequences, `\`) {
			return false
		}

	}

	return true
}
