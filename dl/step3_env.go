package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/Preetam/mal/dl/environment"
	"github.com/Preetam/mal/dl/printer"
	"github.com/Preetam/mal/dl/reader"
	"github.com/Preetam/mal/dl/types"
)

func replRead(line string) (types.MalType, error) {
	return reader.ReadString(line)
}

func replEval(val types.MalType, env *environment.Env) (types.MalType, error) {
	if list, ok := val.(types.MalList); ok {
		if len(list) == 0 {
			return val, nil
		}

		if symbol, ok := list[0].(types.MalSymbol); ok {
			switch symbol {
			case "def!":
				if len(list) < 3 {
					return nil, errors.New("missing def! arguments")
				}
				if keySymbol, ok := list[1].(types.MalSymbol); ok {
					evaluated, err := replEval(list[2], env)
					if err != nil {
						return nil, err
					}
					env.Set(keySymbol, evaluated)
					return evaluated, nil
				} else {
					return nil, errors.New("expected symbol argument")
				}
			case "let*":
				newEnv := environment.New(env)
				if len(list) < 3 {
					return nil, errors.New("missing let* arguments")
				}
				bindings, ok := list[1].(types.MalList)
				if !ok {
					return nil, errors.New("expected a list of bindings")
				}
				if len(bindings)%2 != 0 {
					return nil, errors.New("incomplete bindings list")
				}
				for i := 0; i < len(bindings); i += 2 {
					keySymbol, ok := bindings[i].(types.MalSymbol)
					if !ok {
						return nil, errors.New("expected symbol key")
					}
					evaluated, err := replEval(bindings[i+1], newEnv)
					if err != nil {
						return nil, err
					}
					newEnv.Set(keySymbol, evaluated)
				}
				return replEval(list[2], newEnv)
			}
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

func evalAST(ast types.MalType, env *environment.Env) (types.MalType, error) {
	switch ast.(type) {
	case types.MalSymbol:
		val, err := env.Get(ast.(types.MalSymbol))
		if err != nil {
			return nil, errors.New(string(ast.(types.MalSymbol)) + " not found")
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

func rep(line string, env *environment.Env) (string, error) {
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

func defaultEnv() *environment.Env {
	env := environment.New(nil)
	env.Set("+", types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a, b, err := getTwoMalInts(args)
		if err != nil {
			return nil, err
		}
		return a + b, nil
	}))
	env.Set("-", types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a, b, err := getTwoMalInts(args)
		if err != nil {
			return nil, err
		}
		return a - b, nil
	}))
	env.Set("*", types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a, b, err := getTwoMalInts(args)
		if err != nil {
			return nil, err
		}
		return a * b, nil
	}))
	env.Set("/", types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
		a, b, err := getTwoMalInts(args)
		if err != nil {
			return nil, err
		}
		return a / b, nil
	}))
	return env
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	env := defaultEnv()
	for {
		fmt.Print("dl> ")
		if scanner.Scan() {
			resultString, err := rep(scanner.Text(), env)
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
