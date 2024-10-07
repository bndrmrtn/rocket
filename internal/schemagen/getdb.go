package schemagen

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/rocket/internal/tokenizer"
)

type DB interface {
	Bind(data *tokenizer.Generated)
	Get() string
	Create(out string) error
}

func GetDB(name string) (DB, error) {
	switch strings.ToLower(name) {
	case "mysql":
		return &mysql{}, nil
	}

	return nil, fmt.Errorf("database %s not found", name)
}
