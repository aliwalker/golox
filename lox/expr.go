package lox

type ExprVisitor interface {
	VisitAssignExpr(expr *Assign) interface{}
	VisitBinaryExpr(expr *Binary) interface{}
	VisitCallExpr(expr *Call) interface{}
	VisitGetExpr(expr *Get) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLambdaExpr(expr *Lambda) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitLogicalExpr(expr *Logical) interface{}
	VisitSetExpr(expr *Set) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitVariableExpr(expr *Variable) interface{}
}

type Expr interface {
	Accept(v ExprVisitor) interface{}
}

type Assign struct {
	Name     *Token
	Operator *Token
	Value    Expr
}

func NewAssign(name *Token, operator *Token, value Expr) Expr {
	return &Assign{Name: name, Operator: operator, Value: value}
}
func (expr *Assign) Accept(v ExprVisitor) interface{} {
	return v.VisitAssignExpr(expr)
}

type Binary struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func NewBinary(left Expr, operator *Token, right Expr) Expr {
	return &Binary{Left: left, Operator: operator, Right: right}
}
func (expr *Binary) Accept(v ExprVisitor) interface{} {
	return v.VisitBinaryExpr(expr)
}

type Call struct {
	Callee    Expr
	Paren     *Token
	Arguments []Expr
}

func NewCall(callee Expr, paren *Token, arguments []Expr) Expr {
	return &Call{Callee: callee, Paren: paren, Arguments: arguments}
}
func (expr *Call) Accept(v ExprVisitor) interface{} {
	return v.VisitCallExpr(expr)
}

type Get struct {
	Object Expr
	Name   *Token
}

func NewGet(object Expr, name *Token) Expr {
	return &Get{Object: object, Name: name}
}
func (expr *Get) Accept(v ExprVisitor) interface{} {
	return v.VisitGetExpr(expr)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) Expr {
	return &Grouping{Expression: expression}
}
func (expr *Grouping) Accept(v ExprVisitor) interface{} {
	return v.VisitGroupingExpr(expr)
}

type Lambda struct {
	LambdaFunc *Function
}

func NewLambda(lambdafunc *Function) Expr {
	return &Lambda{LambdaFunc: lambdafunc}
}
func (expr *Lambda) Accept(v ExprVisitor) interface{} {
	return v.VisitLambdaExpr(expr)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(value interface{}) Expr {
	return &Literal{Value: value}
}
func (expr *Literal) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteralExpr(expr)
}

type Logical struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func NewLogical(left Expr, operator *Token, right Expr) Expr {
	return &Logical{Left: left, Operator: operator, Right: right}
}
func (expr *Logical) Accept(v ExprVisitor) interface{} {
	return v.VisitLogicalExpr(expr)
}

type Set struct {
	Object Expr
	Name   *Token
	Value  Expr
}

func NewSet(object Expr, name *Token, value Expr) Expr {
	return &Set{Object: object, Name: name, Value: value}
}
func (expr *Set) Accept(v ExprVisitor) interface{} {
	return v.VisitSetExpr(expr)
}

type Unary struct {
	Operator *Token
	Right    Expr
}

func NewUnary(operator *Token, right Expr) Expr {
	return &Unary{Operator: operator, Right: right}
}
func (expr *Unary) Accept(v ExprVisitor) interface{} {
	return v.VisitUnaryExpr(expr)
}

type Variable struct {
	Name *Token
}

func NewVariable(name *Token) Expr {
	return &Variable{Name: name}
}
func (expr *Variable) Accept(v ExprVisitor) interface{} {
	return v.VisitVariableExpr(expr)
}
