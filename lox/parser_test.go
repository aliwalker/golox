package lox

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

// args passed to check helpers.
// Operator's TokenType is always the first arg, if there is no operator(e.g, primary) pass 0.
// args for each Expr or Stmt are grouped. From root of an AST down to each leaves.
type argsTocheck []interface{}

func runChecks(t *testing.T, pairs map[string]argsTocheck) {
	for source, args := range pairs {
		stmts := parseSingleLine(t, source)
		checkExprAndPrintStmt(t, stmts, args...)
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
		expr = stmt.Expression
	case *Print:
		expr = stmt.Expression
	default:
		t.Error("expect Print or Expression stmt.")
	}

	operator, ok = values[0].(TokenType)
	if ok != true {
		operator = 0
	}
	values = values[1:]
	checkLogical(t, expr, operator, values...)
}

func checkLogical(t *testing.T, expr Expr, operator TokenType, values ...interface{}) []interface{} {
	var checkOperand = func(child *Logical, isLogical, isLeft bool) {
		var ok bool
		if isLogical == true && child != nil {
			if operator, ok = values[0].(TokenType); ok != true {
				t.Error("expect TokenType of an operator for Logical checking.")
			}
			values = values[1:] // consume 1 operator.

			checkLogical(t, child, operator, values...)
			values = values[2:] // consume 2 operands.
			return
		}

		logical, _ := expr.(*Logical)
		operator, ok = values[0].(TokenType)
		if ok == true {
			if isLeft == true {
				values = checkBinary(t, logical.Left, operator, values...)
			} else {
				values = checkBinary(t, logical.Right, operator, values...)
			}
			return
		}

		if isLeft == true {
			checkGroupingAndLiteral(t, logical.Left, values[0])
		} else {
			checkGroupingAndLiteral(t, logical.Right, values[0])
		}
		values = values[1:]
		//return
	}

	switch operator {
	case 0,
		TokenSlash,
		TokenStar,
		TokenPercent,
		TokenPlus,
		TokenMinus,
		TokenBang,
		TokenEqualEqual,
		TokenBangEqual,
		TokenGreater,
		TokenGreaterEqual,
		TokenLess,
		TokenLessEqual:
		checkBinary(t, expr, operator, values...)
	default:
		logical, ok := expr.(*Logical)
		if ok != true {
			t.Error("expect Logical expression.")
		}

		// check operator type.
		if logical.Operator.Type != operator {
			t.Error(fmt.Sprintf("checkLogical error: expect operator type %v, but got %v", operator, logical.Operator.Type))
		}
		left, ok := convertLogical(logical.Left)
		checkOperand(left, ok, true)

		right, ok := convertLogical(logical.Right)
		checkOperand(right, ok, false)
	}
	return values
}

func checkBinary(t *testing.T, expr Expr, operator TokenType, values ...interface{}) []interface{} {
	var checkOperand = func(child *Binary, isBinary bool, left bool) {
		var ok bool
		if isBinary == true && child != nil {
			if operator, ok = values[0].(TokenType); ok != true {
				t.Error("expect an operator's TokenType for Binary checking.")
			}
			// consume operator.
			values = values[1:]

			checkBinary(t, child, operator, values...)
			// consume 2 operands.
			values = values[2:]
			return
		}
		// only gonna have one expected value, ignore the rest if any.
		binary, _ := expr.(*Binary)
		if left == true {
			checkUnary(t, binary.Left, 0, values[0])
		} else {
			checkUnary(t, binary.Right, 0, values[0])
		}
		// consume 1 operands.
		values = values[1:]
		return
	}

	switch operator {
	case 0:
		//TokenBang,
		//TokenMinus:
		checkUnary(t, expr, operator, values[0])
	default:
		binary, ok := expr.(*Binary)
		if ok == true {
			if binary.Operator.Type != operator {
				t.Error(fmt.Sprintf("checkBinary error: expect operator type %v, but got %v", operator, binary.Operator.Type))
			}
			left, ok := convertBinary(binary.Left)
			checkOperand(left, ok, true)

			right, ok := convertBinary(binary.Right)
			checkOperand(right, ok, false)
			return values
		}

		unary, ok := expr.(*Unary)
		if ok == true {
			checkUnary(t, unary, operator, values[0])
			return values[1:]
		}

		t.Error("checkBinary error: expect Unary or Binary.")
	}
	return values
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
		t.Error(fmt.Sprintf("expect value to be %v, but got %v", value, literal.Value))
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
		"2 * 2;":      {TokenStar, float64(2), float64(2)},
		"2 * 3 * 4;":  {TokenStar, TokenStar, float64(2), float64(3), float64(4)},
		"10 / 2;":     {TokenSlash, float64(10), float64(2)},
		"10 / 2 / 5;": {TokenSlash, TokenSlash, float64(10), float64(2), float64(5)},
		"9 % 2;":      {TokenPercent, float64(9), float64(2)},
	}
	runChecks(t, runs)
}

func TestAddition(t *testing.T) {
	var runs = map[string]argsTocheck{
		"2 + 2;":     {TokenPlus, float64(2), float64(2)},
		"5+6;":       {TokenPlus, float64(5), float64(6)},
		"5 + 5 - 2;": {TokenMinus, TokenPlus, float64(5), float64(5), float64(2)},
	}
	runChecks(t, runs)
}

func TestComparison(t *testing.T) {
	var runs = map[string]argsTocheck{
		"1 < 2;":      {TokenLess, float64(1), float64(2)},
		"1 > 2 == 2;": {TokenEqualEqual, TokenGreater, float64(1), float64(2), float64(2)},
		"2 <= 3;":     {TokenLessEqual, float64(2), float64(3)},
		"9 > 5;":      {TokenGreater, float64(9), float64(5)},
		"9 >= 5;":     {TokenGreaterEqual, float64(9), float64(5)},
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
		"true and true;":          {TokenAnd, true, true},
		"false or true;":          {TokenOr, false, true},
		"true and true or false;": {TokenOr, TokenAnd, true, true, false},
	}
	runChecks(t, runs)
}

func TestPrecedence(t *testing.T) {
	var runs = map[string]argsTocheck{
		"1 + 2 * 3;":     {TokenPlus, float64(1), TokenStar, float64(2), float64(3)},
		"1 * 2 / 2;":     {TokenSlash, TokenStar, float64(1), float64(2), float64(2)},
		"1 + 2 < 2;":     {TokenLess, TokenPlus, float64(1), float64(2), float64(2)},
		"1 < 2 == true;": {TokenEqualEqual, TokenLess, float64(1), float64(2), true},
	}
	runChecks(t, runs)
}
