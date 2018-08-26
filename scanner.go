package golox

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
}

var keywords = map[string]TokenType{
	"and":    TokenAnd,
	"class":  TokenClass,
	"else":   TokenElse,
	"false":  TokenFalse,
	"for":    TokenFor,
	"fun":    TokenFun,
	"if":     TokenIf,
	"nil":    TokenNil,
	"or":     TokenOr,
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
		strings.NewReader(source), 0, 0, 1}
}

// ScanTokens returns a list of tokens from the source code.
func (s *Scanner) ScanTokens() []*Token {
	for !s.isAtEnd() {
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
	case ',':
		s.addToken(TokenComma, nil)
	case ';':
		s.addToken(TokenSemi, nil)
	case '-':
		s.addToken(TokenMinus, nil)
	case '+':
		s.addToken(TokenPlus, nil)
	case '*':
		s.addToken(TokenStar, nil)
	case '%':
		s.addToken(TokenPercent, nil)

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
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(TokenSlash, nil)
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
		if isAlpha(c) {
			// s.identifier covers keywords.
			s.identifier()
		} else if isDigit(c) {
			s.number()
		} else {
			LexingError(s.line, "unexpected character.")
		}
	}
}

func (s *Scanner) advance() rune {
	ch, sz, err := s.reader.ReadRune()
	if err != nil {
		panic(fmt.Sprintf("advance error: %v", err))
	}
	s.current += sz
	return ch
}

func (s *Scanner) addToken(t TokenType, literal interface{}) {
	lexeme := s.source[s.start:s.current]
	s.Tokens = append(s.Tokens, NewToken(t, lexeme, literal, s.line))
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	ch, sz, err := s.reader.ReadRune()
	if err != nil {
		panic(fmt.Sprintf("match error: %v", err))
	}

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
	if s.isAtEnd() {
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

	for ; i < 2 && !s.isAtEnd() && err == nil; i++ {
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
	for !s.isAtEnd() {
		if isAlphanumeric(s.peek()) {
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
	for isDigit(s.peek()) {
		s.advance()
	}

	// check if it is a floating point number.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(
		s.source[s.start:s.current],
		32)

	if err != nil {
		LexingError(s.line, "error parsing number.")
		return
	}

	s.addToken(TokenNumber, value)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		LexingError(s.line, "unterminated string.")
		return
	}

	s.advance()
	s.addToken(TokenString, s.source[s.start+1:s.current-1])
}

func isAlpha(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch == '_')
}

func isAlphanumeric(ch rune) bool {
	return isDigit(ch) || isAlpha(ch)
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}
