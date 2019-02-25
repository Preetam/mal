package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Preetam/mal/dl/types"
)

func Print(val types.MalType) string {
	switch malVal := val.(type) {
	case types.MalBool:
		if malVal {
			return "true"
		}
		return "false"
	case types.MalInt:
		return strconv.FormatInt(int64(malVal), 10)
	case types.MalSymbol:
		return string(malVal)
	case types.MalList:
		elems := []string{}
		for _, elem := range malVal {
			elems = append(elems, Print(elem))
		}
		return "(" + strings.Join(elems, " ") + ")"
	case nil:
		return "nil"
	case types.MalFunction:
		return "#<function>"
	case types.MalString:
		return string(malVal)
	}
	panic("unknown type " + fmt.Sprintf("%T", val))
}
