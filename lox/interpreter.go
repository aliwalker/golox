package lox

import (
	"fmt"

	"github.com/fatih/color"
)

// Interpreter is an object interprets our AST.
type Interpreter struct {
	repl            bool         // REPL mode or not.
	hadRuntimeError bool         // indicates runtime error.
	environment     *Environment // current environment.
	global          *Environment // global environment.
	locals          map[Expr]int // for local variable resolution.
}

// NewInterpreter returns an interpreter object.
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
	// evaluate the super class first.
	var superClass *LoxClass
	if stmt.Super != nil {
		var ok bool
		super := i.evaluate(stmt.Super)
		if superClass, ok = super.(*LoxClass); ok != true {
			panic(NewRuntimeError(stmt.Super.Name, "superclass must be a class."))
		}
	}

	i.environment.Define(stmt.Name.Lexeme, nil)

	if stmt.Super != nil {
		// add another scope for "super".
		i.environment = NewEnvironment(i.environment)
		i.environment.Define("super", superClass)
	}

	statics := map[string]*LoxFunction{}
	for _, static := range stmt.Statics {
		statics[static.Name.Lexeme] = NewLoxFunction(static, i.environment)
	}

	methods := map[string]*LoxFunction{}
	for _, method := range stmt.Methods {
		methods[method.Name.Lexeme] = NewLoxFunction(method, i.environment)
	}

	getters := map[string]*LoxFunction{}
	for _, getter := range stmt.Getters {
		getters[getter.Name.Lexeme] = NewLoxFunction(getter, i.environment)
	}

	setters := map[string]*LoxFunction{}
	for _, setter := range stmt.Setters {
		setters[setter.Name.Lexeme] = NewLoxFunction(setter, i.environment)
	}

	class := NewLoxClass(stmt.Name.Lexeme, superClass, statics, methods, getters, setters)

	if stmt.Super != nil {
		// remember to exist the scope created previously,
		// before assigning the current class to current env.
		i.environment = i.environment.enclosing
	}

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
		if val == nil {
			return nil
		}
		color.Cyan("%v", val)
	}
	return nil
}

// VisitFunctionStmt converts function ast node to runtime function object.
// This function adds an entry to the current env, while methods in a class don't.
func (i *Interpreter) VisitFunctionStmt(stmt *Function) interface{} {
	i.environment.Define(stmt.Name.Lexeme, NewLoxFunction(stmt, i.environment))
	return nil
}

// VisitIfStmt interpretes an if statement.
func (i *Interpreter) VisitIfStmt(stmt *If) interface{} {
	if truthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

// VisitPrintStmt prints an expression in Cyan color.
func (i *Interpreter) VisitPrintStmt(stmt *Print) interface{} {
	val := i.evaluate(stmt.Expression)
	color.Cyan("%v", val)
	return nil
}

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

	return function.Call(i, args...)
}

func (i *Interpreter) VisitGetExpr(expr *Get) interface{} {
	if object, ok := i.evaluate(expr.Object).(*LoxInstance); ok {
		return object.Get(i, expr.Name)
	}

	if class, ok := i.evaluate(expr.Object).(*LoxClass); ok {
		return class.FindStatic(expr.Name.Lexeme)
	}

	panic(NewRuntimeError(expr.Name, "unexpected property access."))
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) interface{} {
	innerVal := i.evaluate(expr.Expression)
	return innerVal
}

func (i *Interpreter) VisitLambdaExpr(expr *Lambda) interface{} {
	return NewLoxFunction(expr.LambdaFunc, i.environment)
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

func (i *Interpreter) VisitSetExpr(expr *Set) interface{} {
	var (
		loxInstance *LoxInstance
		value       interface{}
		ok          bool
	)

	value = i.evaluate(expr.Object)

	if loxInstance, ok = value.(*LoxInstance); ok != true {
		panic(NewRuntimeError(expr.Name, "set property on a non Lox instance object."))
	}

	value = i.evaluate(expr.Value)
	loxInstance.Set(i, expr.Name, value)
	return value
}

// VisitSuperExpr interpretes something like "super.foo"
func (i *Interpreter) VisitSuperExpr(expr *Super) interface{} {
	distance := i.locals[expr]
	superClass, _ := i.environment.GetAt(distance, "super").(*LoxClass)
	// the context for the method queryed. This is a little hack since we've known
	// it is there, and we've known it must be a LoxInstance.
	object, _ := i.environment.GetAt(distance-1, "this").(*LoxInstance)
	// TODO: add getter/setter inheritance support.
	method := superClass.FindMethod(object, expr.Method.Lexeme)

	if method == nil {
		panic(NewRuntimeError(expr.Method, "undefined property '"+expr.Method.Lexeme+"'."))
	}

	return method
}

func (i *Interpreter) VisitThisExpr(expr *This) interface{} {
	return i.lookUpVariable(expr, expr.Keyword)
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
	if distance, ok := i.locals[expr]; ok {
		return i.environment.GetAt(distance, name.Lexeme)
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
