package lox

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Scanner for lexing.
type Scanner struct {
	Tokens []*Token
	source string
	reader *strings.Reader

	current int
	start   int
	line    int

	hadError bool
}

var keywords = map[string]TokenType{
	"and":    TokenAnd,
	"break":  TokenBreak,
	"class":  TokenClass,
	"else":   TokenElse,
	"false":  TokenFalse,
	"for":    TokenFor,
	"fun":    TokenFun,
	"get":    TokenGetter,
	"if":     TokenIf,
	"nil":    TokenNil,
	"or":     TokenOr,
	"print":  TokenPrint,
	"return": TokenReturn,
	"super":  TokenSuper,
	"this":   TokenThis,
	"true":   TokenTrue,
	"var":    TokenVar,
	"while":  TokenWhile,
}

// NewScanner returns a new s.
func NewScanner(source string) *Scanner {
	return &Scanner{
		make([]*Token, 0),
		source,
		strings.NewReader(source), 0, 0, 1, false}
}

// ScanTokens returns a list of tokens from the source code.
func (s *Scanner) ScanTokens() ([]*Token, bool) {
	return s.scanTokens(), s.hadError
}

// helper function.
func (s *Scanner) scanTokens() []*Token {
	defer func() {
		if val := recover(); val != nil {
			s.hadError = true
			// repanic if it is not a LexingError.
			lexingError := val.(*LexingError)
			fmt.Println(lexingError.Error())
		}
	}()
	for !s.end() {
		s.start = s.current
		s.scanToken()
	}
	s.Tokens = append(s.Tokens, NewToken(TokenEOF, "", nil, s.line))
	return s.Tokens
}

func (s *Scanner) scanToken() {
	var c = s.advance()

	switch c {
	// single char token
	case '(':
		s.addToken(TokenLeftParen, nil)
	case ')':
		s.addToken(TokenRightParen, nil)
	case '{':
		s.addToken(TokenLeftBrace, nil)
	case '}':
		s.addToken(TokenRightBrace, nil)
	case '.':
		s.addToken(TokenDot, nil)
	case ',':
		s.addToken(TokenComma, nil)
	case ';':
		s.addToken(TokenSemi, nil)

	// one or two chars.
	case '-':
		if s.match('>') {
			s.addToken(TokenArrow, nil)
			return
		}
		s.addIfMatch('=', TokenMinusEqual, TokenMinus)
	case '+':
		s.addIfMatch('=', TokenPlusEqual, TokenPlus)
	case '*':
		s.addIfMatch('=', TokenStarEqual, TokenStar)
	case '%':
		s.addIfMatch('=', TokenPercentEqual, TokenPercent)
	case '!':
		s.addIfMatch('=', TokenBangEqual, TokenBang)
	case '=':
		s.addIfMatch('=', TokenEqualEqual, TokenEqual)
	case '<':
		s.addIfMatch('=', TokenLessEqual, TokenLess)
	case '>':
		s.addIfMatch('=', TokenGreaterEqual, TokenGreater)

	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.end() {
				s.advance()
			}
		} else {
			s.addIfMatch('=', TokenSlashEqual, TokenSlash)
		}

	// ignore white spaces.
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++

	// string
	case '"':
		s.string()

	default:
		if alpha(c) {
			// s.identifier covers keywords.
			s.identifier()
		} else if digit(c) {
			s.number()
		} else {
			panic(NewLexingError(s.line, "unexpected character "+string(s.peek())))
		}
	}
}

func (s *Scanner) advance() rune {
	// the caller of advance will ensure it is not at end.
	ch, sz, _ := s.reader.ReadRune()
	//if err != nil {
	//	panic(fmt.Sprintf("advance error: %v", err))
	//}
	s.current += sz
	return ch
}

func (s *Scanner) addToken(t TokenType, literal interface{}) {
	lexeme := s.source[s.start:s.current]
	s.Tokens = append(s.Tokens, NewToken(t, lexeme, literal, s.line))
}

func (s *Scanner) end() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) match(expected rune) bool {
	if s.end() {
		return false
	}

	ch, sz, _ := s.reader.ReadRune()
	// if err != nil {
	// 	panic(fmt.Sprintf("match error: %v", err))
	// }

	if ch != expected {
		if err := s.reader.UnreadRune(); err != nil {
			panic(fmt.Sprintf("match unread error: %v", err))
		}
		return false
	}

	s.current += sz
	return true
}

func (s *Scanner) addIfMatch(expected rune, expectedType TokenType, altType TokenType) {
	var (
		token     *Token
		tokenType TokenType
	)

	if s.peek() == expected {
		s.advance()
		tokenType = expectedType
	} else {
		tokenType = altType
	}

	token = NewToken(tokenType, s.source[s.start:s.current], nil, s.line)
	s.Tokens = append(s.Tokens, token)
}

func (s *Scanner) peek() rune {
	if s.end() {
		return 0
	}

	ch, _, err := s.reader.ReadRune()
	if err != nil {
		panic(fmt.Sprintf("peek error: %v", err))
	}

	if err := s.reader.UnreadRune(); err != nil {
		panic(fmt.Sprintf("peek unread error: %v", err))
	}
	return ch
}

func (s *Scanner) peekNext() rune {
	var (
		i   int
		err error
		ch  rune
	)

	for ; i < 2 && !s.end() && err == nil; i++ {
		ch, _, err = s.reader.ReadRune()
	}

	// unwind
	if _, err = s.reader.Seek(int64(s.current), io.SeekStart); err != nil {
		panic(fmt.Sprintf("peekNext unread error: %v", err))
	}

	// failed to read ahead 2 runes.
	if i < 2 {
		return 0
	}

	return ch
}

func (s *Scanner) identifier() {
	for !s.end() {
		if alphanumeric(s.peek()) {
			s.advance()
		} else {
			break
		}
	}

	lexeme := s.source[s.start:s.current]
	t := keywords[lexeme]

	if t == NotAKeyword {
		t = TokenIdentifier
	}

	s.Tokens = append(s.Tokens, NewToken(t, lexeme, nil, s.line))
}

func (s *Scanner) number() {
	var isInt = true

	for digit(s.peek()) {
		s.advance()
	}

	// check if it is a floating point number.
	if s.peek() == '.' && digit(s.peekNext()) {
		s.advance()
		isInt = false
		for digit(s.peek()) {
			s.advance()
		}
	}

	// if it is followed by non-whitespace, it is an error.
	if alphanumeric(s.peek()) {
		panic(NewLexingError(s.line, "identifier must start with a letter or underscore."))
	}

	value, err := strconv.ParseFloat(
		s.source[s.start:s.current],
		32)

	if err != nil {
		panic(NewLexingError(s.line, "error parsing number."))
	}

	// If it is an integer, keep the internal representation as integer at scanning.
	// It helps when interpreting modulo expression, e.g. 5 % 2, which requires both operands be integers.
	if isInt == true {
		s.addToken(TokenNumber, int(value))
	} else {
		s.addToken(TokenNumber, value)
	}
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.end() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.end() {
		panic(NewLexingError(s.line, "unterminated string."))
	}

	s.advance()
	s.addToken(TokenString, s.source[s.start+1:s.current-1])
}

func alpha(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch == '_')
}

func alphanumeric(ch rune) bool {
	return digit(ch) || alpha(ch)
}

func digit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func whitespace(ch rune) bool {
	switch ch {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}
