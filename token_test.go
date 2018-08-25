package golox

import (
	"fmt"
	"testing"
)

func TestNewToken(t *testing.T) {
	var token = NewToken(TokenAnd, "and", nil, 0)

	if token.Type != TokenAnd {
		t.Error(fmt.Sprintf("expecting type 'TokenAnd', but got: %v", token.Type))
	}

	if token.Lexeme != "and" {
		t.Error(fmt.Sprintf("expecting lexeme 'and', but got: %v", token.Lexeme))
	}

}
