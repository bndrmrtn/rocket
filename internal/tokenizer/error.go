package tokenizer

import (
	"fmt"
)

type ErrPos struct {
	m string
}

func (e *ErrPos) Error() string {
	return e.m
}

func NewErrorWithPosition(m string, token Token) *ErrPos {
	return &ErrPos{
		m: fmt.Sprintf("Error: %s\n\tfile: %s\n\tline: %d\n\tnear: %s", m, token.FileName, token.Line, token.Value),
	}
}
