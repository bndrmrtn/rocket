package schemagen

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/rocket/internal/query_interpreter"
	"github.com/bndrmrtn/rocket/internal/tokenizer"
)

type DB interface {
	Bind(data *tokenizer.Generated)
	Get() (string, error)
	Create(out string) error
	GetQueryParser() func(query_interpreter.Query) string
}

func GetDB(name string) (DB, error) {
	switch strings.ToLower(name) {
	case "mysql":
		return &mysql{}, nil
	}

	return nil, fmt.Errorf("database %s not found", name)
}
