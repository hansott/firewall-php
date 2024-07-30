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

func looksLikeASecret(str string) bool {
	if len(str) <= MINIMUM_LENGTH {
		return false
	}

	hasNumber := false
	for _, char := range NUMBERS {
		if strings.Contains(str, char) {
			hasNumber = true
			break
		}
	}

	if !hasNumber {
		return false
	}

	hasLower := false
	hasUpper := false
	hasSpecial := false

	for _, char := range LOWERCASE {
		if strings.Contains(str, char) {
			hasLower = true
			break
		}
	}

	for _, char := range UPPERCASE {
		if strings.Contains(str, char) {
			hasUpper = true
			break
		}
	}

	for _, char := range SPECIAL {
		if strings.Contains(str, char) {
			hasSpecial = true
			break
		}
	}

	charsets := []bool{hasLower, hasUpper, hasSpecial}

	// If the string doesn't have at least 2 different charsets, it's not a secret
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

	for _, separator := range KNOWN_WORD_SEPARATORS {
		if strings.Contains(str, separator) {
			return false
		}
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
