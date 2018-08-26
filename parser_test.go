package golox

import (
	"fmt"
	"testing"
)

func convertBinary(expr Expr) (*Binary, bool) {
	binary, ok := expr.(*Binary)
	if ok == true {
		return binary, ok
	}
	return nil, false
}

func convertLogical(expr Expr) (*Logical, bool) {
	logical, ok := expr.(*Logical)
	if ok == true {
		return logical, ok
	}
	return nil, false
}

func runParser(t *testing.T, source string) []Stmt {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	stmts := parser.Parse()

	if len(stmts) != 1 {
		t.Error("expect len(exprs) to be 1")
	}
	return stmts
}

// ===================================== check helpers. =========================================

func checkExprAndPrintStmt(t *testing.T, stmts []Stmt, operator TokenType, values ...interface{}) {
	var expr Expr

	// convert it to either an Expression Stmt or Print Stmt.
	switch stmt := stmts[0].(type) {
	case *Expression:
	case *Print:
		expr = stmt.Expression
		checkLogical(t, expr, operator, values)
	default:
		t.Error("expect Print or Expression stmt.")
	}
}

func checkLogical(t *testing.T, expr Expr, operator TokenType, values ...interface{}) {
	var checkOperand = func(child *Logical, isLogical bool) {
		var ok bool
		if isLogical == true && child != nil {
			if operator, ok = values[0].(TokenType); ok != true {
				t.Error("expect an operator's TokenType for Logical checking.")
			}

			values = values[1:]
			checkLogical(t, child, operator, values)
			return
		} else {
			checkBinary(t, expr, operator, values)
			return
		}

		t.Error("checkLogical error. ")
	}

	switch operator {
	case 0:
	case TokenSlash:
	case TokenStar:
	case TokenPercent:
	case TokenPlus:
	case TokenMinus:
	case TokenBang:
	case TokenEqualEqual:
	case TokenBangEqual:
	case TokenGreater:
	case TokenGreaterEqual:
	case TokenLess:
	case TokenLessEqual:
		checkBinary(t, expr, operator, values)
	default:
		logical, ok := expr.(*Logical)
		if ok != true {
			t.Error("expect Logical expression.")
		}

		// check operator type.
		if logical.Operator.Type != operator {
			t.Error(fmt.Sprintf("checkLogical error: expect operator type %v, but got %v", operator, logical.Operator.Type))
		}
		checkOperand(convertLogical(logical.Left))
		checkOperand(convertLogical(logical.Right))
	}
}

func checkBinary(t *testing.T, expr Expr, operator TokenType, values ...interface{}) {
	var checkOperand = func(child *Binary, isLogical bool) {
		var ok bool
		if isLogical == true && child != nil {
			if operator, ok = values[0].(TokenType); ok != true {
				t.Error("expect an operator's TokenType for Binary checking.")
			}

			values = values[1:]
			checkBinary(t, child, operator, values)
			return
		} else {
			// only gonna have one expected value, ignore the rest if any.
			checkUnary(t, expr, operator, values[0])
			return
		}

		t.Error("checkBinary error. ")
	}

	switch operator {
	case 0:
	case TokenBang:
	case TokenMinus:
		checkUnary(t, expr, operator, values[0])
	default:
		binary, ok := expr.(*Binary)
		if ok != true {
			t.Error("expect Binary expression.")
		}

		if binary.Operator.Type != operator {
			t.Error(fmt.Sprintf("checkBinary error: expect operator type %v, but got %v", operator, binary.Operator.Type))
		}
		checkOperand(convertBinary(binary.Left))
		checkOperand(convertBinary(binary.Right))
	}
}

func checkUnary(t *testing.T, expr Expr, operator TokenType, value interface{}) {
	switch operator {
	case 0:
		checkGroupingAndLiteral(t, expr, value)
	default:
		unary, ok := expr.(*Unary)
		if ok != true {
			t.Error("expect Unary expression.")
		}

		if unary.Operator.Type != operator {
			t.Error(fmt.Sprintf("checkUnary error: expect operator type %v, but got %v", operator, unary.Operator.Type))
		}
		checkGroupingAndLiteral(t, unary.Right, value)
	}
}

// p.primary only returns Grouping & Literal.
func checkGroupingAndLiteral(t *testing.T, expr Expr, value interface{}) {
	var literal *Literal
	var ok bool

	// convert it to either a Literal or Grouping Expr.
	switch v := expr.(type) {
	case *Literal:
		literal = v
	case *Grouping:
		innerExpr := v.Expression
		if literal, ok = innerExpr.(*Literal); ok != true {
			t.Error("expect Expression in Grouping.")
		}
	default:
		t.Error("expect Grouping or Literal expr.")
	}

	if literal.Value != value {
		t.Error(fmt.Sprintf("expect value to be %q", value))
	}
}

func TestPrimary(t *testing.T) {
	stmts := runParser(t, "\"This is a string\";")
	checkExprAndPrintStmt(t, stmts, 0, "This is a string;")

	stmts = runParser(t, "60;")
	checkExprAndPrintStmt(t, stmts, 0, float64(60))

	stmts = runParser(t, "nil;")
	checkExprAndPrintStmt(t, stmts, 0, nil)

	stmts = runParser(t, "true;")
	checkExprAndPrintStmt(t, stmts, 0, true)

	stmts = runParser(t, "false;")
	checkExprAndPrintStmt(t, stmts, 0, false)

	// Grouping.
	stmts = runParser(t, "(60);")
	checkExprAndPrintStmt(t, stmts, 0, float64(60))
}

func TestUnary(t *testing.T) {
	var stmts = runParser(t, "!true;")
	checkExprAndPrintStmt(t, stmts, TokenBang, true)

	stmts = runParser(t, "-1;")
	checkExprAndPrintStmt(t, stmts, TokenMinus, float64(1))
}

func TestMultiplication(t *testing.T) {
	var stmts = runParser(t, "2 * 2;")
	checkExprAndPrintStmt(t, stmts, TokenStar, float64(2), float64(2))

	stmts = runParser(t, "10 / 2;")
	checkExprAndPrintStmt(t, stmts, TokenSlash, float64(10), float64(2))

	stmts = runParser(t, "9 % 2;")
	checkExprAndPrintStmt(t, stmts, TokenPercent, float64(9), float64(2))
}

func TestAddition(t *testing.T) {
	var stmts = runParser(t, "2 + 2;")
	checkExprAndPrintStmt(t, stmts, TokenPlus, float64(2), float64(2))

	stmts = runParser(t, "5+6;")
	checkExprAndPrintStmt(t, stmts, TokenPlus, float64(5), float64(6))
}

func TestComparison(t *testing.T) {
	var stmts = runParser(t, "1 < 2;")
	checkExprAndPrintStmt(t, stmts, TokenLess, float64(1), float64(2))

	stmts = runParser(t, "2 <= 3;")
	checkExprAndPrintStmt(t, stmts, TokenLessEqual, float64(2), float64(3))

	stmts = runParser(t, "9 > 5;")
	checkExprAndPrintStmt(t, stmts, TokenGreater, float64(9), float64(5))

	stmts = runParser(t, "9 >= 5;")
	checkExprAndPrintStmt(t, stmts, TokenGreaterEqual, float64(9), float64(5))
}

func TestEquality(t *testing.T) {
	stmts := runParser(t, "1 == 2;")
	// I know the func name is inappropriate.
	checkExprAndPrintStmt(t, stmts, TokenEqualEqual, float64(1), float64(2))

	stmts = runParser(t, "1 != 2;")
	checkExprAndPrintStmt(t, stmts, TokenBangEqual, float64(1), float64(2))
}

func TestLogical(t *testing.T) {
	var stmts = runParser(t, "true and true;")
	checkExprAndPrintStmt(t, stmts, TokenAnd, true, true)

	stmts = runParser(t, "false or true;")
	checkExprAndPrintStmt(t, stmts, TokenOr, false, true)
}
