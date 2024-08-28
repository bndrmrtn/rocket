package main

import (
	"encoding/json"
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
		log.Fatal(err)
	}

	tokens := t.GetTokens()

	typeT := tokenizer.NewType(tokens)
	err = typeT.Generate()
	if err != nil {
		log.Fatal(err)
	}

	data := typeT.Output()
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	_ = os.WriteFile("./out/user.rocket.json", jsonData, 0644)
}
