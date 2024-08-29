package query_interpreter

import "github.com/bndrmrtn/rocket/internal/tokenizer"

type Interpreter struct {
	data  *tokenizer.Generated
	query tokenizer.Query
}

func New(data *tokenizer.Generated, name string) *Interpreter {
	var query tokenizer.Query
	for _, t := range data.Queries {
		if t.Name == name {
			query = t
			break
		}
	}
	if query.Name == "" {
		panic("Invalid query: " + name)
	}

	return &Interpreter{
		data:  data,
		query: query,
	}
}

func (*Interpreter) Interpret() (Result, error) {
	return Result{}, nil
}
