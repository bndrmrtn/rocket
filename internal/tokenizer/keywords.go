package tokenizer

func keywords() map[string]TokenType {
	return map[string]TokenType{
		"{":         BlockStart,
		"}":         BlockEnd,
		"\"":        StringLiteral,
		"schema":    SchemaType,
		"model":     ModelType,
		"hashing":   HashingType,
		"query":     QueryType,
		"&":         LinkingSymbol,
		"@":         AnnotationSymbol,
		"//":        CommentSymbol,
		"provider":  ProviderSymbol,
		"string":    DataTypeString,
		"json":      DataTypeJSON,
		"number":    DataTypeNumber,
		"date":      DataTypeDate,
		"datetime":  DataTypeDateTime,
		"primary":   PrimaryKeyAttribute,
		"increment": AutoIncrementAttribute,
		"sensitive": AnnotationSensitive,
		"hash":      AnnotationHash,
		"default":   AnnotationDefault,
		"get":       QueryGetter,
	}
}

func typeTokens() map[string]TokenType {
	return map[string]TokenType{
		"schema":  SchemaType,
		"model":   ModelType,
		"hashing": HashingType,
		"query":   QueryType,
	}
}
