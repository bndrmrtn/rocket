package query_interpreter

type Keyword string

const (
	Equal Keyword = "="
	In    Keyword = "in"
)
