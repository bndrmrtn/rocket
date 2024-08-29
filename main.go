package main

import (
	"encoding/json"
	"fmt"
	"github.com/bndrmrtn/rocket/internal/tokenizer"
	"log"
	"os"
)

func main() {
	t, err := tokenizer.New("./code/user.rocket")
	if err != nil {
		panic(err)
	}
	err = t.Tokenize()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tokens := t.GetTokens()

	typeT := tokenizer.NewType(tokens)
	err = typeT.Generate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data := typeT.Output()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	_ = os.WriteFile("./out/user.rocket.json", jsonData, 0644)
}
