package lox

import (
	"fmt"
	"math"
	"testing"
)

func runExpr(t *testing.T, src string, expectedVal interface{}) {
	scanner := NewScanner(src)
	tokens, hadError := scanner.ScanTokens()

	if hadError {
		t.Error("scanning error.")
	}
	parser := NewParser(tokens)
	expr := parser.expression()
	interpreter := NewInterpreter(false)

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
	tokens, hadError := scanner.ScanTokens()
	if hadError {
		t.Error("scanning error.")
	}

	parser := NewParser(tokens)
	stmts, hadError := parser.Parse()

	if hadError == true {
		t.Error("syntax error.")
	}

	interpreter := NewInterpreter(false)
	resolver := NewResolver(interpreter)
	hadError = resolver.Resolve(stmts)

	if hadError == true {
		t.Error("resolve error.")
	}

	interpreter.Interprete(stmts)

	if interpreter.hadRuntimeError != false {
		t.Error("runtime error.")
	}
}

// ==================================== specific error runner ===================================
// These are runners for testing specific errors, `src` passed to them should be ensured to have
// specific errors, therefore some error checking are stripped.
func runSynErrStmt(t *testing.T, src string) {
	scanner := NewScanner(src)
	tokens, hadError := scanner.ScanTokens()
	if hadError {
		return
	}
	parser := NewParser(tokens)
	stmts, hadError := parser.Parse()

	if hadError {
		return
	}
	interpreter := NewInterpreter(false)
	interpreter.Interprete(stmts)

	if interpreter.hadRuntimeError != true {
		t.Error("expect runtime error.")
	}
}

func runResErrStmt(t *testing.T, src string) {
	scanner := NewScanner(src)
	tokens, _ := scanner.ScanTokens()
	parser := NewParser(tokens)
	stmts, _ := parser.Parse()

	interpreter := NewInterpreter(false)
	resolver := NewResolver(interpreter)

	resolver.Resolve(stmts)
	if resolver.hadError != true {
		t.Error("expect resolving error.")
	}
}

func runRuntimeErrStmt(t *testing.T, src string) {
	scanner := NewScanner(src)
	tokens, _ := scanner.ScanTokens()
	parser := NewParser(tokens)
	stmts, _ := parser.Parse()

	interpreter := NewInterpreter(false)
	//resolver.Resolve(stmts)

	//resolver.Resolve(stmts)
	interpreter.Interprete(stmts)
	if interpreter.hadRuntimeError != true {
		t.Error("expect runtime error.")
	}
}

// ================================= Error handler testing ================================
func TestSynError(t *testing.T) {
	runSynErrStmt(t, "true - true;")
	runSynErrStmt(t, "1 + \"a string\";")
	runSynErrStmt(t, "-\"a string\";")
	runSynErrStmt(t, "5.0 % 2;")
	runSynErrStmt(t, "\"a string\" % \"another string\";")
	runSynErrStmt(t, "a;")     // variable `a`` is not defined.
	runSynErrStmt(t, "a = 1;") // variable `a` is not defined.
	runSynErrStmt(t, "var 1err = 123; 123;")
	runSynErrStmt(t, "var a = 1; {\n\tvar b = 2;\n}\n\tprint b;")
	runSynErrStmt(t, "fun foo(1) { print foo; }")
	runSynErrStmt(t, "fun foo(a1, a2, a3, a4, a5, a6, a7, a8, a9){}")
	runSynErrStmt(t, "fun foo(a1, a2) { print a1 + a2; } foo(1, 2, 3)")
	runSynErrStmt(t, "else print 5;")
}
func TestRuntimeErrStmt(t *testing.T) {
	runRuntimeErrStmt(t, "\"notAFun\"();")
	runRuntimeErrStmt(t, "45();")
	runRuntimeErrStmt(t, "foo();")
	runRuntimeErrStmt(t, "var a, b = 1; a + b;")                                      // uninitialized variable a.
	runRuntimeErrStmt(t, "a + b;")                                                    // undefined variable a and b.
	runRuntimeErrStmt(t, "5.0 % 2;")                                                  // modulo arithmetic on floating point numbers.
	runRuntimeErrStmt(t, "2();")                                                      // call a non callable.
	runRuntimeErrStmt(t, "fun increment(arg1) { return arg1 + 1; } increment(1, 2);") // unmatched args and params.
}

func TestResError(t *testing.T) {
	runResErrStmt(t, "var a; break;")
	runResErrStmt(t, "return a;")
	runResErrStmt(t, "{var a = a + 1;}")
}

// ========================================================================================
func TestRunStmt(t *testing.T) {
	runStmt(t, "123;")
	runStmt(t, "var a; a = 2;") // test var & assign stmts.
	runStmt(t, "{ var a = 1; print a; }")
	runStmt(t, "for (var i = 0; i < 3; i = i + 1) { print i; }")
	runStmt(t, "fun foo() { print foo }")

}

func TestLambda(t *testing.T) {
	runStmt(t, "var a = () -> 1 + 2; print a()")
}

func TestClosure(t *testing.T) {
	runStmt(t, "fun foo() { fun bar() { print \"bar\" } return bar; } var bar = foo() bar()")
}
func TestBinding(t *testing.T) {
	runStmt(t, "var a = 0; { fun foo() { print a; } foo(); var a = 1; foo(); }")
	runStmt(t, "fun foo() { fun bar() { print \"bar\" } bar() }")
}

func TestVarListStmt(t *testing.T) {
	runStmt(t, "var a = 1, b = 1; print a + b;")
	runStmt(t, "var a, b, c; a = 1; b = 2; c = 3;")
}

func TestBreakStmt(t *testing.T) {
	runStmt(t, "var a = 0;\n while (true) {\n print a; a += 1; \nif (a == 3) \nbreak; \n}")
}

func TestFuncStmt(t *testing.T) {
	runStmt(t, "fun foo() { print \"bar\"; }")
	runStmt(t, "fun bar(foobar) { print foobar;  }")
	runStmt(t, "fun foo() { print \"bar\"; } \nfoo();")
	runStmt(t, "fun foo() { fun bar() { print \"ok\"; } bar(); }")
	runStmt(t, "fun foo() { return 5; } print foo();")
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
