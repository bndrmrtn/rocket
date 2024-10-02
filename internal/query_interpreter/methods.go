package query_interpreter

type Method interface {
	Is(tokens []string) bool
	Build(*Query)
	Args() [][2]string
}

type Where struct{}

func (*Where) Is(tokens []string) bool {
	return false
}

func (*Where) Token() string {
	return "where"
}

func (*Where) Args() [][2]string {
	return [][2]string{
		{"condition", "core#CondBuilder"},
	}
}
