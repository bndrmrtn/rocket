package tokenizer

import (
	"strings"
	"unicode"
)

func parseFunctionCall(fn string) (string, []string) {
	if !strings.Contains(fn, "(") {
		return fn, []string{}
	}

	var (
		name string
		args []string
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
			args = parseFunctionCallArgs(fn[i:])
			break
		}
	}

	return name, args
}

func parseFunctionCallArgs(args string) []string {
	var (
		pos       int
		generated []string
	)

	for pos < len(args) {
		char := rune(args[pos])

		if char == '(' || unicode.IsSpace(char) || char == '\n' {
			pos++
			continue
		}

		if char == '"' {
			pos++
			var buf string
			for pos+1 < len(args) {
				if args[pos] == '"' && args[pos-1] != '\\' {
					break
				}
				buf += string(args[pos])
				pos++
			}

			if buf != "" {
				generated = append(generated, buf)
			}
			continue
		}

		if char == ',' {
			pos++
			continue
		}

		var buf string
		for pos+1 < len(args) {
			if args[pos] == ',' {
				break
			}
			buf += string(args[pos])
			pos++
		}
		if buf != "" {
			generated = append(generated, buf)
		}

		if char == ')' {
			break
		}
	}

	return generated
}
