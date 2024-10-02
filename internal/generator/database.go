package generator

import (
	"slices"
	"strings"
)

type Database struct{}

func GetDatabase(s string) (string, error) {
	s = strings.ToLower(s)

	if slices.Contains([]string{"mysql"}, s) {
		return s, nil
	}

	return "", unsupported("database not supported: " + s)
}
