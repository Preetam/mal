package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/Preetam/mal/dl/printer"

	"github.com/Preetam/mal/dl/reader"
	"github.com/Preetam/mal/dl/types"
)

type Environment map[string]types.MalType

func replRead(line string) (types.MalType, error) {
	return reader.ReadString(line)
}

func replEval(val types.MalType, env Environment) (types.MalType, error) {
	if list, ok := val.(types.MalList); ok {
		if len(list) == 0 {
			return val, nil
		}
		evaluated, err := evalAST(val, env)
		if err != nil {
			return nil, err
		}
		evaluatedList := evaluated.(types.MalList)
		f, ok := evaluatedList[0].(types.MalFunction)
		if !ok {
			return nil, errors.New("first argument is not a function")
		}
		return f(evaluatedList[1:]...)
	}
	return evalAST(val, env)
}

func replPrint(val types.MalType) string {
	return printer.Print(val)
}

func evalAST(ast types.MalType, env Environment) (types.MalType, error) {
	switch ast.(type) {
	case types.MalSymbol:
		symbolStr := string(ast.(types.MalSymbol))
		val, ok := env[symbolStr]
		if !ok {
			return nil, errors.New("unknown symbol " + symbolStr)
		}
		return val, nil
	case types.MalList:
		evaluatedList := types.MalList{}
		for _, v := range ast.(types.MalList) {
			evaluated, err := replEval(v, env)
			if err != nil {
				return nil, err
			}
			evaluatedList = append(evaluatedList, evaluated)
		}
		return evaluatedList, nil
	}
	return ast, nil
}

func rep(line string, env Environment) (string, error) {
	val, err := replRead(line)
	if err != nil {
		return "", err
	}
	val, err = replEval(val, env)
	if err != nil {
		return "", err
	}
	return replPrint(val), nil
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

var defaultEnv = Environment{
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
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("dl> ")
		if scanner.Scan() {
			resultString, err := rep(scanner.Text(), defaultEnv)
			if err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println(resultString)
			}
		} else {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
				os.Exit(1)
			}
			fmt.Println()
			return
		}
	}
}
