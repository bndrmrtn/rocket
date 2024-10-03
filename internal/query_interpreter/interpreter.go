package query_interpreter

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bndrmrtn/rocket/internal/tokenizer"
)

type Interpreter struct {
	data *tokenizer.Generated
}

func NewInterpreter(data *tokenizer.Generated) *Interpreter {
	return &Interpreter{
		data: data,
	}
}

func (i *Interpreter) InterpretAll() ([]Query, error) {
	var queries []Query

	for _, tQuery := range i.data.Queries {
		query, err := i.Interpret(tQuery.Name)
		if err != nil {
			return nil, err
		}
		queries = append(queries, *query)
	}

	return queries, nil
}

func (i *Interpreter) Interpret(name string) (*Query, error) {
	var tQuery tokenizer.Query
	for _, tQuery = range i.data.Queries {
		if tQuery.Name == name {
			break
		}
	}
	if tQuery.Name == "" {
		return nil, fmt.Errorf("query %s not found", name)
	}

	tokens := tQuery.Tokens

	var query Query

	for inx, token := range tokens {
		if inx == 0 && token == string(MultiResult) {
			query.MultiResult = true
			continue
		}

		if inx == 0 && query.MultiResult == false &&
			slices.Contains(QueryTypeKeywords, token) ||
			inx == 1 && slices.Contains(QueryTypeKeywords, token) {
			continue
		} else if inx == 0 && query.MultiResult == false || inx == 1 && query.MultiResult == true {
			return nil, errors.New("invalid token: " + token)
		}

		if inx == 1 && query.MultiResult == false ||
			inx == 2 && query.MultiResult == true {
			query.Fields = i.makeFields(token)
			continue
		}

		if inx == 2 && query.MultiResult == false || inx == 3 && query.MultiResult == true {
			query.From = token
			continue
		}

		if inx > 2 && query.MultiResult == false || inx > 3 && query.MultiResult == true {
			if slices.Contains(MethodKeywords, token) {
				var parantheses string
				if len(tokens) > inx+1 &&
					strings.HasPrefix(tokens[inx+1], "(") && strings.HasSuffix(tokens[inx+1], ")") {
					parantheses = tokens[inx+1]
				}

				err := i.parseQueryMethod(token, parantheses, &query)
				if err != nil {
					return nil, fmt.Errorf("error parsing method %s: %w", token, err)
				}
			}
		}

		if token == "." {
			continue
		}

	}

	return &query, nil
}

func (i *Interpreter) makeFields(token string) map[string][]string {
	token = strings.Trim(token, "{}")
	rawFields := strings.Split(token, ",")

	var fields = make(map[string][]string)

	for _, field := range rawFields {
		rawFieldData := strings.SplitN(field, ".", 2)
		if len(rawFieldData) == 1 {
			fields["#fieldFrom#"] = append(fields["#fieldFrom#"], rawFieldData[0])
		} else {
			fields[rawFieldData[0]] = append(fields[rawFieldData[0]], rawFieldData[1])
		}
	}

	return fields
}

func (i *Interpreter) parseQueryMethod(methodName string, parantheses string, query *Query) error {
	var err error

	switch methodName {
	case "Where":
		err = parseWhereFunc(parantheses, query)
	}

	return err
}
