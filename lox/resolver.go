package lox

import (
	"fmt"
)

type FuncType int

const (
	_ FuncType = iota
	FuncNone
	FuncFunc
	FuncMeth
)

// Resolver resolves bindings.
type Resolver struct {
	scopes      *Scopes
	interpreter *Interpreter
	curFunc     FuncType
	inLoop      bool
	hadError    bool
}

// NewResolver returns a new resolver.
func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{NewScopes(), interpreter, FuncNone, false, false}
}

// BeginScope is called when resolver enters a new scope.
func (r *Resolver) BeginScope() {
	r.scopes.Push()
}

// EndScope is called when resolver exists a scope.
func (r *Resolver) EndScope() {
	r.scopes.Pop()
}

// Declare marks a variable is being declared yet available for used in current scope.
func (r *Resolver) Declare(name *Token) {
	// don't check global scope.
	if r.scopes.Empty() {
		return
	}

	scope := r.scopes.Peek()

	if exists := scope[name.Lexeme]; exists == varDeclared || exists == varDefined {
		panic(NewLoxError(name, "variable redeclared."))
	}

	scope[name.Lexeme] = varDeclared
}

// Define marks a variable is available for used in current scope.
func (r *Resolver) Define(name *Token) {
	if r.scopes.Empty() {
		return
	}

	scope := r.scopes.Peek()
	scope[name.Lexeme] = varDefined
}

func (r *Resolver) resolve(node interface{}) {
	// We know for sure `node` is either Stmt or Expr.
	if stmt, ok := node.(Stmt); ok {
		stmt.Accept(r)
	} else if stmts, ok := node.([]Stmt); ok {
		for _, stmt := range stmts {
			r.resolve(stmt)
		}
	} else {
		expr, _ := node.(Expr)
		expr.Accept(r)
	}
}

func (r *Resolver) resolveStmts(stmts []Stmt) {
	defer func() {
		if val := recover(); val != nil {
			r.hadError = true
			error := val.(*LoxError)
			fmt.Println(error.Error())
		}
	}()

	for _, stmt := range stmts {
		stmt.Accept(r)
	}
}

// Resolve resolves names referenced in `stmts`.
func (r *Resolver) Resolve(stmts []Stmt) bool {
	r.resolveStmts(stmts)
	return r.hadError
}

// ResolveLocal resolves the `name` referenced by `expr`.
// Once it is resolved, the resolver informs the interpreter, the distance from
// the current scope to resolved scope.
func (r *Resolver) resolveLocal(expr Expr, name *Token) {
	for i := r.scopes.Len() - 1; i >= 0; i-- {
		scope := r.scopes.Get(i)
		if scope.HasName(name.Lexeme) {
			r.interpreter.resolve(expr, r.scopes.Len()-1-i)
			return
		}
	}
	// if the execution reaches out of the `for` loop above, assume `name` is defined in global.
}

func (r *Resolver) VisitAssignExpr(expr *Assign) interface{} {
	// check the operator at runtime.
	r.resolve(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *Binary) interface{} {
	r.resolve(expr.Left)
	r.resolve(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *Call) interface{} {
	r.resolve(expr.Callee)

	for _, arg := range expr.Arguments {
		r.resolve(arg)
	}
	return nil
}

func (r *Resolver) VisitGetExpr(expr *Get) interface{} {
	r.resolve(expr.Object)
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *Grouping) interface{} {
	r.resolve(expr.Expression)
	return nil
}

func (r *Resolver) VisitLambdaExpr(expr *Lambda) interface{} {
	r.resolveFunction(expr.LambdaFunc, FuncFunc)
	return nil
}

// VisitLiteralExpr doesn't do anything because there's nothing to resolve.
func (r *Resolver) VisitLiteralExpr(expr *Literal) interface{} {
	return nil
}

// VisitLogicalExpr resolves Left & Right.
func (r *Resolver) VisitLogicalExpr(expr *Logical) interface{} {
	r.resolve(expr.Left)
	r.resolve(expr.Right)
	return nil
}

// VisitSetExpr resolves set expression.
func (r *Resolver) VisitSetExpr(expr *Set) interface{} {
	r.resolve(expr.Value)
	r.resolve(expr.Object)
	return nil
}

// VisitThisExpr resolves "this"
func (r *Resolver) VisitThisExpr(expr *This) interface{} {
	r.resolveLocal(expr, expr.Keyword)
	return nil
}

// VisitUnaryExpr resolves Right.
func (r *Resolver) VisitUnaryExpr(expr *Unary) interface{} {
	r.resolve(expr.Right)
	return nil
}

// VisitVariableExpr makes sure a variable is not referenced during being declared.
func (r *Resolver) VisitVariableExpr(expr *Variable) interface{} {
	if !r.scopes.Empty() && r.scopes.Peek()[expr.Name.Lexeme] == varDeclared {
		panic(NewLoxError(expr.Name, "cannot read variable being declared."))
	}
	r.resolveLocal(expr, expr.Name)
	return nil
}

// statements

func (r *Resolver) VisitBlockStmt(stmt *Block) interface{} {
	r.BeginScope()
	r.resolve(stmt.Stmts)
	r.EndScope()
	return nil
}

func (r *Resolver) VisitClassStmt(stmt *Class) interface{} {
	r.Declare(stmt.Name)
	r.Define(stmt.Name)

	// Since we added "this", we need another layer between the scope containing the class
	// and the method scope.
	r.BeginScope()
	r.scopes.Peek()["this"] = varDefined

	for _, f := range stmt.Methods {
		r.resolveFunction(f, FuncMeth)
	}

	r.EndScope()
	return nil
}

// VisitControlStmt interpretes "break" & "return" statements.
func (r *Resolver) VisitControlStmt(stmt *Control) interface{} {
	if stmt.CtrlType == ControlReturn {
		if r.curFunc == FuncNone {
			panic(NewLoxError(stmt.Keyword, "illegal return statement."))
		}

		if stmt.Value != nil {
			r.resolve(stmt.Value)
		}
	}

	if stmt.CtrlType == ControlBreak && r.inLoop == false {
		panic(NewLoxError(stmt.Keyword, "illegal break statement."))
	}
	return nil
}

// VisitFunctionStmt resolves function declaration statement.
func (r *Resolver) VisitFunctionStmt(stmt *Function) interface{} {
	r.Declare(stmt.Name)
	r.Define(stmt.Name)

	r.resolveFunction(stmt, FuncFunc)
	return nil
}

// This function is used in resolving function and method.
// `fType` is the passed from caller to indicate whether this is a method or a
// function.
func (r *Resolver) resolveFunction(function *Function, fType FuncType) {
	enclosingFunc := r.curFunc
	r.curFunc = fType
	defer func() {
		r.curFunc = enclosingFunc
	}()

	r.BeginScope()
	for _, param := range function.Params {
		r.Declare(param)
		r.Define(param)
	}
	r.resolve(function.Body)
	r.EndScope()

}

func (r *Resolver) VisitExpressionStmt(stmt *Expression) interface{} {
	r.resolve(stmt.Expression)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *If) interface{} {
	r.resolve(stmt.Condition)
	r.resolve(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolve(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *Print) interface{} {
	r.resolve(stmt.Expression)
	return nil
}

func (r *Resolver) VisitVarStmt(stmt *Var) interface{} {
	r.Declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolve(stmt.Initializer)
	}
	r.Define(stmt.Name)
	return nil
}

func (r *Resolver) VisitVarListStmt(stmt *VarList) interface{} {
	varDecs := stmt.stmts

	for _, varDec := range varDecs {
		r.resolve(varDec)
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *While) interface{} {
	r.resolve(stmt.Condition)

	preLoop := r.inLoop
	r.inLoop = true
	defer func() {
		r.inLoop = preLoop
	}()

	r.resolve(stmt.Body)
	return nil
}
