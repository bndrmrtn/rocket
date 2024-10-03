package query_interpreter

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var MethodKeywords = []string{"Where", "OrderBy", "Limit", "Offset"}

func parseLimitFunc(parantheses string, query *Query) error {
	parantheses = strings.Trim(parantheses, "()")
	parantheses = strings.TrimSpace(parantheses)
	num, err := strconv.Atoi(parantheses)
	if err != nil {
		return fmt.Errorf("value error: limit is not a number: %w", err)
	}

	query.Limit = num
	return nil
}

func parseOffsetFunc(parantheses string, query *Query) error {
	parantheses = strings.Trim(parantheses, "()")
	parantheses = strings.TrimSpace(parantheses)
	num, err := strconv.Atoi(parantheses)
	if err != nil {
		return fmt.Errorf("value error: offset is not a number: %w", err)
	}

	query.Offset = num
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
		query.Order = OrderBy{Model: mData[0], Field: mData[1], Order: data[0]}
	} else {
		query.Order = OrderBy{Model: "#modelFrom#", Field: mData[0], Order: data[0]}
	}

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
		}
	}

	return nil
}
