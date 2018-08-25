package golox

import (
	"testing"
)

func TestPrimaryString(t *testing.T) {
	// string
	scanner := NewScanner("\"This is a string\"")
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	exprs := parser.Parse()

	if len(exprs) != 1 {
		t.Error("expect len(exprs) to be 1")
	}

	literal, ok := exprs[0].(*Literal)
	if ok != true {
		t.Error("expect type 'Literal'")
	}

	if literal.Value != "This is a string" {
		t.Error("expect value to be \"This is a string\"")
	}
}

func TestPrimaryNumber(t *testing.T) {
}
