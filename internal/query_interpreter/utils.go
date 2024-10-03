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

func tokenizeOperation(input string) []string {
	pattern := `\w+|==|!=|>|>=|<|<=|in|not\s+in|\?\?|[0-9]+|\|\|`
	re := regexp.MustCompile(pattern)
	tokens := re.FindAllString(input, -1)

	var result []string
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token != "" {
			result = append(result, token)
		}
	}

	return result
}
