package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/Preetam/mal/dl/core"
	"github.com/Preetam/mal/dl/environment"
	"github.com/Preetam/mal/dl/printer"
	"github.com/Preetam/mal/dl/reader"
	"github.com/Preetam/mal/dl/types"
)

func replRead(line string) (types.MalType, error) {
	return reader.ReadString(line)
}

func replEval(val types.MalType, env *environment.Env) (types.MalType, error) {
	if _, ok := val.(types.MalList); !ok {
		// Not a list
		return evalAST(val, env)
	}
	list := val.(types.MalList)
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
			newEnv, _ := environment.New(env, nil, nil)
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
		case "do":
			var returnVal types.MalType
			var err error
			for _, elem := range list[1:] {
				returnVal, err = evalAST(elem, env)
				if err != nil {
					return nil, err
				}
			}
			return returnVal, nil
		case "if":
			result, err := replEval(list[1], env)
			if err != nil {
				return nil, err
			}
			boolResult, ok := result.(types.MalBool)
			if result == nil || (ok && !bool(boolResult)) {
				return replEval(list[2], env)
			}
			if len(list) < 4 {
				return nil, nil
			}
			return replEval(list[3], env)
		case "fn*":
			return types.MalFunction(func(args ...types.MalType) (types.MalType, error) {
				newEnv, err := environment.New(env, list[1], types.MalList(args))
				if err != nil {
					return nil, err
				}
				return replEval(list[2], newEnv)
			}), nil
		}
	}

	evaluated, err := evalAST(val, env)
	if err != nil {
		return nil, err
	}
	evaluatedList := evaluated.(types.MalList)
	f, ok := evaluatedList[0].(types.MalFunction)
	if !ok {
		return evaluatedList, nil
		//return nil, errors.New("first argument is not a function")
	}
	return f(evaluatedList[1:]...)
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

func defaultEnv() *environment.Env {
	env, _ := environment.New(nil, nil, nil)
	for name, f := range core.NS {
		env.Set(types.MalSymbol(name), f)
	}
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
