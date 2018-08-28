package lox

import (
	"fmt"
	"testing"
)

func runExpr(t *testing.T, src string, expectedVal interface{}) {
	scanner := NewScanner(src)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	expr := parser.expression()
	interpreter := NewInterpreter()

	value := interpreter.evaluate(expr)

	if value != expectedVal {
		t.Error(fmt.Sprintf("expect i.evaluate(expr) to be %v, but got %v", expectedVal, value))
	}
}

func TestLiteralExpr(t *testing.T) {
	runExpr(t, "\"a test string.\"", "a test string.")
	runExpr(t, "5", float64(5))
}
func TestUnaryExpr(t *testing.T) {
	runExpr(t, "!true", false)
	runExpr(t, "-5", float64(-5))
}
func TestGroupingExpr(t *testing.T) {
	runExpr(t, "(1 + 2)", float64(3))
	runExpr(t, "-(1 + 2)", float64(-3))
}
func TestBinaryExpr(t *testing.T) {
	runExpr(t, "2 * 3 + 2", float64(8))
	runExpr(t, "1 + 2 / 2", float64(2))
	runExpr(t, "1 < 2", true)
	runExpr(t, "1 == 2", false)
	runExpr(t, "3 > 3", false)
	runExpr(t, "3 >= 3", true)
	runExpr(t, "1 != 2", true)
	runExpr(t, "1 - 2", float64(-1))
	runExpr(t, "5 % 2", 1)
}

func TestLogicalExpr(t *testing.T) {
	runExpr(t, "true and false", false)
	runExpr(t, "nil or 1", float64(1))
}
