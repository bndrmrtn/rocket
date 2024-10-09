package query_interpreter

type Query struct {
	Name        string
	Fields      map[string][]string
	MultiResult bool
	From        string
	Limit       string
	Offset      string
	Conditions  []ConditionBuilder
	Order       []OrderBy
	FuncParams  []FuncParam
}

type FuncParam struct {
	Name string
	Type string
}

// ConditionBuilder holds more ConditionTypes
type ConditionBuilder []ConditionType

// ConditionType is a struct that can be an operator or a condition
// For example: (a == b) && (c == d)
// In this example, (a == b) is a condition and && is an operator
type ConditionType struct {
	// Type can be an operator or a condition
	Type string
	// Operator is an operator that separates two conditions
	Operator string
	// Cond is a Condition
	Cond *Condition
}

type Condition struct {
	Model string
	// Operator can be ==, !=, >, <, >=, <=, in, not in, like, not like
	// May be updated in the future
	Operator string
	// Value (or values) can be a single value or multiple values
	Value []string
	// SingleValue is a boolean that is used to determine if the value is a single value
	SingleValue bool
	// Compare is a string that is used to compare the value with the field
	Compare string
	// CompareArg is a bool that is used to determine if the value is a func argument or not
	CompareArg bool
}

type OrderBy struct {
	Model string
	Field string
	Order string
}

func NewQuery() *Query {
	return &Query{
		Fields:      make(map[string][]string),
		MultiResult: false,
	}
}
