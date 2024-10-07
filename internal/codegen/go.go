package codegen

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/bndrmrtn/rocket/internal/query_interpreter"
	"github.com/bndrmrtn/rocket/internal/tokenizer"
	"github.com/bndrmrtn/rocket/utils"
)

type Go struct {
	generated *tokenizer.Generated
	queries   []query_interpreter.Query
	sqlGen    func(query_interpreter.Query) string

	imports    []string
	headerout  string
	modelsout  string
	queriesout string
}

func (g *Go) Bind(gen *tokenizer.Generated, queries []query_interpreter.Query) {
	g.generated = gen
	g.queries = queries
}

func (g *Go) Get() string {
	var imports string

	if len(g.imports) > 0 {
		imports = "import ("
		for _, imp := range g.imports {
			imports += fmt.Sprintf("\n\t\"%s\"", imp)
		}
		imports += ")\n\n"
	}

	return g.headerout + imports + g.modelsout + g.queriesout
}

func (g *Go) Generate(sqlGen func(query_interpreter.Query) string) error {
	g.sqlGen = sqlGen
	pkg := os.Getenv("GOPACKAGE")
	if pkg == "" {
		return errors.New("GOPACKAGE env is required to generate Go code")
	}

	g.headerout += fmt.Sprintf("package %s\n\n", pkg)

	g.createModels()
	g.createQueries()

	return nil
}

func (g *Go) Save(file string) error {
	return nil
}

func (g *Go) createModels() error {
	for name, model := range g.generated.Models {
		err := g.createModel(name, model, g.generated.ModelKeys[name])
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Go) createModel(name string, model map[string]tokenizer.ModelConfig, keys []string) error {
	name = utils.PascalCase(name)
	var out = fmt.Sprintf("type %s struct {\n", name)
	for _, key := range keys {
		if strings.ToLower(key) == "id" {
			key = "ID"
		} else {
			key = utils.PascalCase(key)
		}

		var nullStr string
		if slices.Contains(model[key].Attributes, "nullable") {
			nullStr = "*"
		}

		out += fmt.Sprintf("\t%s %s%s", key, nullStr, g.getFieldType(model[key].Type))

		jsonKey, err := model[key].Annotations.Get("json")
		if err != nil {
			jsonKey = utils.SnakeCase(key)
		}
		out += fmt.Sprintf(" `json:\"%s\"`", jsonKey)

		out += "\n"
	}
	out += "}\n\n"

	g.modelsout += out

	return nil
}

func (g *Go) getFieldType(s string) string {
	strs := strings.SplitN(s, "(", 2)
	s = strings.ToLower(strs[0])

	switch s {
	case "number", "int":
		return "int"
	case "string", "text", "varchar":
		return "string"
	case "bool", "boolean", "bit":
		return "bool"
	case "date", "datetime", "time":
		if !slices.Contains(g.imports, "time") {
			g.imports = append(g.imports, "time")
		}
		return "time.Time"
	}

	return "string"
}

func (g *Go) createQueries() error {
	for _, query := range g.queries {
		var (
			model  = "*" + utils.PascalCase(query.From)
			multi  string
			fnName = "FindOne"
		)
		if query.MultiResult {
			multi = "[]"
			fnName = "FindMany"
			model = model[1:]
		}
		g.queriesout += fmt.Sprintf("func %s() (%s%s, error) {\n", utils.PascalCase(query.Name), multi, model)
		g.queriesout += fmt.Sprintf("\trawQuery := `%s`\n", g.sqlGen(query))
		g.queriesout += fmt.Sprintf("\treturn rocket.%s(rawQuery)\n", fnName)
		g.queriesout += "}\n\n"
	}
	return nil
}
