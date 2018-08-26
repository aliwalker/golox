package golox

type StmtVisitor interface {
	VisitExpressionStmt(expr *Expression) interface{}
	VisitPrintStmt(expr *Print) interface{}
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
