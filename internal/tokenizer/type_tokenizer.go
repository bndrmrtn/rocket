package tokenizer

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode"
)

type TypeTokenizer struct {
	output *Generated

	tokens  []Token
	schemas []BuildToken
	hashing []BuildToken
	queries []BuildToken
}

func NewType(tokens []Token) *TypeTokenizer {
	t := &TypeTokenizer{
		output: NewGenerated(),
		tokens: tokens,
	}
	t.sortTokens()
	return t
}

func (t *TypeTokenizer) Generate() error {
	for _, token := range t.hashing {
		err := t.parseHashing(token)
		if err != nil {
			return err
		}
	}

	for _, token := range t.schemas {
		err := t.parseModels(token)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TypeTokenizer) parseModels(token BuildToken) error {
	if token.Value == "" {
		fmt.Println("WARN: empty " + token.Key)
		return nil
	}

	schema, err := t.parseModel(token)

	if err != nil {
		return err
	}

	if token.Type == ModelType {
		t.output.Models[token.Key] = schema
	}

	return nil
}

func (t *TypeTokenizer) parseModel(token BuildToken) (Model, error) {
	var schema Model = map[string]ModelConfig{}

	for i, line := range strings.Split(token.Value, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, string(LinkingSymbol)) {
			var err error
			schema, err = t.createLink(schema, line, i, token)
			if err != nil {
				return nil, err
			}
			continue
		}

		data := strings.SplitN(line, " ", 2)
		val, err := t.parseModelConfig(strings.TrimSpace(data[1]), token.Line+i+1, token.File)

		if err != nil {
			return nil, err
		}

		schema[data[0]] = val
	}

	return schema, nil
}

func (t *TypeTokenizer) parseModelConfig(s string, line int, file string) (ModelConfig, error) {
	var (
		pos  int
		Type string
	)

	for pos+1 <= len(s) {
		char := rune(s[pos])

		if unicode.IsSpace(char) {
			pos++
			break
		}

		if char == '(' {
			for pos+1 <= len(s) {
				Type += string(s[pos])
				if s[pos] == ')' {
					break
				}
				pos++
			}
			break
		}

		if unicode.IsLetter(char) || char == '_' || (char == '$' && Type == "") {
			Type += string(char)
			pos++
		}
	}

	modelConfig := ModelConfig{}

	attrs, annotations := t.getAttrsAndAnnotations(strings.TrimSpace(strings.TrimPrefix(s, Type)), line, file)

	modelConfig.Attributes = attrs
	modelConfig.Annotations = annotations

	if strings.HasPrefix(Type, "$") {
		t, args := parseFunctionCall(strings.TrimPrefix(Type, "$"))
		if !slices.Contains(RelationTypes, t) {
			return modelConfig, NewErrorWithPosition(fmt.Sprintf("%s is not a valid relation type. Try: %v", t, RelationTypes), Token{
				Line:     line,
				FileName: file,
				Value:    s,
			})
		}

		if len(args) == 2 {
			modelConfig.Type = "model"
			modelConfig.Relation = ModelRelation{
				Model: args[0],
				Field: args[1],
				Type:  t,
			}
			modelConfig.Ignore = true
		} else {
			return modelConfig, NewErrorWithPosition(fmt.Sprintf("[%s] %v arguments passed, %v is required at line %v", t, len(args), 2, line), Token{
				Line:     line,
				FileName: file,
				Value:    s,
			})
		}
	} else {
		modelConfig.Type = Type
	}

	return modelConfig, nil
}

func (t *TypeTokenizer) getAttrsAndAnnotations(s string, line int, file string) ([]string, []ModelAnnotation) {
	var (
		pos int
		buf string

		attrs       []string
		annotations []ModelAnnotation
	)

	for pos < len(s) {
		char := rune(s[pos])

		if unicode.IsSpace(char) {
			pos++
			if buf != "" {
				attrs = append(attrs, buf)
				buf = ""
			}
			continue
		}

		if string(char) == string(AnnotationSymbol) {
			var val string
			var braceCount = 0
			pos++

			for pos < len(s) {
				val += string(s[pos])
				pos++

				if unicode.IsSpace(rune(s[pos])) {
					break
				}

				if s[pos] == '(' {
					for pos < len(s) {
						if s[pos] == '(' {
							braceCount++
						}

						val += string(s[pos])
						pos++
						if s[pos-1] == ')' {
							if braceCount == 0 {
								break
							}
							braceCount--
						}
					}
					break
				}
			}

			fn, args := parseFunctionCall(val)
			annotations = append(annotations, ModelAnnotation{
				Annotation: fn,
				Arguments:  args,
			})
			continue
		}

		if unicode.IsLetter(char) || char == '_' {
			buf += string(char)
		}
		pos++
	}

	if buf != "" {
		attrs = append(attrs, buf)
	}

	return attrs, annotations
}

func (t *TypeTokenizer) getSchema(link string) (Model, error) {
	var (
		exists bool
		token  BuildToken
	)

	for _, s := range t.schemas {
		if s.Key == link {
			exists = true
			token = s
		}
	}

	if !exists {
		return nil, errors.New("Invalid schema \"" + link + "\"")
	}

	return t.parseModel(token)
}

func (t *TypeTokenizer) createLink(schema Model, line string, i int, token BuildToken) (Model, error) {
	if strings.Contains(line, " ") {
		return nil, NewErrorWithPosition("Attributes cannot be added to a schema linking", Token{
			Value:     line,
			TokenType: token.Type,
			Line:      token.Line + i,
			FileName:  token.File,
		})
	}
	data, err := t.getSchema(strings.TrimPrefix(line, string(LinkingSymbol)))

	if err != nil {
		return nil, err
	}

	for key, value := range data {
		schema[key] = value
	}

	return schema, nil
}

func (t *TypeTokenizer) parseHashing(token BuildToken) error {
	if token.Value == "" {
		fmt.Println("WARN: empty", token.Key)
		return nil
	}

	var (
		algoList = getHashAlgoList()
		algo     string
		args     []string
	)

	for i, line := range strings.Split(token.Value, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, string(Algorithm)) {
			line = strings.TrimPrefix(line, string(Algorithm))
		} else {
			return NewErrorWithPosition("Invalid data", Token{
				Value:     line,
				TokenType: token.Type,
				Line:      token.Line + i,
				FileName:  token.File,
			})
		}

		line = strings.TrimSpace(line)

		algo, args = parseFunctionCall(line)
		if !slices.Contains(algoList, algo) {
			return NewErrorWithPosition("Invalid hashing algo", Token{
				Value:     line,
				TokenType: token.Type,
				Line:      token.Line + i + 1,
				FileName:  token.File,
			})
		}
	}

	if algo == "" {
		return nil
	}

	hashing := Hashing{
		Name:      token.Key,
		Provider:  algo,
		Arguments: args,
	}
	t.output.Hashing = append(t.output.Hashing, hashing)

	return nil
}

func (t *TypeTokenizer) sortTokens() {
	for i, token := range t.tokens {
		if token.TokenType == TypeValue {
			continue
		}
		switch token.TokenType {
		case SchemaType, ModelType:
			t.schemas = append(t.schemas, t.createToken(i, token))
		case HashingType:
			t.hashing = append(t.schemas, t.createToken(i, token))
		case QueryType:
			t.queries = append(t.schemas, t.createToken(i, token))
		}
	}
}

func (t *TypeTokenizer) createToken(inx int, token Token) BuildToken {
	var val string

	if inx+1 < len(t.tokens) && t.tokens[inx+1].TokenType == TypeValue {
		val = t.tokens[inx+1].Value
	}

	return BuildToken{
		Type:  token.TokenType,
		Key:   token.Value,
		Value: val,
		Line:  token.Line,
		File:  token.FileName,
	}
}

func (t *TypeTokenizer) Output() *Generated {
	return t.output
}
