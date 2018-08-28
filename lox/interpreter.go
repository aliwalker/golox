package lox

import (
	"fmt"
)

type Interpreter struct {
	hadRuntimeError bool
}

func NewInterpreter() *Interpreter {
	return &Interpreter{false}
}

func (i *Interpreter) Interprete(stmts []Stmt) (hadRuntimeError bool) {
	defer func() {
		if r := recover(); r != nil {
			i.hadRuntimeError = true
		}
		hadRuntimeError = i.hadRuntimeError
	}()

	for _, stmt := range stmts {
		i.execute(stmt)
	}
	return
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) VisitExpressionStmt(stmt *Expression) interface{} {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) interface{} {
	val := i.evaluate(stmt.Expression)
	fmt.Println(val)
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *Assign) interface{} {
	value := i.evaluate(expr.Value)
	return value
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	// comparison.
	case TokenBangEqual:
		return !isEqual(left, right)
	case TokenEqualEqual:
		return isEqual(left, right)
	case TokenGreater:
		lval, rval := convertNumberOperands(expr.Operator, left, right)
		return lval > rval
	case TokenGreaterEqual:
		lval, rval := convertNumberOperands(expr.Operator, left, right)
		return lval >= rval
	case TokenLess:
		lval, rval := convertNumberOperands(expr.Operator, left, right)
		return lval < rval
	case TokenLessEqual:
		lval, rval := convertNumberOperands(expr.Operator, left, right)
		return lval <= rval

	// arithmetics
	case TokenMinus:
		lval, rval := convertNumberOperands(expr.Operator, left, right)
		return lval - rval
	case TokenPlus:
		// string concat is supported.
		lval, ok1 := left.(string)
		rval, ok2 := right.(string)
		if ok1 == true && ok2 == true {
			return lval + rval
		}

		lval2, rval2 := convertNumberOperands(expr.Operator, left, right)
		return lval2 + rval2
	case TokenStar:
		lval, rval := convertNumberOperands(expr.Operator, left, right)
		return lval * rval
	case TokenSlash:
		lval, rval := convertNumberOperands(expr.Operator, left, right)
		return lval / rval
	case TokenPercent:
		lval, rval := convertNumberOperands(expr.Operator, left, right)
		return int(lval) % int(rval)
	default:
		return nil
	}
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) interface{} {
	innerVal := i.evaluate(expr.Expression)
	return innerVal
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitLogicalExpr(expr *Logical) interface{} {
	left := i.evaluate(expr.Left)

	if expr.Operator.Type == TokenOr {
		if isTruthy(left) {
			return left
		}
	} else { // TokenAnd
		if !isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) interface{} {
	operator := expr.Operator
	value := i.evaluate(expr.Right)
	switch operator.Type {
	case TokenMinus:
		num := convertNumberOperand(operator, value)
		return 0 - num
	case TokenBang:
		return !isTruthy(value)
	}
	return nil
}

func convertNumberOperand(operator *Token, operand interface{}) float64 {
	val, ok := operand.(float64)
	if ok == true {
		return val
	}
	errmsg := "Operand must be number."
	RuntimeError(operator, errmsg)
	panic(errmsg)
}

func convertNumberOperands(operator *Token, left, right interface{}) (float64, float64) {
	lval, ok1 := left.(float64)
	rval, ok2 := right.(float64)

	if ok1 == true && ok2 == true {
		return lval, rval
	}
	errmsg := "Operands must be number."
	RuntimeError(operator, errmsg)
	panic(errmsg)
}

func isEqual(left, right interface{}) bool {
	switch lval := left.(type) {
	case string, float64:
		rval, ok := right.(float64)
		if ok != true || lval != rval {
			return false
		}
		return true
	}
	return false
}

func isTruthy(object interface{}) bool {
	switch val := object.(type) {
	case nil:
		return false
	case bool:
		return val
	default:
		return true
	}
}
