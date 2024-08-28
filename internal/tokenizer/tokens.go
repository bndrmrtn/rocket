package tokenizer

type TokenType string

const (
	BlockStart TokenType = "{"
	BlockEnd   TokenType = "}"

	StringLiteral TokenType = "\""

	SchemaType  TokenType = "schema"
	ModelType   TokenType = "model"
	HashingType TokenType = "hashing"
	QueryType   TokenType = "query"
	TypeValue   TokenType = "#typeValue#"

	LinkingSymbol    TokenType = "&"
	AnnotationSymbol TokenType = "@"
	CommentSymbol    TokenType = "//"
	ProviderSymbol   TokenType = "provider"

	DataTypeString   TokenType = "string"
	DataTypeJSON     TokenType = "json"
	DataTypeNumber   TokenType = "number"
	DataTypeDate     TokenType = "date"
	DataTypeDateTime TokenType = "datetime"

	PrimaryKeyAttribute    TokenType = "primary"
	AutoIncrementAttribute TokenType = "increment"

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
