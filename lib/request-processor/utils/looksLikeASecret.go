package utils

import (
	"regexp"
	"strings"
)

var (
	LOWERCASE             = strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	UPPERCASE             = strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	NUMBERS               = strings.Split("0123456789", "")
	SPECIAL               = strings.Split("!#$%^&*|;:<>", "")
	KNOWN_WORD_SEPARATORS = []string{"-"}
	WHITE_SPACE           = regexp.MustCompile(`\s+`)
	MINIMUM_LENGTH        = 10
)

func contains(input string, chars []string) bool {
	found := false
	for _, char := range chars {
		if strings.Contains(input, char) {
			found = true
			break
		}
	}
	return found
}

func looksLikeASecret(str string) bool {
	if len(str) <= MINIMUM_LENGTH {
		return false
	}

	if !contains(str, NUMBERS) {
		return false
	}

	// If the string doesn't have at least 2 different charsets, it's not a secret
	charsets := []bool{contains(str, LOWERCASE), contains(str, UPPERCASE), contains(str, SPECIAL)}

	numCharsets := 0
	for _, charset := range charsets {
		if charset {
			numCharsets++
		}
	}
	if numCharsets < 2 {
		return false
	}

	// If the string has white space, it's not a secret
	if WHITE_SPACE.MatchString(str) {
		return false
	}

	if contains(str, KNOWN_WORD_SEPARATORS) {
		return false
	}

	// Check uniqueness of characters in a window of 10 characters
	windowSize := MINIMUM_LENGTH
	var ratios []float64
	for i := 0; i <= len(str)-windowSize; i++ {
		window := str[i : i+windowSize]
		uniqueChars := make(map[rune]struct{})
		for _, char := range window {
			uniqueChars[char] = struct{}{}
		}
		ratios = append(ratios, float64(len(uniqueChars))/float64(windowSize))
	}

	sum := 0.0
	for _, ratio := range ratios {
		sum += ratio
	}
	averageRatio := sum / float64(len(ratios))

	return averageRatio > 0.75
}
