package main

import (
	"bufio"
	"fmt"
	"os"
)

func replRead(line string) string {
	return line
}

func replEval(line string) string {
	return line
}

func replPrint(line string) string {
	return line
}

func rep(line string) string {
	return replPrint(replEval(replRead(line)))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("dl> ")
		if scanner.Scan() {
			fmt.Println(rep(scanner.Text()))
		} else {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
				return
			}
			fmt.Println()
			return
		}

	}
}
