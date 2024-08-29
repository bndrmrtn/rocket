package tokenizer

import (
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

type ErrPos struct {
	message string
	file    string
	line    int
	near    string
	t       string
}

func NewErrorWithPosition(m string, token Token) *ErrPos {
	return &ErrPos{
		message: m,
		file:    token.FileName,
		line:    token.Line,
		near:    strings.TrimSpace(token.Value),
		t:       "error",
	}
}

func (e *ErrPos) Error() string {
	var content string
	content += color.RedString("%s: %s\n", e.t, e.message)
	content += color.RedString("file: ") + color.GreenString("\"%s\"", e.file) + "\n"
	content += color.RedString("line: ") + color.BlueString(strconv.Itoa(e.line)) + "\n"
	if e.near != "" {
		content += color.RedString("near: ") + color.YellowString(e.near)
	}
	return content
}

func (e *ErrPos) SetType(t string) *ErrPos {
	e.t = t
	return e
}

func WarnWithPos(message string, token Token) {
	var content string
	content += color.YellowString("warning: %s\n", message)
	content += color.GreenString("\"%s\"", token.FileName) + ":" + color.BlueString(strconv.Itoa(token.Line))
	fmt.Println(content)
}
