package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aliwalker/golox/lox"
)

var interpreter = lox.NewInterpreter()

func main() {
	run("fun foo(a1, a2) { print a1 + a2; } foo(1, 2);")
}

func run(source string) (hadError, hadRuntimeError bool) {
	scanner := lox.NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := lox.NewParser(tokens)
	stmts, hadError := parser.Parse()

	if hadError {
		return
	}
	hadRuntimeError = interpreter.Interprete(stmts)
	return
}

// RunFile runs a lox script file.
func RunFile(path string) {
	var (
		dat    []byte
		source string
		err    error
	)

	if dat, err = ioutil.ReadFile(path); err != nil {
		panic(fmt.Sprintf("Unable to read from file: %v.\n %v", path, err))
	}
	source = string(dat)
	hadError, hadRuntimeError := run(source)

	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(79)
	}
}

// RunPrompt provides a lox REPL environment.
func RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for true {
		fmt.Print("> ")
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("error reading from stdin.")
			os.Exit(80)
		}
		run(string(line))
	}
}
