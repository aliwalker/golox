package lox

import "fmt"

// LexingError represents error in lexing phase.
type LexingError struct {
	line    int
	message string
}

// NewLexingError returns a new lexing error.
func NewLexingError(line int, message string) error {
	return &LexingError{line, message}
}

// error interface.
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

// RuntimeError represents Lox runtime error.
type RuntimeError struct {
	token   *Token
	message string
}

// NewRuntimeError is a constructor.
func NewRuntimeError(token *Token, message string) error {
	return &RuntimeError{token, message}
}

// error interface.
func (err *RuntimeError) Error() string {
	line := err.token.Line
	where := err.token.Lexeme
	message := err.message

	if err.token.Type == TokenEOF {
		where = "end"
	}

	return fmt.Sprintf("[line %v] Runtime Error at %v: %v\n", line, where, message)
}
