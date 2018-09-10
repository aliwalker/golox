package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aliwalker/golox/lox"
	"github.com/chzyer/readline"
)

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

func run(interpreter *lox.Interpreter, source string) (hadError, hadRuntimeError bool) {
	scanner := lox.NewScanner(source)
	tokens, hadError := scanner.ScanTokens()

	if hadError {
		return
	}

	parser := lox.NewParser(tokens)
	stmts, hadError := parser.Parse()

	if hadError {
		return
	}
	resolver := lox.NewResolver(interpreter)
	hadError = resolver.Resolve(stmts)

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

	interpreter := lox.NewInterpreter(false)

	if dat, err = ioutil.ReadFile(path); err != nil {
		fmt.Printf("Unable to read from file: %v.\n %v", path, err.Error())
		os.Exit(1)
	}
	source = string(dat)
	hadError, hadRuntimeError := run(interpreter, source)

	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(79)
	}
}

// RunPrompt provides a lox REPL environment.
func RunPrompt() {
	//reader := bufio.NewReader(os.Stdin)
	reader, err := readline.New("> ")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(80)
	}
	interpreter := lox.NewInterpreter(true)

	for {
		fmt.Print("> ")
		line, err := reader.Readline()
		if err != nil {
			// if the user press Ctrl + C or Ctrl + D, err will be io.EOF.
			if err == io.EOF {
				os.Exit(0)
			}
			// FIXME: I know this is ugly...
			if err.Error() == "Interrupt" {
				os.Exit(0)
			}
			fmt.Println("error reading from stdin.")
			os.Exit(80)
		}
		run(interpreter, string(line))
	}
}
