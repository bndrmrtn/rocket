package query

type GetterResult struct {
	Fields      map[string][]string
	MultiResult bool
}

type Conditions struct {
	DefaultModel string
	Where        WhereCondition
	Limit        int
	Offset       int
}

type WhereCondition struct {
	Model    string
	Operator string
	Value    []string
	Compare  string
}
