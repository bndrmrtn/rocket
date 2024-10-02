package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bndrmrtn/rocket/internal/generator"
	"github.com/bndrmrtn/rocket/internal/schemagen"
	"github.com/bndrmrtn/rocket/internal/tokenizer"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g", "gen"},
	Short:   "Generate code from a Rocket file.",
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()
		out := cmd.Flag("output").Value.String()
		language := cmd.Flag("language").Value.String()
		database := cmd.Flag("database").Value.String()

		if strings.Contains(out, "*") {
			outPath := strings.TrimSuffix(file, filepath.Ext(file))
			out = strings.ReplaceAll(out, "*", outPath)
		}

		_, err := generator.GetLanguage(language)
		if err != nil {
			log.Fatal(err)
		}

		/*if strings.Contains(out, "{ext}") {
		out = strings.ReplaceAll(out, "{ext}", lang.Extension)
		}*/

		_, err = generator.GetDatabase(database)
		if err != nil {
			log.Fatal(err)
		}

		t, err := tokenizer.New(file)
		if err != nil {
			log.Fatal(err)
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

		db, err := schemagen.GetDB(database)
		if err != nil {
			log.Fatal(err)
		}

		d, _ := yaml.Marshal(data)
		_ = os.WriteFile("./out/data.yaml", d, os.ModePerm)

		db.Bind(data)
		db.Create(strings.ReplaceAll(out, "{ext}", "sql"))
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("file", "f", "", "Rocket file to generate code from.")
	generateCmd.Flags().StringP("language", "l", "go", "Language to generate code in.")
	generateCmd.Flags().StringP("database", "d", "mysql", "Database to generate code for.")
	generateCmd.Flags().StringP("output", "o", "*.{ext}", "Output file for generated code.")

	generateCmd.MarkFlagRequired("file")
}
