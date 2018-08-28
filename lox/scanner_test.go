package lox

import (
	"fmt"
	"math"
	"testing"
)

func errmsg(expectedType TokenType, token *Token) string {
	return fmt.Sprintf("expect type %v, but got token: %v", expectedType, token.Lexeme)
}

func checkLiteralToken(t *testing.T, src string, tokenType TokenType, val interface{}) {
	scanner := NewScanner(src)
	tokens := scanner.ScanTokens()

	if len(tokens) != 2 {
		t.Error(fmt.Sprintf("expect len(tokens) to be 2, but got %v", len(tokens)))
	}

	if tokens[0].Type != tokenType {
		t.Error(fmt.Sprintf("expect TokenType %v, but got %v", tokenType, tokens[0].Type))
	}

	if tokens[0].Type == TokenNumber {
		switch v := tokens[0].Literal.(type) {
		case float64:
			// suppose we only test floating number with precision 2.
			tokens[0].Literal = math.Floor(v*100) / 100
		}
	}
	if tokens[0].Literal != val {
		t.Error(fmt.Sprintf("expect Literal %v, but got %v", val, tokens[0].Literal))
	}
}

func TestScanString(t *testing.T) {
	checkLiteralToken(t, "\"This is a string\"", TokenString, "This is a string")
	checkLiteralToken(t, "\"First line.\nSecond line.\"", TokenString, "First line.\nSecond line.")
}

func TestScanNumber(t *testing.T) {
	checkLiteralToken(t, "62", TokenNumber, 62)
	checkLiteralToken(t, "62.22", TokenNumber, 62.22)
}

func TestScanSingleLine(t *testing.T) {
	var scanner = NewScanner("var a = 1;")
	var expectedTypes = []TokenType{TokenVar, TokenIdentifier, TokenEqual, TokenNumber, TokenSemi, TokenEOF}
	var tokens = scanner.ScanTokens()

	if len(tokens) != len(expectedTypes) {
		t.Error(fmt.Sprintf("expect %v tokens, but got: %v", len(expectedTypes), len(tokens)))
	}

	for i, token := range tokens {
		if token.Type != expectedTypes[i] {
			t.Error(errmsg(expectedTypes[i], token))
		}
	}
}

func TestScanMultilines(t *testing.T) {
	var scanner = NewScanner("var a;\nvar b;")
	var tokens = scanner.ScanTokens()
	var types = []TokenType{TokenVar, TokenIdentifier, TokenSemi, TokenVar, TokenIdentifier, TokenSemi, TokenEOF}

	if len(tokens) != len(types) {
		t.Error(fmt.Sprintf("expect %v tokens, but got: %v", len(types), len(tokens)))
	}

	for i, typ := range types {
		if typ != tokens[i].Type {
			t.Error(errmsg(typ, tokens[i]))
		}
	}

	if tokens[3].Line != 2 {
		t.Error(fmt.Sprintf("expect line 2, but got: %v", tokens[3].Line))
	}
}

func TestIdentifierName(t *testing.T) {
	var scanner = NewScanner("varidentifier")
	var tokens = scanner.ScanTokens()

	if tokens[0].Type != TokenIdentifier {
		t.Error(errmsg(TokenIdentifier, tokens[0]))
	}
}
