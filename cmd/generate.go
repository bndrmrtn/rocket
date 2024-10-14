package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bndrmrtn/rocket/internal/codegen"
	"github.com/bndrmrtn/rocket/internal/generator"
	"github.com/bndrmrtn/rocket/internal/query_interpreter"
	"github.com/bndrmrtn/rocket/internal/schemagen"
	"github.com/bndrmrtn/rocket/internal/tokenizer"
	"github.com/bndrmrtn/rocket/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g", "gen"},
	Short:   "Generate code from a Rocket file.",
	Run:     execGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("file", "f", "", "Rocket file to generate code from.")
	generateCmd.Flags().StringP("language", "l", "go", "Language to generate code in.")
	generateCmd.Flags().StringP("database", "d", "mysql", "Database to generate code for.")
	generateCmd.Flags().StringP("out", "o", "{filename}_rgen.{ext}", "Output file for generated code.")
	generateCmd.Flags().BoolP("no-color", "c", false, "Disable colors.")

	_ = generateCmd.MarkFlagRequired("file")
}

func execGenerate(cmd *cobra.Command, args []string) {
	start := time.Now()
	defer func(start time.Time) {
		fmt.Printf("Generated in %s\n", time.Since(start))
	}(start)
	fmt.Println("🚀 Rocket - Generate Code")

	file := cmd.Flag("file").Value.String()
	out := cmd.Flag("out").Value.String()
	language := cmd.Flag("language").Value.String()
	database := cmd.Flag("database").Value.String()
	noColor := cmd.Flag("no-color").Value.String()
	if noColor == "true" {
		color.NoColor = true
	}

	if strings.Contains(out, "{filename}") {
		outPath := strings.TrimSuffix(file, filepath.Ext(file))
		out = strings.ReplaceAll(out, "{filename}", outPath)
	}

	lang, err := generator.GetLanguage(language)
	if err != nil {
		log.Fatal(err)
	}

	_, err = generator.GetDatabase(database)
	if err != nil {
		log.Fatal(err)
	}

	tokens, hash, err := getTokens(file)
	if err != nil {
		log.Fatalf("Failed to get tokens from file: %v", err)
	}

	if len(tokens) == 0 {
		fmt.Println("No tokens. Nothing to do.")
		os.Exit(0)
	}

	fmt.Println("Tokens Hash (sha256): " + hash)

	typeT := tokenizer.NewType(tokens)

	err = typeT.Generate()
	if err != nil {
		log.Fatal(err)
	}

	data := typeT.Output()

	db, err := schemagen.GetDB(database)
	if err != nil {
		log.Fatal(err)
	}

	db.Bind(data)
	sqlDump, err := db.Get()
	if err != nil {
		log.Fatal("Failed to generate SQL Code: ", err)
	}
	success("SQL Code generated from models successfully.")

	b, _ := yaml.Marshal(data)
	os.WriteFile("./out/queries.yaml", b, os.ModePerm)

	interpreter := query_interpreter.NewInterpreter(data)
	queries, err := interpreter.InterpretAll()
	if err != nil {
		log.Fatal("Failed to interpret query tokens: ", err)
	}

	success("Queries interpreted successfully.")

	cg, err := codegen.GetLang(lang.Extension)
	if err != nil {
		log.Fatal(err)
	}

	cg.Bind(sqlDump, data, queries)
	err = cg.Generate(db.GetQueryParser())
	if err != nil {
		log.Fatal("Failed to generate output: ", err)
	}
	cg.Save(strings.ReplaceAll(out, "{ext}", lang.Extension))
}

// Helpers

func getTokens(file string) ([]tokenizer.Token, string, error) {
	var (
		files  []string
		tokens []tokenizer.Token
	)

	st, err := os.Stat(file)
	if err != nil {
		return nil, "", err
	}

	if st.IsDir() {
		files, err = utils.WalkDir(file, "rocket")
		if err != nil {
			return nil, "", err
		}
	} else {
		files = []string{file}
	}

	for _, file := range files {
		t, err := tokenizer.New(file)
		if err != nil {
			return nil, "", err
		}

		err = t.Tokenize()
		if err != nil {
			return nil, "", err
		}

		tokens = append(tokens, t.GetTokens()...)
	}

	var b bytes.Buffer
	_ = gob.NewEncoder(&b).Encode(tokens)

	h := sha256.New()
	_, err = h.Write(b.Bytes())

	return tokens, hex.EncodeToString(h.Sum(nil)), nil
}
