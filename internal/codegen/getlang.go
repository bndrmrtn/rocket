package codegen

import (
	"errors"

	"github.com/bndrmrtn/rocket/internal/query_interpreter"
	"github.com/bndrmrtn/rocket/internal/tokenizer"
)

type Lang interface {
	Bind(data *tokenizer.Generated, queries []query_interpreter.Query)
	Generate(func(query_interpreter.Query) string) error
	Get() string
	Save(file string) error
}

func GetLang(lang string) (Lang, error) {
	switch lang {
	case "go":
		return &Go{}, nil
	default:
		return nil, errors.New("Language not supported")
	}
}
