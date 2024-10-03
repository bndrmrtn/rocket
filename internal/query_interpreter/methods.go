package query_interpreter

import "fmt"

var MethodKeywords = []string{"Where", "OrderBy", "Limit", "Offset"}

func parseWhereFunc(parantheses string, query *Query) error {
	data := explodeMultiOperations(parantheses)
	for _, d := range data {
		fmt.Println(tokenizeOperation(d))
	}
	return nil
}
