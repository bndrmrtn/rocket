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
		})
	}
	return nil
}

func (q *QueryTokenizer) tokenizeData(input string) []string {
	pattern := `[\w]+|[=().{}]|\s+|==|_id`
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
