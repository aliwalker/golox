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

// LoxError occurs when there's syntax error.
type LoxError struct {
	token   *Token
	message string
}

// NewLoxError returns a parsing error.
func NewLoxError(token *Token, message string) error {
	return &LoxError{token, message}
}

// Error implements the built-in error interface.
func (err *LoxError) Error() string {
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
