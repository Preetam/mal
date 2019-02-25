package core

import (
	"errors"

	"github.com/Preetam/mal/dl/printer"
	"github.com/Preetam/mal/dl/types"
)

var NS = map[string]types.MalType{
	"+": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a, b, err := getTwoMalInts(args)
		if err != nil {
			return nil, err
		}
		return a + b, nil
	}),
	"-": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a, b, err := getTwoMalInts(args)
		if err != nil {
			return nil, err
		}
		return a - b, nil
	}),
	"*": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a, b, err := getTwoMalInts(args)
		if err != nil {
			return nil, err
		}
		return a * b, nil
	}),
	"/": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a, b, err := getTwoMalInts(args)
		if err != nil {
			return nil, err
		}
		return a / b, nil
	}),
	"prn": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		if len(args) == 0 {
			return nil, nil
		}
		str := printer.Print(args[0])
		return types.MalSymbol(str), nil
	}),
	"pr-str": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		if len(args) == 0 {
			return nil, nil
		}
		str := printer.Print(args[0])
		return types.MalString(str), nil
	}),
	"list": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		return types.MalList(args), nil
	}),
	"list?": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		if _, ok := args[0].(types.MalList); ok {
			return types.MalBool(true), nil
		}
		return types.MalBool(false), nil
	}),
	"empty?": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		if _, ok := args[0].(types.MalList); ok {
			return types.MalBool(len(args[0].(types.MalList)) == 0), nil
		}
		return types.MalBool(false), nil
	}),
	"count": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		if _, ok := args[0].(types.MalList); ok {
			return types.MalInt(len(args[0].(types.MalList))), nil
		}
		return types.MalInt(0), nil
	}),
	"=": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a1, a2 := args[0], args[1]
		var areEqual func(a, b types.MalType) types.MalBool

		areEqual = func(a, b types.MalType) types.MalBool {
			if a == nil || b == nil {
				return types.MalBool(false)
			}
			if a.IsMalType() != b.IsMalType() {
				return types.MalBool(false)
			}
			switch a.(type) {
			case types.MalBool:
				return types.MalBool(a.(types.MalBool) == b.(types.MalBool))
			case types.MalInt:
				return types.MalBool(a.(types.MalInt) == b.(types.MalInt))
			case types.MalSymbol:
				return types.MalBool(a.(types.MalSymbol) == b.(types.MalSymbol))
			case types.MalFunction:
				return types.MalBool(false)
			case types.MalList:
				aList, bList := a.(types.MalList), b.(types.MalList)
				if len(aList) != len(bList) {
					return types.MalBool(false)
				}
				for i := range aList {
					if !areEqual(aList[i], bList[i]) {
						return types.MalBool(false)
					}
				}
				return types.MalBool(true)
			}
			return types.MalBool(false)
		}

		return areEqual(a1, a2), nil
	}),
	"<": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a, b, err := getTwoMalInts(args)
		if err != nil {
			return nil, err
		}
		return types.MalBool(a < b), nil
	}),
	"<=": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a, b, err := getTwoMalInts(args)
		if err != nil {
			return nil, err
		}
		return types.MalBool(a <= b), nil
	}),
	">": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a, b, err := getTwoMalInts(args)
		if err != nil {
			return nil, err
		}
		return types.MalBool(a > b), nil
	}),
	">=": types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a, b, err := getTwoMalInts(args)
		if err != nil {
			return nil, err
		}
		return types.MalBool(a >= b), nil
	}),
}

func getTwoMalInts(args []types.MalType) (types.MalInt, types.MalInt, error) {
	if len(args) != 2 {
		return 0, 0, errors.New("expected 2 args")
	}
	a, b := args[0], args[1]

	if _, ok := a.(types.MalInt); !ok {
		return 0, 0, errors.New(printer.Print(a) + " is not an int")
	}
	if _, ok := b.(types.MalInt); !ok {
		return 0, 0, errors.New(printer.Print(b) + " is not an int")
	}
	return a.(types.MalInt), b.(types.MalInt), nil
}
