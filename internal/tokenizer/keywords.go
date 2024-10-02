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
		"settings":  SettingsType,
		"&":         LinkingSymbol,
		"@":         AnnotationSymbol,
		"//":        CommentSymbol,
		"$":         RelationSymbol,
		"algo":      Algorithm,
		"string":    DataTypeString,
		"json":      DataTypeJSON,
		"number":    DataTypeNumber,
		"date":      DataTypeDate,
		"datetime":  DataTypeDateTime,
		"primary":   PrimaryKeyAttribute,
		"increment": AutoIncrementAttribute,
		"nullable":  NullableAttribute,
		"sensitive": AnnotationSensitive,
		"hash":      AnnotationHash,
		"default":   AnnotationDefault,
		"get":       QueryGetter,
	}
}

func typeTokens() map[string]TokenType {
	return map[string]TokenType{
		"schema":   SchemaType,
		"enum":     EnumType,
		"model":    ModelType,
		"hashing":  HashingType,
		"query":    QueryType,
		"settings": SettingsType,
	}
}
