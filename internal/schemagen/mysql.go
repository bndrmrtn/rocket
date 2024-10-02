package schemagen

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/bndrmrtn/rocket/internal/tokenizer"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type mysql struct {
	data         *tokenizer.Generated
	out          string
	helperTables string
}

func (m *mysql) Bind(data *tokenizer.Generated) {
	m.data = data
}

func (m *mysql) Create(out string) error {
	for name, model := range m.data.Models {
		err := m.createTable(name, model, m.data.ModelKeys[name])
		if err != nil {
			return err
		}
	}
	return os.WriteFile(out, []byte(m.out), os.ModePerm)
}

func (m *mysql) createTable(name string, model map[string]tokenizer.ModelConfig, keys []string) error {
	var out = "CREATE TABLE IF NOT EXISTS `" + name + "` (\n"
	for _, field := range keys {
		out += m.getField(name, field, model[field], false)
	}
	out = strings.TrimSuffix(out, ",\n")
	out += "\n);"
	m.out += out + "\n\n"
	return nil
}

func (m *mysql) getField(name, field string, config tokenizer.ModelConfig, disablePrimaryCreate bool) string {
	var out string
	if config.Type != "model" {
		t := m.getType(config.Type)
		out += "\t`" + field + "` " + t
		if !disablePrimaryCreate && slices.Contains(config.Attributes, "increment") {
			out += " AUTO_INCREMENT"
		}
		if !slices.Contains(config.Attributes, "nullable") {
			out += " NOT NULL"
		}
		out += m.getDefaultValue(config.Annotations)
		out += ",\n"
		if !disablePrimaryCreate {
			if slices.Contains(config.Attributes, "primary") {
				out += "\tPRIMARY KEY (`" + field + "`),\n"
			}
		}
	} else {
		if rel := config.Relation; rel != nil {
			if rel.Field != "#" {
				out += m.getRelationQuery(name, field, config)
			}
		}
	}
	return out
}

func (m *mysql) getType(s string) string {
	if strings.HasPrefix(s, "enum:") {
		return m.getEnumType(strings.TrimPrefix(s, "enum:"))
	}

	var (
		tp   string = s
		size string
	)
	if strings.Contains(s, "(") && strings.HasSuffix(s, ")") {
		sp := strings.SplitN(strings.TrimSuffix(s, ")"), "(", 2)
		if len(sp) == 2 {
			tp = sp[0]
			size = sp[1]
		}
	}

	if size != "" {
		size = "(" + size + ")"
	}

	switch tp {
	case "string", "text", "varchar":
		if size == "" {
			return "TEXT"
		}
		return "VARCHAR" + size
	case "number", "int":
		if size == "" {
			size = "(32)"
		}
		return "INT" + size
	case "bool", "boolean", "bit":
		if size == "" {
			size = "(1)"
		}
		return "BIT" + size
	case "datetime", "date", "time":
		return strings.ToUpper(s)
	}

	return tp + size
}

func (m *mysql) getEnumType(s string) string {
	enum, ok := m.data.Enums[s]
	if !ok {
		log.Fatal("Enum does not exists:", s)
	}

	var out = "ENUM("

	for _, val := range enum {
		out += "'" + val + "',"
	}
	out = strings.TrimSuffix(out, ",")
	out += ")"

	return out
}

func (m *mysql) getDefaultValue(anns []tokenizer.ModelAnnotation) string {
	var out string

	for _, ann := range anns {
		if ann.Annotation == "default" && len(ann.Arguments) == 1 {
			val := ann.Arguments[0]
			if strings.HasPrefix(val, "&") {
				val = m.getEnumField(strings.TrimPrefix(val, "&"))
			}
			if val == "now()" {
				val = "NOW()"
			}
			out = " DEFAULT " + val
			break
		}
	}

	return out
}

func (m *mysql) getEnumField(e string) string {
	parts := strings.SplitN(e, ".", 2)
	if len(parts) != 2 {
		return e
	}

	enum, ok := m.data.Enums[parts[0]]
	if !ok {
		log.Fatal("Enum not found:", parts[0])
	}

	field, ok := enum[parts[1]]
	if !ok {
		log.Fatal("Field not found for enum:", parts[0], parts[1])
	}

	return "'" + field + "'"
}

func (m *mysql) getRelationQuery(model string, field string, config tokenizer.ModelConfig) string {
	var (
		out string

		constraint bool
		onDelete   string
		onUpdate   string

		relatedModel string = config.Relation.Model
		relatedField string = config.Relation.Field
		// relationType         string = config.Relation.Type
		madeRelatedFieldName string
	)

	constraint = slices.Contains(config.Attributes, "constraint")
	if constraint {
		for _, ann := range config.Annotations {
			switch ann.Annotation {
			case "cascadeOnDelete":
				onDelete = "cascade"
			case "cascadeOnUpdate":
				onUpdate = "cascade"
			case "onDelete", "onUpdate":
				if len(ann.Arguments) != 1 {
					log.Fatal("Not enough argument for onDelete/onUpdate annotation")
				}
				if ann.Annotation == "onDelete" {
					onDelete = strings.Trim(ann.Arguments[0], "\"")
				} else if ann.Annotation == "onUpdate" {
					onUpdate = strings.Trim(ann.Arguments[0], "\"")
				}
				break
			}
		}
	}

	madeRelatedFieldName = strings.ToLower(field) + "_" + strings.ToLower(relatedField)

	var relatedModelConfig tokenizer.ModelConfig
	mv, ok := m.data.Models[relatedModel]
	if !ok {
		log.Fatal("related model does not exists", relatedModel)
	}
	relatedModelConfig, ok = mv[relatedField]
	if !ok {
		log.Fatal("related model field does not exists", relatedModel, relatedField)
	}

	out += m.getField(relatedModel, madeRelatedFieldName, relatedModelConfig, true)

	if constraint {
		caser := cases.Title(language.English)
		out += fmt.Sprintf("\tCONSTRAINT `fk_%s%s`", caser.String(strings.ToLower(model)), caser.String(strings.ToLower(relatedModel)))
		out += " "
	}

	if !strings.Contains(out, "CONSTRAINT") {
		out += "\t"
	}
	out += fmt.Sprintf("FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`)", madeRelatedFieldName, relatedModel, relatedField)

	if constraint {
		if onDelete != "" || onUpdate != "" {
			out += "\n\t"
		}

		if onDelete != "" {
			out += " ON DELETE " + strings.ToUpper(onDelete)
		}
		if onUpdate != "" {
			out += " ON UPDATE " + strings.ToUpper(onUpdate)
		}
	}

	return out + ",\n"
}
