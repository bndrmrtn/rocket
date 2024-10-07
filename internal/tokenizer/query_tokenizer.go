package tokenizer

import (
	"regexp"
	"strings"
)

type QueryTokenizer struct {
	tokens []BuildToken

	queries []Query
}

func NewQueryTokenizer(tokens []BuildToken) *QueryTokenizer {
	return &QueryTokenizer{
		tokens:  tokens,
		queries: []Query{},
	}
}

func (q *QueryTokenizer) Generate() error {
	for _, token := range q.tokens {
		name, args, err := parseFunction(token.Key)
		if err != nil {
			return NewErrorWithPosition(err.Error(), Token{
				Value:    token.Value,
				FileName: token.File,
				Line:     token.Line,
			})
		}

		q.queries = append(q.queries, Query{
			Name:      name,
			Arguments: q.makeArgs(args),
			Tokens:    q.tokenizeData(token.Value),
			BT:        token,
		})
	}
	return nil
}

func (q *QueryTokenizer) tokenizeData(input string) []string {
	// pattern := `\[\]|\{[^}]+\}|\(|\)|\w+|==|\|\||[=().{}]|\s+|,|[0-9]+`
	pattern := `\[\w*(?:\.\w*)*\]|\{[^}]*\}|\(|\)|\w+|.|==|!=|>|>=|->|in|not\s+in|\?\?|[0-9]+|\|\|`
	re := regexp.MustCompile(pattern)
	tokens := re.FindAllString(input, -1)

	var result []string
	var buffer []string
	openParens := 0

	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		if token == "(" {
			openParens++
		}

		if token == ")" {
			openParens--
		}

		if openParens > 0 || (openParens == 0 && len(buffer) > 0) {
			buffer = append(buffer, token)
		}

		if openParens == 0 && len(buffer) > 0 {
			result = append(result, strings.Join(buffer, " "))
			buffer = nil
		} else if openParens == 0 {
			result = append(result, token)
		}
	}

	return result
}

func (q *QueryTokenizer) makeArgs(args map[string]string) []QueryArg {
	var out []QueryArg
	for key, value := range args {
		out = append(out, QueryArg{
			Name: key,
			Type: value,
		})
	}
	return out
}

func (q *QueryTokenizer) Output() []Query {
	return q.queries
}
