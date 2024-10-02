package generator

import (
	"errors"

	"github.com/fatih/color"
)

func unsupported(err string) error {
	red := color.RedString("unsupported: %s", err)
	return errors.New(red)
}
