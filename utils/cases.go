package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func PascalCase(input string) string {
	var result strings.Builder

	input = strings.ReplaceAll(input, "_", " ")
	input = strings.ReplaceAll(input, "-", " ")

	capitalizeNext := true
	for i, ch := range input {
		if i > 0 && unicode.IsUpper(ch) && unicode.IsLower(rune(input[i-1])) {
			capitalizeNext = true
		}

		if capitalizeNext {
			result.WriteRune(unicode.ToUpper(ch))
			capitalizeNext = false
		} else {
			result.WriteRune(unicode.ToLower(ch))
		}

		if unicode.IsSpace(ch) {
			capitalizeNext = true
		}
	}

	return strings.ReplaceAll(strings.ReplaceAll(result.String(), " ", ""), "Id", "ID")
}

func SnakeCase(input string) string {
	input = PascalCase(input)
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(input, "${1}_${2}")
	return strings.ToLower(snake)
}
