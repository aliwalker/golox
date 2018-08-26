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

func parseSingleLine(t *testing.T, source string) []Stmt {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	stmts := parser.Parse()

	if len(stmts) != 1 {
		t.Error("expect len(exprs) to be 1")
	}
	return stmts
}

type argsTocheck []interface{}

func runChecks(t *testing.T, pairs map[string]argsTocheck) {
	for source, args := range pairs {
		stmts := parseSingleLine(t, source)
		checkExprAndPrintStmt(t, stmts, args)
	}
}

// ===================================== check helpers. =========================================

func checkExprAndPrintStmt(t *testing.T, stmts []Stmt, values ...interface{}) {
	var expr Expr
	var ok bool
	var operator TokenType

	// convert it to either an Expression Stmt or Print Stmt.
	switch stmt := stmts[0].(type) {
	case *Expression:
	case *Print:
		expr = stmt.Expression
		operator, ok = values[0].(TokenType)
		if ok != true {
			t.Error("checkExprAndPrintStmt error: expect first arg to be TokenType of an operator.")
		} else {
			checkLogical(t, expr, operator, values)
		}
	default:
		t.Error("expect Print or Expression stmt.")
	}
}

func checkLogical(t *testing.T, expr Expr, operator TokenType, values ...interface{}) {
	var checkOperand = func(child *Logical, isLogical bool) {
		var ok bool
		if isLogical == true && child != nil {
			if operator, ok = values[0].(TokenType); ok != true {
				t.Error("expect TokenType of an operator for Logical checking.")
			}
			values = values[1:]
			checkLogical(t, child, operator, values)
			return
		} else {
			checkBinary(t, expr, operator, values)
			return
		}
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
	var NoOperator = 0
	var runs = map[string]argsTocheck{
		"\"This is a string\";": {NoOperator, "This is a string"},
		"60;":    {NoOperator, float64(60)},
		"nil;":   {NoOperator, nil},
		"true;":  {NoOperator, true},
		"false;": {NoOperator, false},
		"(60);":  {NoOperator, float64(60)},
	}

	runChecks(t, runs)
}

func TestUnary(t *testing.T) {
	var runs = map[string]argsTocheck{
		"!true;": {TokenBang, true},
		"-1;":    {TokenMinus, float64(1)},
	}

	runChecks(t, runs)
}

func TestMultiplication(t *testing.T) {
	var runs = map[string]argsTocheck{
		"2 * 2;":  {TokenStar, float64(2), float64(2)},
		"10 / 2;": {TokenSlash, float64(10), float64(2)},
		"9 % 2;":  {TokenPercent, float64(9), float64(2)},
	}
	runChecks(t, runs)
}

func TestAddition(t *testing.T) {
	var runs = map[string]argsTocheck{
		"2 + 2;": {TokenPlus, float64(2), float64(2)},
		"5+6;":   {TokenPlus, float64(5), float64(6)},
	}
	runChecks(t, runs)
}

func TestComparison(t *testing.T) {
	var runs = map[string]argsTocheck{
		"1 < 2;":  {TokenLess, float64(1), float64(2)},
		"2 <= 3;": {TokenLessEqual, float64(2), float64(3)},
		"9 > 5;":  {TokenGreater, float64(9), float64(5)},
		"9 >= 5;": {TokenGreaterEqual, float64(9), float64(5)},
	}
	runChecks(t, runs)
}

func TestEquality(t *testing.T) {
	var runs = map[string]argsTocheck{
		"1 == 2;": {TokenEqualEqual, float64(1), float64(2)},
		"1 != 2;": {TokenBangEqual, float64(1), float64(2)},
	}

	runChecks(t, runs)
}

func TestLogical(t *testing.T) {
	var runs = map[string]argsTocheck{
		"true and true;": {TokenAnd, true, true},
		"false or true;": {TokenOr, false, true},
	}
	runChecks(t, runs)
}
