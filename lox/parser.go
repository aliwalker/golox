package lox

import (
	"fmt"
)

// Parser parses the tokens into an AST.
type Parser struct {
	tokens   []*Token
	current  int
	hadError bool
}

// NewParser creates a parser.
func NewParser(tokens []*Token) *Parser {
	return &Parser{tokens, 0, false}
}

func (p *Parser) advance() *Token {
	if !p.end() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(t TokenType) bool {
	if p.end() {
		return false
	}

	currentType := p.tokens[p.current].Type
	if currentType == t {
		return true
	}
	return false
}

func (p *Parser) consume(t TokenType, message string) *Token {
	if p.check(t) {
		return p.advance()
	}
	panic(NewLoxError(p.peek(), message))
}

func (p *Parser) end() bool {
	if p.tokens[p.current].Type == TokenEOF {
		return true
	}
	return false
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.end() {
		if p.previous().Type == TokenSemi {
			return
		}

		switch p.peek().Type {
		case TokenClass,
			TokenFun,
			TokenVar,
			TokenFor,
			TokenIf,
			TokenWhile,
			TokenPrint,
			TokenReturn:
			return
		default:
			break
		}

		p.advance()
	}
}

// program			-> declaration* EOF ;
// declaration		-> varDeclaration | funDeclaration ;
// funDeclaration	-> "fun" function ;
// function			-> IDENTIFIER "(" parameters? ")" block ;
// parameters		-> IDENTIFIER ( "," IDENTIFIER )* ;
// varDeclaration	-> "var" nameDeclaration ;
// nameDeclaration	-> IDENTIFIER ( "=" expression ) ( "," IDENTIFIER ( "=" expression )? )* ";"
// statement		-> block | expreStmt | printStmt | "break" | returnStmt ;
// block			-> "{" declaration* "}" ;
// printStmt		-> "print" expression ;
// expreStmt		-> expression ;
// forStmt			-> "for" "(" ( varDeclaration | expreStmt | ";" ) expression? ";" expression? ")" statement ;
// IfStmt			-> "if" "(" expression ")" statement ( "else" statement  )? ;
// returnStmt		-> "return" expression? ";" ;
// WhileStmt		-> "while" "(" expression ")" statement
// expression		-> assignment ;
// asignment		-> identifier ( "=" | "+=" | "-=" | "*=" | "/=" ) expression | logical_or ;
// logical_or		-> logical_and ( "or" logical_and )* ;
// logical_and		-> equality ( "and" equality )* ;
// equality			-> comparison ( ( "==" | "!=" ) comparison )* ;
// comparison		-> addition ( ( "<" | "<=" | ">" | ">=" ) addition )* ;
// addition			-> multiplication ( ( "+" | "-" ) multiplication )* ;
// multiplication 	-> unary ( ( "*" | "/" | "%" ) unary )* ;
// unary			-> ( "!" | "-" ) unary | call ;
// call				-> primary ( "(" expression ( "," expression )* "}" ) ;
// primary 			-> IDENTIFIER | NUMBER | STRING | "(" expression ")" | "true" | "false" | "nil" ;

// Parse is the entry point of Parser.
func (p *Parser) Parse() ([]Stmt, bool) {
	var stmts []Stmt
	for !p.end() {
		stmts = append(stmts, p.declaration())
	}

	return stmts, p.hadError
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if val := recover(); val != nil {
			// might trigger another panic if it is not a Parsing Error.
			parsingError := val.(*LoxError)
			fmt.Println(parsingError.Error())

			p.hadError = true
			p.synchronize()
		}
	}()

	switch {
	case p.match(TokenVar):
		return p.varDeclaration()
	case p.match(TokenFun):
		return p.function("function")
	default:
		return p.statement()
	}
}

func (p *Parser) function(kind string) Stmt {
	var (
		name   *Token
		params []*Token
		body   []Stmt
	)

	name = p.consume(TokenIdentifier, "expect IDENTIFIER after 'fun'.")
	params = make([]*Token, 0)

	// parameters.
	p.consume(TokenLeftParen, "expect '(' after IDENTIFIER.")
	if !p.check(TokenRightParen) {
		for true {
			param := p.consume(TokenIdentifier, "expect TokenIdentifier as param.")
			params = append(params, param)

			if len(params) > 8 {
				panic(NewLoxError(p.peek(), "cannot have more than 8 parameters."))
			}
			if !p.match(TokenComma) {
				break
			}
		}
	}
	p.consume(TokenRightParen, "expect ')' after param list.")

	// body.
	p.consume(TokenLeftBrace, "expect '{' before function body.")
	body = p.block()
	return NewFunction(name, params, body)
}

func (p *Parser) varDeclaration() Stmt {
	varDec := p.nameDeclaration()

	if p.check(TokenSemi) {
		p.advance()
		return varDec
	}

	varDecs := make([]*Var, 0)
	varDecs = append(varDecs, varDec)
	for !p.check(TokenSemi) {
		p.consume(TokenComma, "unexpected token.")
		varDec = p.nameDeclaration()
		varDecs = append(varDecs, varDec)
	}
	p.consume(TokenSemi, "expect ';' after variable declaration.")

	return NewVarList(varDecs)
}

// a helper function for dealing with multi-var declarations.
func (p *Parser) nameDeclaration() *Var {
	var (
		name        *Token
		initializer Expr
	)

	name = p.consume(TokenIdentifier, "expect variable name.")
	if p.match(TokenEqual) {
		initializer = p.expression()
	}

	// We convert it eagerly because we know it's a *Var.
	varDec, _ := NewVar(name, initializer).(*Var)
	return varDec
}

func (p *Parser) statement() Stmt {
	switch {
	case p.match(TokenBreak):
		keyword := p.previous()
		p.consume(TokenSemi, "expect ';' after 'break'.")
		return NewControl(keyword, ControlBreak, nil)
	case p.match(TokenFor):
		return p.forStmt()
	case p.match(TokenIf):
		return p.ifStmt()
	case p.match(TokenPrint):
		return p.printStmt()
	case p.match(TokenReturn):
		var value Expr

		keyword := p.previous()
		if !p.check(TokenSemi) {
			value = p.expression()
		}
		p.consume(TokenSemi, "expect ';' after return value.")
		return NewControl(keyword, ControlReturn, value)
	case p.match(TokenLeftBrace):
		return NewBlock(p.block())
	case p.match(TokenWhile):
		return p.whileStmt()
	default:
		return p.expressionStmt()
	}
}

func (p *Parser) block() []Stmt {
	stmts := make([]Stmt, 0)

	for !p.check(TokenRightBrace) && !p.end() {
		stmts = append(stmts, p.declaration())
	}

	p.consume(TokenRightBrace, "expect '}' after block.")
	return stmts
}

/*
 for statement is a syntax sugar for while statement.

 		for (var i = 0; i < 10; i = i + 1) {
 				print i;
 		}

 will be translated into

		{
			var i = 0;
			while (i < 10) {
				print i;
				i = i + 1
			}
		}
*/
func (p *Parser) forStmt() Stmt {
	var (
		forBlock    Stmt
		initializer Stmt
		condition   Expr
		increment   Expr
		body        Stmt
		forBody     []Stmt
	)

	p.consume(TokenLeftParen, "expect '(' after 'for'.")
	if p.match(TokenVar) {
		initializer = p.varDeclaration()
	} else if p.match(TokenSemi) {
		initializer = nil
	} else {
		initializer = p.expressionStmt()
	}

	if p.check(TokenSemi) {
		condition = NewLiteral(true)
	} else {
		condition = p.expression()
	}

	p.consume(TokenSemi, "expect ';' after 'for' initializer.")
	if p.check(TokenRightParen) {
		increment = nil
	} else {
		increment = p.expression()
	}
	p.match(TokenRightParen)

	body = p.statement()
	// If it is a block, we strip it first.
	if stmts, ok := body.(*Block); ok {
		forBody = stmts.Stmts
	} else {
		forBody = make([]Stmt, 2)
		forBody = append(forBody, body)
	}

	if increment != nil {
		forBody = append(forBody, NewExpression(increment))
	}

	innerWhile := NewWhile(condition, NewBlock(forBody))
	if initializer != nil {
		forBlock = NewBlock([]Stmt{initializer, innerWhile})
	} else {
		forBlock = NewBlock([]Stmt{innerWhile})
	}
	return forBlock
}

func (p *Parser) ifStmt() Stmt {
	var (
		condition  Expr
		thenBranch Stmt
		elseBranch Stmt
	)

	p.consume(TokenLeftParen, "expect '(' after 'if'.")
	condition = p.expression()
	p.consume(TokenRightParen, "expect ')' after if condition.")

	thenBranch = p.statement()
	if p.check(TokenElse) {
		p.advance()
		elseBranch = p.statement()
	}

	return NewIf(condition, thenBranch, elseBranch)
}

func (p *Parser) printStmt() Stmt {
	expr := p.expression()
	p.consume(TokenSemi, "expect ';' after print expression.")
	return NewPrint(expr)
}

func (p *Parser) expressionStmt() Stmt {
	expr := p.expression()

	p.consume(TokenSemi, "expect ';' after expression.")
	return NewExpression(expr)
}

func (p *Parser) whileStmt() Stmt {
	var (
		condition Expr
		body      Stmt
	)

	p.consume(TokenLeftParen, "expect '(' after 'while'.")
	condition = p.expression()
	p.consume(TokenRightParen, "expect ')' after while condition.")

	body = p.statement()
	return NewWhile(condition, body)
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	// call assignment recursively because it's right associative.
	if p.match(TokenEqual, TokenPlusEqual, TokenMinusEqual, TokenStarEqual, TokenSlashEqual) {
		operator := p.previous()
		value := p.assignment()

		if varExpr, ok := expr.(*Variable); ok {
			name := varExpr.Name
			return NewAssign(name, operator, value)
		}

		errmsg := "invalid assign target."
		panic(NewLoxError(operator, errmsg))
	}
	// let's skip it until we define variable.
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(TokenOr) {
		operator := p.previous()
		right := p.and()
		expr = NewLogical(expr, operator, right)
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(TokenAnd) {
		operator := p.previous()
		right := p.equality()
		expr = NewLogical(expr, operator, right)
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(TokenEqualEqual, TokenBangEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.addition()

	if p.match(TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) addition() Expr {
	expr := p.multiplication()

	for p.match(TokenPlus, TokenMinus) {
		operator := p.previous()
		right := p.multiplication()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) multiplication() Expr {
	expr := p.unary()

	for p.match(TokenStar, TokenSlash, TokenPercent) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(TokenBang, TokenMinus) {
		operator := p.previous()
		value := p.unary()
		return NewUnary(operator, value)
	}

	return p.call()
}

func (p *Parser) call() Expr {
	var (
		expr      Expr
		arguments []Expr
		paren     *Token
	)

	expr = p.primary()

	if p.check(TokenLeftParen) {
		p.advance()
		paren = p.previous()
		arguments = p.arguments()
		expr = NewCall(expr, paren, arguments)
	}

	return expr
}

func (p *Parser) arguments() []Expr {
	exprs := make([]Expr, 0)

	for !p.check(TokenRightParen) {
		expr := p.expression()
		exprs = append(exprs, expr)
		if p.check(TokenRightParen) {
			break
		}
		p.consume(TokenComma, "expect ';' to separate arguments.")
	}
	p.consume(TokenRightParen, "expect ')' after argument list.")
	return exprs
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(TokenIdentifier):
		return NewVariable(p.previous())
	case p.match(TokenFalse):
		return NewLiteral(false)
	case p.match(TokenTrue):
		return NewLiteral(true)
	case p.match(TokenNil):
		return NewLiteral(nil)
	case p.match(TokenNumber, TokenString):
		return NewLiteral(p.previous().Literal)
	case p.match(TokenLeftParen): // grouping.
		expr := p.expression()
		p.consume(TokenRightParen, "expect ')' after expression.")
		return NewGrouping(expr)
	default:
		panic(NewLoxError(p.peek(), "expect expression."))
	}
}
