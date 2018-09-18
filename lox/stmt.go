package lox

type StmtVisitor interface {
	VisitBlockStmt(stmt *Block) interface{}
	VisitClassStmt(stmt *Class) interface{}
	VisitControlStmt(stmt *Control) interface{}
	VisitFunctionStmt(stmt *Function) interface{}
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitIfStmt(stmt *If) interface{}
	VisitPrintStmt(stmt *Print) interface{}
	VisitVarStmt(stmt *Var) interface{}
	VisitVarListStmt(stmt *VarList) interface{}
	VisitWhileStmt(stmt *While) interface{}
}

type Stmt interface {
	Accept(v StmtVisitor) interface{}
}

type Block struct {
	Stmts []Stmt
}

func NewBlock(stmts []Stmt) Stmt {
	return &Block{Stmts: stmts}
}
func (expr *Block) Accept(v StmtVisitor) interface{} {
	return v.VisitBlockStmt(expr)
}

type Class struct {
	Name    *Token
	Super   *Variable
	Statics []*Function
	Methods []*Function
	Getters []*Function
	Setters []*Function
}

func NewClass(name *Token, super *Variable, statics []*Function, methods []*Function, getters []*Function, setters []*Function) Stmt {
	return &Class{Name: name, Super: super, Statics: statics, Methods: methods, Getters: getters, Setters: setters}
}
func (expr *Class) Accept(v StmtVisitor) interface{} {
	return v.VisitClassStmt(expr)
}

type Control struct {
	Keyword  *Token
	CtrlType ControlType
	Value    Expr
}

func NewControl(keyword *Token, ctrltype ControlType, value Expr) Stmt {
	return &Control{Keyword: keyword, CtrlType: ctrltype, Value: value}
}
func (expr *Control) Accept(v StmtVisitor) interface{} {
	return v.VisitControlStmt(expr)
}

type Function struct {
	Name   *Token
	Params []*Token
	Body   []Stmt
}

func NewFunction(name *Token, params []*Token, body []Stmt) Stmt {
	return &Function{Name: name, Params: params, Body: body}
}
func (expr *Function) Accept(v StmtVisitor) interface{} {
	return v.VisitFunctionStmt(expr)
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

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIf(condition Expr, thenbranch Stmt, elsebranch Stmt) Stmt {
	return &If{Condition: condition, ThenBranch: thenbranch, ElseBranch: elsebranch}
}
func (expr *If) Accept(v StmtVisitor) interface{} {
	return v.VisitIfStmt(expr)
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
	Name        *Token
	Initializer Expr
}

func NewVar(name *Token, initializer Expr) Stmt {
	return &Var{Name: name, Initializer: initializer}
}
func (expr *Var) Accept(v StmtVisitor) interface{} {
	return v.VisitVarStmt(expr)
}

type VarList struct {
	stmts []*Var
}

func NewVarList(stmts []*Var) Stmt {
	return &VarList{stmts: stmts}
}
func (expr *VarList) Accept(v StmtVisitor) interface{} {
	return v.VisitVarListStmt(expr)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func NewWhile(condition Expr, body Stmt) Stmt {
	return &While{Condition: condition, Body: body}
}
func (expr *While) Accept(v StmtVisitor) interface{} {
	return v.VisitWhileStmt(expr)
}
