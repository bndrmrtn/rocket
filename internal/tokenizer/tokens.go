package tokenizer

type TokenType string

const (
	BlockStart TokenType = "{"
	BlockEnd   TokenType = "}"

	StringLiteral TokenType = "\""

	SchemaType   TokenType = "schema"
	ModelType    TokenType = "model"
	EnumType     TokenType = "enum"
	HashingType  TokenType = "hashing"
	QueryType    TokenType = "query"
	SettingsType TokenType = "settings"
	TypeValue    TokenType = "#typeValue#"

	Algorithm TokenType = "algo"

	LinkingSymbol    TokenType = "&"
	AnnotationSymbol TokenType = "@"
	CommentSymbol    TokenType = "//"
	RelationSymbol   TokenType = "$"

	DataTypeString   TokenType = "string"
	DataTypeJSON     TokenType = "json"
	DataTypeNumber   TokenType = "number"
	DataTypeDate     TokenType = "date"
	DataTypeDateTime TokenType = "datetime"

	PrimaryKeyAttribute    TokenType = "primary"
	AutoIncrementAttribute TokenType = "increment"
	NullableAttribute      TokenType = "nullable"

	AnnotationSensitive TokenType = "sensitive"
	AnnotationHash      TokenType = "hash"
	AnnotationDefault   TokenType = "default"

	QueryGetter TokenType = "query"
)

type Token struct {
	Value     string
	TokenType TokenType
	TokenPos  int
	FileName  string
	Line      int
}

type BuildToken struct {
	Type  TokenType
	Key   string
	Value string
	Line  int
	File  string
}

func (b BuildToken) ToToken() Token {
	return Token{
		Value:     b.Value,
		TokenType: b.Type,
		TokenPos:  b.Line,
		FileName:  b.File,
		Line:      b.Line,
	}
}
