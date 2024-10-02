package query

import (
	"fmt"

	"github.com/bndrmrtn/rocket/internal/tokenizer"
)

type Interpreter struct {
	data *tokenizer.Generated
	pos  int
	src  []string
}

func NewInterpreter(data *tokenizer.Generated) *Interpreter {
	return &Interpreter{
		data: data,
	}
}

func (i *Interpreter) Interpret(name string) (*Result, error) {
	query, err := i.query(name)
	if err != nil {
		return nil, err
	}
	// setup interpreter
	i.pos = 0
	i.src = query.Tokens

	var (
		getter *GetterResult
	)

	for i.pos < len(i.src) {
		val := i.src[i.pos]

		if val == Get.String() {
			i.pos++
			getter, err = i.parseGetter()
		}
	}

	fmt.Println(getter)

	return nil, nil
}

func (i *Interpreter) parseGetter() (*GetterResult, error) {
	var res GetterResult

	if len(i.src) > 1 && i.src[0] == "[" && i.src[0+1] == "]" {
		res.MultiResult = true
	}

	for i.nextOK() {
		if res.MultiResult && i.pos < 2 {
			continue
		}

		val := i.src[i.pos]

		if val == "{" {
			var buf []string
			for i.nextOK() {
				buf = append(buf, i.src[i.pos])
				if i.src[i.pos] == "}" {
					break
				}
				i.pos++
			}
			res.Fields = i.parseFields(buf)
		}

		conditions, _ := i.parseConditions(i.src[i.pos:], i.getModelNames())
		fmt.Println(conditions)
		break
	}

	fmt.Println(res.Fields)

	return nil, nil
}
