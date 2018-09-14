package lox

import (
	"fmt"
)

type AstPrinter struct {
	indents int
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{indents: 0}
}

func getIndents(num int) string {
	indents := ""

	for i := 0; i < num; i++ {
		indents += "	"
	}
	return indents
}

func (p *AstPrinter) Print(stmts []Stmt) {
	ast := ""

	for _, stmt := range stmts {
		val, _ := stmt.Accept(p).(string)
		ast += val
	}

	fmt.Println(ast)
}

func (p *AstPrinter) VisitBlockStmt(stmt *Block) interface{} {
	ast := getIndents(p.indents) + "(block \n"
	p.indents++

	defer func() {
		p.indents--
	}()

	for _, st := range stmt.Stmts {
		val, _ := st.Accept(p).(string)
		ast += val
	}
	ast += ")\n"
	return ast
}
func (p *AstPrinter) VisitClassStmt(stmt *Class) interface{} {
	ast := getIndents(p.indents) + "(class " + stmt.Name.Lexeme + "\n"

	p.indents++
	defer func() {
		p.indents--
	}()

	for _, method := range stmt.Methods {
		val, _ := method.Accept(p).(string)
		ast += val
	}
	ast += ")\n"
	return ast
}
func (p *AstPrinter) VisitControlStmt(stmt *Control) interface{} {
	ast := getIndents(p.indents) + "(" + stmt.Keyword.Lexeme
	if stmt.Keyword.Lexeme == "break" {
		return ast
	}

	// return statement.
	if stmt.Value != nil {
		ast += " "

		val, _ := stmt.Value.Accept(p).(string)
		ast += val
		ast += ")\n"
	}
	return ast
}

func (p *AstPrinter) VisitFunctionStmt(stmt *Function) interface{} {
	ast := getIndents(p.indents) + "(fun " + stmt.Name.Lexeme + "("
	for i, param := range stmt.Params {
		if i != 0 {
			ast += " "
		}
		ast += param.Lexeme
	}
	ast += ")\n"

	p.indents++
	defer func() {
		p.indents--
	}()

	for _, st := range stmt.Body {
		val, _ := st.Accept(p).(string)
		ast += val
	}

	return ast
}

func (p *AstPrinter) VisitExpressionStmt(stmt *Expression) interface{} {
	return getIndents(p.indents) + p.parenthesize(";", stmt.Expression) + "\n"
}

func (p *AstPrinter) VisitIfStmt(stmt *If) interface{} {
	if stmt.ElseBranch == nil {
		return getIndents(p.indents) + p.parenthesize("if", stmt.Condition, stmt.ThenBranch) + "\n"
	}
	return getIndents(p.indents) +
		p.parenthesize("if-else", stmt.Condition, stmt.ThenBranch, stmt.ElseBranch) + "\n"
}

func (p *AstPrinter) VisitPrintStmt(stmt *Print) interface{} {
	return getIndents(p.indents) + p.parenthesize("print", stmt.Expression) + "\n"
}

func (p *AstPrinter) VisitVarStmt(stmt *Var) interface{} {
	if stmt.Initializer == nil {
		return getIndents(p.indents) + p.parenthesize("var", stmt.Name) + "\n"
	}
	return getIndents(p.indents) + p.parenthesize("var", stmt.Name, "=", stmt.Initializer) + "\n"
}

func (p *AstPrinter) VisitVarListStmt(stmt *VarList) interface{} {
	ast := ""

	for _, v := range stmt.stmts {
		val, _ := p.VisitVarStmt(v).(string)
		ast += val
	}

	return ast
}

func (p *AstPrinter) VisitWhileStmt(stmt *While) interface{} {
	return p.parenthesize("while", stmt.Condition, stmt.Body)
}

func (p *AstPrinter) VisitAssignExpr(expr *Assign) interface{} {
	return p.parenthesize("=", expr.Name.Lexeme, expr.Value)
}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitCallExpr(expr *Call) interface{} {
	return p.parenthesize("call", expr.Callee, expr.Arguments)
}

func (p *AstPrinter) VisitGetExpr(expr *Get) interface{} {
	return p.parenthesize("get", expr.Object, expr.Name)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}

	switch v := expr.Value.(type) {
	case string:
		return v
	case int, float64:
		return fmt.Sprintf("%v", v)
	default:
		panic("unknown literal.")
	}
}

func (p *AstPrinter) VisitLambdaExpr(expr *Lambda) interface{} {
	lambda := expr.LambdaFunc

	ast := getIndents(p.indents) + "(lambda " + "("
	for i, param := range lambda.Params {
		if i != 0 {
			ast += " "
		}
		ast += param.Lexeme
	}
	ast += ")\n"

	p.indents++
	defer func() {
		p.indents--
	}()

	for _, st := range lambda.Body {
		val, _ := st.Accept(p).(string)
		ast += val
	}

	return ast
}

func (p *AstPrinter) VisitLogicalExpr(expr *Logical) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitSetExpr(expr *Set) interface{} {
	return p.parenthesize("set", expr.Object, expr.Name, expr.Value)
}

func (p *AstPrinter) VisitThisExpr(expr *This) interface{} {
	return "this"
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) VisitVariableExpr(expr *Variable) interface{} {
	return expr.Name.Lexeme
}

func (p *AstPrinter) parenthesize(name string, values ...interface{}) string {
	ast := "(" + name

	for _, obj := range values {
		ast += " "
		switch v := obj.(type) {
		case Expr:
			val, _ := v.Accept(p).(string)
			ast += val
		case Stmt:
			val, _ := v.Accept(p).(string)
			ast += val
		case *Token:
			ast += v.Lexeme
		default:
			val, _ := obj.(string)
			ast += val
		}
	}
	ast += ")"
	return ast
}
