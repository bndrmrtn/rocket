package query_interpreter

type Query struct {
	Fields      map[string][]string
	MultiResult bool
	From        string
	Limit       int
	Offset      int
	Conditions  []Condition
}

type Condition struct {
	Model       string
	Operator    string
	Value       []string
	Compare     string
	SingleValue bool
}

func NewQuery() *Query {
	return &Query{
		Fields:      make(map[string][]string),
		MultiResult: false,
	}
}
