package query_interpreter

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/bndrmrtn/rocket/internal/tokenizer"
)

func parseLimitFunc(parantheses string, query *Query, args []tokenizer.QueryArg) error {
	parantheses = strings.Trim(parantheses, "()")
	parantheses = strings.TrimSpace(parantheses)
	_, err := strconv.Atoi(parantheses)
	if err != nil {
		var ok bool
		for _, arg := range args {
			if arg.Name == parantheses {
				if arg.Type == "number" {
					query.Limit = arg.Name
					ok = true
				} else {
					return fmt.Errorf("param error: limit is not number: %s, type: %s", arg.Name, arg.Type)
				}
			}
		}

		if !ok {
			return fmt.Errorf("value error: limit is not a number: %w", err)
		}
	}

	query.Limit = parantheses
	return nil
}

func parseOffsetFunc(parantheses string, query *Query, args []tokenizer.QueryArg) error {
	parantheses = strings.Trim(parantheses, "()")
	parantheses = strings.TrimSpace(parantheses)
	_, err := strconv.Atoi(parantheses)
	if err != nil {
		var ok bool
		for _, arg := range args {
			if arg.Name == parantheses {
				if arg.Type == "number" {
					query.Limit = arg.Name
					ok = true
				} else {
					return fmt.Errorf("param error: offset is not number: %s, type: %s", arg.Name, arg.Type)
				}
			}
		}

		if !ok {
			return fmt.Errorf("value error: offset is not a number: %w", err)
		}
	}

	query.Offset = parantheses
	return nil
}

func parseOrderByFunc(parantheses string, query *Query) error {
	parantheses = strings.Trim(parantheses, "()")
	parantheses = strings.TrimSpace(parantheses)
	var types = []string{"asc", "desc"}
	data := strings.SplitN(parantheses, " ", 2)

	if len(data) != 2 {
		return fmt.Errorf("value error: order by must have a field and an order type keyword (asc/desc) - OrderBy(asc Model.Field)")
	}

	if !slices.Contains(types, data[0]) {
		return fmt.Errorf("value error: order type must be asc or desc")
	}

	modelField := data[1]
	mData := strings.SplitN(modelField, " . ", 2)
	if len(mData) == 2 {
		query.Order = append(query.Order, OrderBy{Model: mData[0], Field: mData[1], Order: data[0]})
	} else {
		query.Order = append(query.Order, OrderBy{Model: "#modelFrom#", Field: mData[0], Order: data[0]})
	}

	return nil
}

func parseOrderByRandFunc(parantheses string, query *Query) error {
	data := strings.Trim(parantheses, "()")
	data = strings.TrimSpace(data)
	if data != "" {
		return errors.New("parameter error: order by rand does not accept any parameters")
	}
	query.Order = append(query.Order, OrderBy{Order: "rand()"})
	return nil
}

// TODO: implement this function
func parseWhereFunc(parantheses string, _ *Query) error {
	var cond ConditionBuilder

	data := explodeMultiOperations(parantheses)
	for _, d := range data {
		if slices.Contains(Operators(), Operator(d)) {
			cond = append(cond, ConditionType{Type: "operator", Operator: d})
			continue
		} else {
			err := parseWhereCondition(d, &cond)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func parseWhereCondition(d string, _ *ConditionBuilder) error {
	_ = tokenizeOperation(d)
	return nil
}
