package tokenizer

import (
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
	enums   []BuildToken
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
	for _, token := range t.enums {
		err := t.parseEnum(token)
		if err != nil {
			return err
		}

	}

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

	query := NewQueryTokenizer(t.queries)
	err := query.Generate()
	if err != nil {
		return err
	}

	t.output.Queries = query.Output()

	return nil
}

func (t *TypeTokenizer) parseEnum(token BuildToken) error {
	var enum = make(map[string]string)

	data := strings.Split(token.Value, "\n")
	for _, line := range data {
		parts := strings.Split(line, "//")
		if len(parts) > 1 {
			line = parts[0]
		}

		parts = strings.Split(line, "=")
		var (
			key, value string
			err        error
		)
		if len(parts) == 1 {
			key = strings.TrimSpace(parts[0])
			value = strings.TrimSpace(parts[0])
		} else if len(parts) == 2 {
			key = strings.TrimSpace(parts[0])
			value, err = parseValue(parts[1])
			if err != nil {
				return NewErrorWithPosition(fmt.Sprintf("invalid value: \"%s\"", value), token.ToToken()).SetType("syntax error")
			}
		} else {
			return NewErrorWithPosition("invalid length of items between \"=\"", token.ToToken()).SetType("syntax error")
		}

		if key == "" {
			continue
		}

		enum[key] = value
	}

	t.output.Enums[token.Key] = enum

	return nil
}

func (t *TypeTokenizer) parseModels(token BuildToken) error {
	if token.Value == "" {
		WarnWithPos("empty code block", Token{
			FileName: token.File,
			Line:     token.Line,
		})
		return nil
	}

	schema, keys, err := t.parseModel(token)

	if err != nil {
		return err
	}

	if token.Type == ModelType {
		t.output.Models[token.Key] = schema
		t.output.ModelKeys[token.Key] = keys
	}

	return nil
}

func (t *TypeTokenizer) parseModel(token BuildToken) (Model, []string, error) {
	var (
		schema Model = map[string]ModelConfig{}
		keys   []string
	)

	for i, line := range strings.Split(token.Value, "\n") {
		if strings.Contains(line, "//") {
			sp := strings.SplitN(line, "//", 2)
			if len(sp) > 1 {
				line = sp[0]
			}
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, string(LinkingSymbol)) {
			var (
				err error
				k   []string
			)
			schema, k, err = t.createLink(schema, line, i, token)
			if err != nil {
				return nil, nil, err
			}
			keys = append(keys, k...)
			continue
		}

		data := strings.SplitN(line, " ", 2)
		val, err := t.parseModelConfig(strings.TrimSpace(data[1]), token.Line+i+1, token.File)

		if err != nil {
			return nil, nil, err
		}

		keys = append(keys, data[0])
		schema[data[0]] = val
	}

	return schema, keys, nil
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

		if unicode.IsLetter(char) || char == '&' || char == '_' || (char == '$' && Type == "") {
			Type += string(char)
			pos++
		}
	}

	modelConfig := ModelConfig{}

	attrs, annotations := t.getAttrsAndAnnotations(strings.TrimSpace(strings.TrimPrefix(s, Type)), line, file)

	modelConfig.Attributes = attrs
	modelConfig.Annotations = annotations

	if strings.HasPrefix(Type, "$") {
		rel, args := parseFunctionCall(strings.TrimPrefix(Type, "$"))
		if !slices.Contains(RelationTypes, rel) {
			return modelConfig, NewErrorWithPosition(fmt.Sprintf("%v is not a valid relation type. try: %v", t, RelationTypes), Token{
				Line:     line,
				FileName: file,
				Value:    s,
			}).SetType("invalid relation")
		}

		err := t.checkRelationArgs(rel, args, Token{
			Line:     line,
			FileName: file,
			Value:    s,
		})

		if err != nil {
			return modelConfig, err
		}

		if len(args) == 1 {
			args = append(args, "#")
		}

		modelConfig.Type = "model"
		modelConfig.Relation = &ModelRelation{
			Model: args[0],
			Field: args[1],
			Type:  rel,
		}
		modelConfig.Ignore = true
	} else if strings.HasPrefix(Type, "&") {
		enumName := strings.TrimPrefix(Type, "&")
		_, ok := t.output.Enums[enumName]
		if !ok {
			return modelConfig, NewErrorWithPosition(fmt.Sprintf("enum \"%s\" not found", enumName), Token{
				Line:     line,
				FileName: file,
				Value:    s,
			}).SetType("typegetter error")
		}
		modelConfig.Type = "enum:" + enumName
	} else {
		modelConfig.Type = Type
	}

	return modelConfig, nil
}

func (t *TypeTokenizer) getAttrsAndAnnotations(s string, _ int, _ string) ([]string, []ModelAnnotation) {
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

				if pos >= len(s) {
					continue
				}

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

func (t *TypeTokenizer) getSchema(link string, bt BuildToken) (Model, []string, error) {
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
		return nil, nil, NewErrorWithPosition(fmt.Sprintf("invalid schema \"%s\"", link), bt.ToToken()).SetType("typegetter error")
	}

	return t.parseModel(token)
}

func (t *TypeTokenizer) createLink(schema Model, line string, i int, token BuildToken) (Model, []string, error) {
	if strings.Contains(line, " ") {
		return nil, nil, NewErrorWithPosition("attributes cannot be added to a schema linking", Token{
			Value:     line,
			TokenType: token.Type,
			Line:      token.Line + i,
			FileName:  token.File,
		}).SetType("schema error")
	}
	data, keys, err := t.getSchema(strings.TrimPrefix(line, string(LinkingSymbol)), token)

	if err != nil {
		return nil, nil, err
	}

	for key, value := range data {
		schema[key] = value
	}

	return schema, keys, nil
}

func (t *TypeTokenizer) parseHashing(token BuildToken) error {
	if token.Value == "" {
		WarnWithPos("empty code block", Token{
			FileName: token.File,
			Line:     token.Line,
		})
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
			return NewErrorWithPosition("invalid data", Token{
				Value:     line,
				TokenType: token.Type,
				Line:      token.Line + i,
				FileName:  token.File,
			}).SetType("hash algo error")
		}

		line = strings.TrimSpace(line)

		algo, args = parseFunctionCall(line)
		if !slices.Contains(algoList, algo) {
			return NewErrorWithPosition("unsupported hashing algorithm", Token{
				Value:     line,
				TokenType: token.Type,
				Line:      token.Line + i + 1,
				FileName:  token.File,
			}).SetType("hash algo error")
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
			t.hashing = append(t.hashing, t.createToken(i, token))
		case QueryType:
			t.queries = append(t.queries, t.createToken(i, token))
		case EnumType:
			t.enums = append(t.enums, t.createToken(i, token))
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

func (t *TypeTokenizer) checkRelationArgs(relation string, args []string, token Token) error {
	if slices.Contains([]string{"belongsTo"}, relation) {
		if len(args) != 2 {
			return NewErrorWithPosition(fmt.Sprintf("[%s] %v arguments passed, %v is required", relation, len(args), 2), token).SetType("relation error")
		}
	}
	return nil
}

func (t *TypeTokenizer) Output() *Generated {
	return t.output
}
