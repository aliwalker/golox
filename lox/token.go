package lox

import (
	"fmt"
)

// TokenType represents the type of a token.
type TokenType int

// Token types.
const (
	NotAKeyword TokenType = iota
	// single char
	TokenLeftParen
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenComma
	TokenDot
	TokenMinus
	TokenPlus
	TokenPercent
	TokenPrint
	TokenSemi
	TokenSlash
	TokenStar

	// one or two
	TokenBang
	TokenBangEqual
	TokenEqual
	TokenEqualEqual
	TokenGreater
	TokenGreaterEqual
	TokenLess
	TokenLessEqual

	// literal
	TokenString
	TokenIdentifier
	TokenNumber

	// keywords
	TokenAnd
	TokenClass
	TokenFalse
	TokenElse
	TokenFor
	TokenFun
	TokenIf
	TokenNil
	TokenOr
	TokenReturn
	TokenSuper
	TokenThis
	TokenTrue
	TokenVar
	TokenWhile

	TokenEOF
)

// Token represents a single unit.
type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

// NewToken creates a token structure.
func NewToken(t TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{t, lexeme, literal, line}
}

func (t Token) String() string {
	return fmt.Sprintf("%v %v %v", t.Type, t.Lexeme, t.Literal)
}
