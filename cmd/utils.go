package cmd

import (
	"fmt"

	"github.com/fatih/color"
)

func success(s ...interface{}) {
	c := color.New(color.FgHiGreen)

	fmt.Println("🆗", c.Sprint(s...))
}
