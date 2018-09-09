package lox

import (
	"fmt"
)

type Interpreter struct {
	repl            bool         // REPL mode or not.
	hadRuntimeError bool         // indicates runtime error.
	environment     *Environment // current environment.
	global          *Environment // global environment.
	locals          map[Expr]int // for local variable resolution.
}

func NewInterpreter(repl bool) *Interpreter {
	global := NewEnvironment(nil)
	environment := global

	return &Interpreter{
		repl:            repl,
		hadRuntimeError: false,
		environment:     global,
		global:          environment,
		locals:          map[Expr]int{},
	}
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

// execute a block in `env`.
func (i *Interpreter) executeBlock(stmts []Stmt, env *Environment) {
	prevEnv := i.environment
	i.environment = env

	for _, stmt := range stmts {
		i.execute(stmt)
	}

	i.environment = prevEnv
}

func (i *Interpreter) resolve(expr Expr, distance int) {
	i.locals[expr] = distance
}

func (i *Interpreter) VisitBlockStmt(stmt *Block) interface{} {
	i.executeBlock(stmt.Stmts, NewEnvironment(i.environment))
	return nil
}

func (i *Interpreter) VisitClassStmt(stmt *Class) interface{} {
	i.environment.Define(stmt.Name.Lexeme, nil)
	class := NewLoxClass(stmt.Name.Lexeme)
	i.environment.Assign(stmt.Name, class)
	return nil
}

func (i *Interpreter) VisitControlStmt(stmt *Control) interface{} {
	// throw it to unwind the call stack.
	panic(stmt)
}

func (i *Interpreter) VisitExpressionStmt(stmt *Expression) interface{} {
	val := i.evaluate(stmt.Expression)
	if i.repl {
		fmt.Println(val)
	}
	return nil
}

// convert function ast node to runtime function object.
func (i *Interpreter) VisitFunctionStmt(stmt *Function) interface{} {
	i.environment.Define(stmt.Name.Lexeme, NewLoxFunction(stmt, i.environment))
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *If) interface{} {
	if truthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) interface{} {
	val := i.evaluate(stmt.Expression)
	fmt.Println(val)
	return nil
}

/*
func stringify(value interface{}) interface{} {
	switch v := value.(type) {
	case int, float64, string:
		return v
	case printable:
		return v.String()
	}
	return value
}*/

func (i *Interpreter) VisitVarStmt(stmt *Var) interface{} {
	var (
		identifier  *Token
		initializer Expr
		initVal     interface{}
	)

	identifier = stmt.Name
	if initializer = stmt.Initializer; initializer != nil {
		initVal = i.evaluate(initializer)
	}

	i.environment.Define(identifier.Lexeme, initVal)
	return nil
}

func (i *Interpreter) VisitVarListStmt(stmt *VarList) interface{} {
	varDecs := stmt.stmts

	for _, varDec := range varDecs {
		i.execute(varDec)
	}

	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *While) interface{} {
	defer func() {
		if val := recover(); val != nil {
			control, ok := val.(*Control)
			// repanic if it is not a Control.
			if ok != true {
				panic(val)
			}
			// repanic if it is a ControlReturn
			if control.CtrlType != ControlBreak {
				panic(val)
			}
		}
	}()

	for truthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}

	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *Assign) interface{} {
	var operator TokenType

	value := i.evaluate(expr.Value)
	switch expr.Operator.Type {
	case TokenPlusEqual:
		operator = TokenPlus
	case TokenMinusEqual:
		operator = TokenMinus
	case TokenStarEqual:
		operator = TokenStar
	case TokenSlashEqual:
		operator = TokenSlash
	case TokenPercentEqual:
		operator = TokenPercent
	default:
		operator = 0
	}

	if operator != 0 {
		lval := NewLiteral(i.lookUpVariable(expr, expr.Name))
		rval := NewLiteral(value)

		binary, _ := NewBinary(lval, NewToken(operator, "", nil, expr.Operator.Line), rval).(*Binary)
		value = i.VisitBinaryExpr(binary)
	}

	distance, ok := i.locals[expr]
	if ok {
		i.environment.AssignAt(distance, expr.Name, value)
	} else {
		i.global.Assign(expr.Name, value)
	}
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

func (i *Interpreter) VisitCallExpr(expr *Call) interface{} {
	callee := i.evaluate(expr.Callee)

	function, ok := callee.(Callable)
	if ok != true {
		panic(NewRuntimeError(expr.Paren, "callee is not callable."))
	}

	if len(expr.Arguments) != function.Arity() {
		panic(NewRuntimeError(expr.Paren, fmt.Sprintf("expect %v arguments, but got %v", function.Arity(), len(expr.Arguments))))
	}

	args := make([]interface{}, 0)
	for _, arg := range expr.Arguments {
		args = append(args, i.evaluate(arg))
	}

	// TODO: add return value.
	return function.Call(i, args...)
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
	val := i.lookUpVariable(expr, expr.Name)
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

// TODO: this might be problematic.
func equal(left, right interface{}) bool {
	if left == right {
		return true
	}

	return false
}

func (i *Interpreter) lookUpVariable(expr Expr, name *Token) interface{} {
	distance, ok := i.locals[expr]
	if ok {
		return i.environment.GetAt(distance, name)
	}
	return i.global.Get(name)
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
