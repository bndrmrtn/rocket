package tokenizer

import (
	"errors"
	"strings"
)

func parseValue(s string) (string, error) {
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return "", errors.New("empty value, unknown type")
	}

	// FIXIT: not ok
	if s[0] == '"' {
		return s[1 : len(s)-1], nil
	}

	if strings.Contains(s, " ") {
		return "", errors.New("not string type cannot have spaces")
	}

	return s, nil
}
