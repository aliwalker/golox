package lox

import "fmt"

func RuntimeError(token *Token, message string) {
	report(token.Line, " at '"+token.Lexeme+"'", "Runtime Error! "+message)
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

// Report to stdin.
func report(line int, where, message string) {
	fmt.Printf("[line %v] Error %v: %v\n", line, where, message)
}
