package lox

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	hadError        = false // Lexing or parsing error.
	hadRuntimeError = false // Runtime error.
)

var interpreter = NewInterpreter()

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	stmts := parser.Parse()

	if hadError {
		return
	}
	interpreter.Interprete(stmts)
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
	run(source)
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
		hadError = false
	}
}

// ParsingError reports a parsing error.
func ParsingError(token *Token, message string) {
	if token.Type == TokenEOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

// LexingError reports a lexing error.
func LexingError(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Printf("[line %v] Error %v: %v\n", line, where, message)
	hadError = true
}
