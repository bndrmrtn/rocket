package query

type TokenType string

const (
	Get TokenType = "get"
)

func (t TokenType) String() string {
	return string(t)
}

type Method string

const (
	Where       Method = "where"
	Limit       Method = "limit"
	Offset      Method = "offset"
	OrderBy     Method = "orderBy"
	OrderByDesc Method = "orderByDesc"
	OrderByAsc  Method = "orderByAsc"
	Random      Method = "random"
)

func (m Method) String() string {
	return string(m)
}
