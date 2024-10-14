package query_interpreter

import (
	"regexp"
	"strings"
)

func explodeMultiOperations(input string) []string {
	pattern := `\(([^()]+)\)|\|\||&&`
	re := regexp.MustCompile(pattern)
	tokens := re.FindAllStringSubmatch(input, -1)

	var result []string
	for _, token := range tokens {
		if token[1] != "" {
			result = append(result, strings.TrimSpace(token[1]))
		} else {
			result = append(result, strings.TrimSpace(token[0]))
		}
	}

	return result
}

// FIXIT: Bad regex
func tokenizeOperation(input string) []string {
	cleanedInput := regexp.MustCompile(`\s*\.\s*`).ReplaceAllString(input, ".")
	cleanedInput = regexp.MustCompile(`\s*:\s*`).ReplaceAllString(cleanedInput, ":")

	re := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_.]*|==|!=|>=|<=|>|<|&&|\|\||:\w+|\d+|[^\s]+)`)
	tokens := re.FindAllString(cleanedInput, -1)

	return tokens
}
