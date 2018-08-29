package lox

import "fmt"

type LexingError struct {
	line    int
	message string
}

func NewLexingError(line int, message string) error {
	return &LexingError{line, message}
}

func (err *LexingError) Error() string {
	return fmt.Sprintf("[line %v] Error: %v\n", err.line, err.message)
}

// ParsingError occurs when there's syntax error.
type ParsingError struct {
	token   *Token
	message string
}

// NewParsingError returns a parsing error.
func NewParsingError(token *Token, message string) error {
	return &ParsingError{token, message}
}

// Error implements the built-in error interface.
func (err *ParsingError) Error() string {
	line := err.token.Line
	where := err.token.Lexeme
	message := err.message

	if err.token.Type == TokenEOF {
		where = "end"
	}
	return fmt.Sprintf("[line %v] Error at %v: %v\n", line, where, message)
}

type RuntimeError struct {
	token   *Token
	message string
}

func NewRuntimeError(token *Token, message string) error {
	return &RuntimeError{token, message}
}

func (err *RuntimeError) Error() string {
	line := err.token.Line
	where := err.token.Lexeme
	message := err.message

	if err.token.Type == TokenEOF {
		where = "end"
	}

	return fmt.Sprintf("[line %v] Runtime Error at %v: %v\n", line, where, message)
}
