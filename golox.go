package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aliwalker/golox/lox"
)

var interpreter = lox.NewInterpreter()

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [script]")
	} else if len(os.Args) == 2 {
		sourcePath, err := filepath.Abs(os.Args[1])
		if err != nil {
			fmt.Println("Unable to find path ")
			os.Exit(-1)
		}
		RunFile(sourcePath)
	} else {
		RunPrompt()
	}
}

func run(source string) (hadError, hadRuntimeError bool) {
	scanner := lox.NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := lox.NewParser(tokens)
	stmts, hadError := parser.Parse()

	if hadError {
		return
	}
	resolver := lox.NewResolver(interpreter)
	resolver.Resolve(stmts)

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
		fmt.Printf("Unable to read from file: %v.\n %v", path, err.Error())
		os.Exit(1)
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
