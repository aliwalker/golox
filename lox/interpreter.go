package lox

import (
	"fmt"
)

type Interpreter struct {
	hadRuntimeError bool
	environment     *Environment
	global          *Environment
}

func NewInterpreter() *Interpreter {
	global := NewEnvironment(nil)
	environment := global

	return &Interpreter{false, global, environment}
}

func (i *Interpreter) Interprete(stmts []Stmt) (hadRuntimeError bool) {
	defer func() {
		if val := recover(); val != nil {
			runtimeError := val.(*RuntimeError)
			fmt.Println(runtimeError.Error())
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

func (i *Interpreter) VisitBlockStmt(stmt *Block) interface{} {
	enclosingEnv := i.environment
	i.environment = NewEnvironment(enclosingEnv)

	defer func() {
		i.environment = enclosingEnv
	}()

	for _, stmt := range stmt.Stmts {
		i.execute(stmt)
	}

	return nil
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

func (i *Interpreter) VisitVarStmt(stmt *Var) interface{} {
	// TODO: add scope variable.
	var (
		identifier  *Token
		initializer Expr
		initVal     interface{}
	)

	identifier = stmt.name
	if initializer = stmt.expr; initializer != nil {
		initVal = i.evaluate(initializer)
	}

	i.environment.Define(identifier.Lexeme, initVal)
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *Assign) interface{} {
	value := i.evaluate(expr.Value)
	i.environment.Assign(expr.Name, value)
	return value
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	var lval, rval float64
	var bothInt bool

	switch expr.Operator.Type {
	// comparison.
	case TokenBangEqual:
		return !equal(left, right)
	case TokenEqualEqual:
		return equal(left, right)
	case TokenGreater:
		lval, rval, _ = convertFloatOperands(expr.Operator, left, right)
		return lval > rval
	case TokenGreaterEqual:
		lval, rval, _ = convertFloatOperands(expr.Operator, left, right)
		return lval >= rval
	case TokenLess:
		lval, rval, _ = convertFloatOperands(expr.Operator, left, right)
		return lval < rval
	case TokenLessEqual:
		lval, rval, _ = convertFloatOperands(expr.Operator, left, right)
		return lval <= rval

	// arithmetics
	case TokenMinus:
		if lval, rval, bothInt = convertFloatOperands(expr.Operator, left, right); bothInt == true {
			return int(lval) - int(rval)
		}
		return lval - rval
	case TokenPlus:
		// string concat is supported.
		lvalString, ok1 := left.(string)
		rvalString, ok2 := right.(string)
		if ok1 == true && ok2 == true {
			return lvalString + rvalString
		}

		if lval, rval, bothInt = convertFloatOperands(expr.Operator, left, right); bothInt == true {
			return int(lval) + int(rval)
		}
		return lval + rval
	case TokenStar:
		if lval, rval, bothInt = convertFloatOperands(expr.Operator, left, right); bothInt == true {
			return int(lval) * int(rval)
		}
		return lval * rval
	case TokenSlash:
		if lval, rval, bothInt = convertFloatOperands(expr.Operator, left, right); bothInt == true {
			return int(lval) / int(rval)
		}
		return lval / rval
	case TokenPercent:
		if lval, rval, bothInt = convertFloatOperands(expr.Operator, left, right); bothInt == false {
			panic(NewRuntimeError(expr.Operator, "both operands both be integers."))
		}
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
		if truthy(left) {
			return left
		}
	} else { // TokenAnd
		if !truthy(left) {
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
		num, isInt := convertNumberOperand(operator, value)
		if isInt == true {
			return int(-num)
		}
		return -num
	case TokenBang:
		return !truthy(value)
	}
	return nil
}

func (i *Interpreter) VisitVariableExpr(expr *Variable) interface{} {
	val := lookUpVariable(i.environment, expr.name)
	return val
}

// returns a number & a bool indicates if this number is int.
func convertNumberOperand(operator *Token, operand interface{}) (float64, bool) {
	if val, ok := operand.(float64); ok == true {
		return val, false
	}

	if val, ok := operand.(int); ok == true {
		return float64(val), true
	}

	panic(NewRuntimeError(operator, "Operand must be number."))
}

// returns both operands as float64, & a bool value indicates if both operands are int.
func convertFloatOperands(operator *Token, left, right interface{}) (float64, float64, bool) {
	lval, isInt1 := convertNumberOperand(operator, left)
	rval, isInt2 := convertNumberOperand(operator, right)

	return lval, rval, isInt1 && isInt2
}

func equal(left, right interface{}) bool {
	switch lval := left.(type) {
	case string, float64, int:
		rval, ok := right.(float64)
		if ok != true || lval != rval {
			return false
		}
		return true
	}
	return false
}

// TODO: add support to resolution pass.
func lookUpVariable(env *Environment, name *Token) interface{} {
	return env.Get(name)
}

func truthy(object interface{}) bool {
	switch val := object.(type) {
	case nil:
		return false
	case bool:
		return val
	default:
		return true
	}
}
