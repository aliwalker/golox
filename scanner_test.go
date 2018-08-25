package golox

import (
	"fmt"
	"testing"
)

func errmsg(expectedType TokenType, token *Token) string {
	return fmt.Sprintf("expect type %v, but got token: %v", expectedType, token.Lexeme)
}

func TestScanString(t *testing.T) {
	var scanner = NewScanner("\"This is a string\"")
	var tokens = scanner.ScanTokens()

	// TokenString, TokenEOF.
	if len(tokens) != 2 {
		t.Error(fmt.Sprintf("expect len(tokens) to be 1, but got: %v.", len(tokens)))
	}

	// scanner should strip the double quotes
	if tokens[0].Literal != "This is a string" {
		t.Error(fmt.Sprintf("expect string value \"This is a string\", but got: %q", tokens[0].Literal))
	}
}

func TestScanNumber(t *testing.T) {
	var scanner = NewScanner("62")
	var tokens = scanner.ScanTokens()

	if len(tokens) != 2 {
		t.Error(fmt.Sprintf("expect len(tokens) to be 1, but got: %v.", len(tokens)))
	}

	if tokens[0].Literal != float64(62) {
		t.Error(fmt.Sprintf("expect number 62, but got: %v", tokens[0].Literal))
	}
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
