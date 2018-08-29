package lox

type StmtVisitor interface {
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitPrintStmt(stmt *Print) interface{}
	VisitVarStmt(stmt *Var) interface{}
}

type Stmt interface {
	Accept(v StmtVisitor) interface{}
}

type Expression struct {
	Expression Expr
}

func NewExpression(expression Expr) Stmt {
	return &Expression{Expression: expression}
}
func (expr *Expression) Accept(v StmtVisitor) interface{} {
	return v.VisitExpressionStmt(expr)
}

type Print struct {
	Expression Expr
}

func NewPrint(expression Expr) Stmt {
	return &Print{Expression: expression}
}
func (expr *Print) Accept(v StmtVisitor) interface{} {
	return v.VisitPrintStmt(expr)
}

type Var struct {
	name *Token
	expr Expr
}

func NewVar(name *Token, expr Expr) Stmt {
	return &Var{name: name, expr: expr}
}
func (expr *Var) Accept(v StmtVisitor) interface{} {
	return v.VisitVarStmt(expr)
}
