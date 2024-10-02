package query

type Result struct {
	rawQuery string
}

func (*Result) Execute(args map[string]string) error {
	return nil
}

func (r *Result) Query() string {
	return r.rawQuery
}
