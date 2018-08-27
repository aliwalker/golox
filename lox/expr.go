package lox

type ExprVisitor interface {
	VisitAssignExpr(expr *Assign) interface{}
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitLogicalExpr(expr *Logical) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
}

type Expr interface {
	Accept(v ExprVisitor) interface{}
}

type Assign struct {
	Name  *Token
	Value Expr
}

func NewAssign(name *Token, value Expr) Expr {
	return &Assign{Name: name, Value: value}
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

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) Expr {
	return &Grouping{Expression: expression}
}
func (expr *Grouping) Accept(v ExprVisitor) interface{} {
	return v.VisitGroupingExpr(expr)
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
