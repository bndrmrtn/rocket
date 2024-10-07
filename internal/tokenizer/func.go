package tokenizer

import (
	"errors"
	"slices"
	"strings"
	"unicode"
)

func parseFunction(fn string) (string, map[string]string, error) {
	if !strings.Contains(fn, "(") {
		return fn, map[string]string{}, nil
	}

	var (
		name   string
		argMap map[string]string
	)

	for i, s := range fn {
		if unicode.IsSpace(s) {
			continue
		}

		if unicode.IsLetter(s) || s == '_' {
			name += string(s)
			continue
		}

		if len(name) > 0 && unicode.IsDigit(s) {
			name += string(s)
			continue
		}

		if s == '(' {
			var err error
			argMap, err = parseFunctionArgs(fn[i:])
			if err != nil {
				return "", nil, err
			}
			break
		}
	}

	return name, argMap, nil
}

func parseFunctionArgs(args string) (map[string]string, error) {
	var (
		pos       int
		generated = make(map[string]string)
	)

	for pos < len(args) {
		char := rune(args[pos])

		if char == '(' || unicode.IsSpace(char) || char == '\n' {
			pos++
			continue
		}

		if char == ')' {
			break
		}

		if unicode.IsLetter(char) || char == '_' {
			var buf string
			for pos < len(args) {
				if args[pos] == ',' || args[pos] == ')' {
					break
				}
				buf += string(args[pos])
				pos++
			}
			buf = strings.TrimSpace(buf)
			data := strings.SplitN(buf, " ", 2)
			if len(data) != 2 {
				return nil, errors.New("argument must have a name and a type")
			}

			if !isTypeOK(data[1]) {
				return nil, errors.New("invalid type: " + data[1])
			}

			generated[data[0]] = data[1]
			pos++
			continue
		}
	}

	return generated, nil
}

func isTypeOK(t string) bool {
	t = strings.TrimPrefix(t, "[]")

	return slices.Contains([]string{"number", "string", "byte", "bool"}, t)
}
