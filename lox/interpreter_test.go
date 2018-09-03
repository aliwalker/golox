package lox

import (
	"fmt"
	"math"
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
		// I know this is ugly...
		if v1, ok1 := value.(float64); ok1 {
			if v2, ok2 := expectedVal.(float64); ok2 {
				if math.Floor(v1*100)/100 == math.Floor(v2*100)/100 {
					return
				}
			}
		}
		t.Error(fmt.Sprintf("expect i.evaluate(expr) to be %v, but got %v", expectedVal, value))
	}
}

func runStmt(t *testing.T, src string) {
	scanner := NewScanner(src)
	parser := NewParser(scanner.ScanTokens())
	stmts, _ := parser.Parse()
	interpreter := NewInterpreter()

	if parser.hadError == true {
		t.Error("syntax error.")
	}
	interpreter.Interprete(stmts)

	if interpreter.hadRuntimeError != false {
		t.Error("runtime error.")
	}
}

func runErrStmt(t *testing.T, src string) {
	scanner := NewScanner(src)
	parser := NewParser(scanner.ScanTokens())
	stmts, hadError := parser.Parse()

	if hadError {
		return
	}
	interpreter := NewInterpreter()
	interpreter.Interprete(stmts)

	if interpreter.hadRuntimeError != true {
		t.Error("expect runtime error.")
	}
}

func TestRunStmt(t *testing.T) {
	runStmt(t, "123;")
	runStmt(t, "var a; a = 2;") // test var & assign stmts.
	runStmt(t, "{ var a = 1; print a; }")
	runStmt(t, "for (var i = 0; i < 3; i = i + 1) { print i; }")
}

func TestBreakStmt(t *testing.T) {
	runStmt(t, "var a = 0;\n while (true) {\n print a; a += 1; \nif (a == 3) \nbreak; \n}")
}

func TestFuncStmt(t *testing.T) {
	runStmt(t, "fun foo() { print \"bar\"; }")
	runStmt(t, "fun bar(foobar) { print foobar;  }")
	runStmt(t, "fun foo() { print \"bar\"; } \nfoo();")
	runStmt(t, "fun foo() { fun bar() { print \"ok\"; } bar(); }")
}
func TestIfStmt(t *testing.T) {
	runStmt(t, "if (true) print \"true\";")
	runStmt(t, "var a = 5; var b = 4; if (a > b) print a; else print b;")
	runStmt(t, "if (5 > 4) { print \"a\"; print \"b\"; }")
}
func TestPrintStmt(t *testing.T) {
	runStmt(t, "print \"hello\";")
}
func TestVarStmt(t *testing.T) {
	runStmt(t, "var a = 1;")
	runStmt(t, "var a = 1; a;")
}
func TestWhileStmt(t *testing.T) {
	runStmt(t, "var i = 0;\nwhile (i < 3) { \nprint i; \ni = i + 1; \n}")
	runStmt(t, "var i = 0; while (i < 3) i = i + 1;")
}
func TestLiteralExpr(t *testing.T) {
	runExpr(t, "\"a test string.\"", "a test string.")
	runExpr(t, "5", 5)
}
func TestUnaryExpr(t *testing.T) {
	runExpr(t, "!true", false)
	runExpr(t, "-5", -5)
	runExpr(t, "-1.1", -1.1)
	runExpr(t, "!1", false)
}
func TestGroupingExpr(t *testing.T) {
	runExpr(t, "(1 + 2)", 3)
	runExpr(t, "-(1 + 2)", -3)
}
func TestBinaryExpr(t *testing.T) {
	runExpr(t, "2 * 3 + 2", 8)
	runExpr(t, "1.0 + 2.0", 3.0)
	runExpr(t, "2.1 + 2.5", 4.6)
	runExpr(t, "1.1 - 0.1", 1.0)
	runExpr(t, "2.0 * 3.0", 6.0)
	runExpr(t, "2.2 / 2.0", 1.1)
	runExpr(t, "1 + 2 / 2", 2)
	runExpr(t, "1 < 2", true)
	runExpr(t, "1 <= 1", true)
	runExpr(t, "1 == 2", false)
	runExpr(t, "1.1 == 1.1", true)
	runExpr(t, "3 > 3", false)
	runExpr(t, "3 >= 3", true)
	runExpr(t, "1 != 2", true)
	runExpr(t, "1 - 2", -1)
	runExpr(t, "5 % 2", 1)
	runExpr(t, "\"adorable\" + \" lady\"", "adorable lady")
}

func TestLogicalExpr(t *testing.T) {
	runExpr(t, "true and false", false)
	runExpr(t, "false and true", false)
	runExpr(t, "nil or 1", 1)
}

func TestAssignExpr(t *testing.T) {
	runStmt(t, "var a = 1;\na += 2; \nprint a;")
	runStmt(t, "var a = \"head\"; a += \" tail\"; print a;")
	runStmt(t, "var b = 1; b -= 1; print b;")
}

func TestRuntimeError(t *testing.T) {
	runErrStmt(t, "true - true;")
	runErrStmt(t, "1 + \"a string\";")
	runErrStmt(t, "-\"a string\";")
	runErrStmt(t, "5.0 % 2;")
	runErrStmt(t, "\"a string\" % \"another string\";")
	runErrStmt(t, "a;")     // due to variable `a`` is not defined.
	runErrStmt(t, "a = 1;") // due to variable `a` is not defined.
	runErrStmt(t, "var 1err = 123; 123;")
	runErrStmt(t, "var a = 1; {\n\tvar b = 2;\n}\n\tprint b;")
	runErrStmt(t, "fun foo() { print foo }")
	runErrStmt(t, "fun foo(1) { print foo; }")
	runErrStmt(t, "fun foo(a1, a2, a3, a4, a5, a6, a7, a8, a9){}")
	runErrStmt(t, "\"notAFun\"();")
	runErrStmt(t, "45();")
	runErrStmt(t, "foo();")
	runErrStmt(t, "fun foo(a1, a2) { print a1 + a2; } foo(1, 2, 3)")
	runErrStmt(t, "else print 5;")
}
