package query_interpreter

type Keyword string

const (
	MultiResult Keyword = "[]"
)

var QueryTypeKeywords = []string{"get"}

type Operator string

const (
	Equal        Operator = "=="
	NotEqual     Operator = "!="
	Greater      Operator = ">"
	GreaterEqual Operator = ">="
	Less         Operator = "<"
	LessEqual    Operator = "<="
	In           Operator = "in"
	NotIn        Operator = "not in"
	Like         Operator = "??"
)
