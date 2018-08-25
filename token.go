package golox

import (
	"fmt"
)

type TokenType int

// Token types.
const (
	NotAKeyword TokenType = iota
	// single char
	TokenLeftParen  //1
	TokenRightParen //2
	TokenLeftBrace  //3
	TokenRightBrace //4
	TokenComma      //5
	TokenDot        //6
	TokenMinus      //7
	TokenPlus       //8
	TokenSemi       //9
	TokenSlash      //10
	TokenStar       //11

	// one or two
	TokenBang         //12
	TokenBangEqual    //13
	TokenEqual        //14
	TokenEqualEqual   //15
	TokenGreater      //16
	TokenGreaterEqual //17
	TokenLess         //18
	TokenLessEqual    //19

	// literal
	TokenString     //20
	TokenIdentifier //21
	TokenNumber     //22

	// keywords
	TokenAnd    //23
	TokenClass  //24
	TokenFalse  //25
	TokenElse   //26
	TokenFor    //27
	TokenFun    //28
	TokenIf     //29
	TokenNil    //30
	TokenOr     //31
	TokenReturn //32
	TokenSuper  //33
	TokenThis   //34
	TokenTrue   //35
	TokenVar    //36
	TokenWhile  //37

	TokenEOF //38
)

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
