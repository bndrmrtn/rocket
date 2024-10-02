package query

import (
	"errors"
	"fmt"
	"github.com/bndrmrtn/rocket/internal/tokenizer"
	"slices"
)

func (i *Interpreter) query(name string) (tokenizer.Query, error) {
	for _, q := range i.data.Queries {
		if q.Name == name {
			return q, nil
		}
	}
	return tokenizer.Query{}, errors.New("query does not exist")
}

func (i *Interpreter) nextOK() bool {
	return i.pos+1 <= len(i.src)
}

func (i *Interpreter) parseFields(f []string) map[string][]string {
	fields := make(map[string][]string)
	var (
		pos int
	)

	for pos+1 < len(f) {
		q := f[pos]
		pos++

		if q == "{" || q == "," {
			continue
		}

		if q == "}" {
			break
		}

		if pos+1 < len(f) && f[pos] == "." {
			fields[f[pos-1]] = append(fields[f[pos-1]], f[pos+1])
			pos += 2
			continue
		}

		fields["#"] = append(fields["#"], f[pos-1])
	}

	fmt.Println(fields)

	return fields
}

func (i *Interpreter) parseConditions(tokens []string, models []string) (*Conditions, error) {
	var (
		con Conditions
		pos int
	)

	for pos+1 < len(tokens) {
		token := tokens[pos]
		pos++

		if pos == 1 {
			if !slices.Contains(models, token) {
				return nil, errors.New(token + " model does not exist")
			}
			con.DefaultModel = token
			continue
		}

		if token == "." {
			continue
		}

		if token == Where.String() {
			pos++
			var bracketCount int

			for pos+1 < len(tokens) {
				if tokens[pos] == "(" {
					bracketCount++
					continue
				}

				if tokens[pos] == ")" {
					bracketCount--
					if bracketCount == 0 {
						break
					}
				}

				if pos+1 < len(token) && tokens[pos] == "." {
					pos += 2
					continue
				}
			}

			continue
		}
	}

	return &con, nil
}

func (i *Interpreter) getModelNames() []string {
	var keys = make([]string, len(i.data.Models))
	for key := range i.data.Models {
		keys = append(keys, key)
	}
	return keys
}
