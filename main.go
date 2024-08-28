package main

import (
	"fmt"
	"github.com/bndrmrtn/rocket/internal/tokenizer"
	"log"
)

func main() {
	t, err := tokenizer.New("./code/user.rocketdb")
	if err != nil {
		panic(err)
	}
	err = t.Tokenize()
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range t.GetTokens() {
		fmt.Printf("Token: %v\n", t)
	}
}
