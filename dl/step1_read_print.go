package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Preetam/mal/dl/printer"

	"github.com/Preetam/mal/dl/reader"
	"github.com/Preetam/mal/dl/types"
)

func replRead(line string) (types.MalType, error) {
	return reader.ReadString(line)
}

func replEval(val types.MalType) types.MalType {
	return val
}

func replPrint(val types.MalType) string {
	return printer.Print(val)
}

func rep(line string) (string, error) {
	val, err := replRead(line)
	if err != nil {
		return "", err
	}
	val = replEval(val)
	return replPrint(val), nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("dl> ")
		if scanner.Scan() {
			resultString, err := rep(scanner.Text())
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
