package lox

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
	ParsingError(p.peek(), message)
	panic(message)
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
// declaration		-> varDeclaration ;
// varDeclaration	-> "var" IDENTIFIER ( "=" expression )? ";" ;
// statement		-> expreStmt | printStmt ;
// printStmt		-> "print" expression ;
// expreStmt		-> expression ;
// expression		-> assignment ;
// asignment		-> identifier "=" expression | logical_or ;
// logical_or		-> logical_and ( "or" logical_and )* ;
// logical_and		-> equality ( "and" equality )* ;
// equality			-> comparison ( ( "==" | "!=" ) comparison )* ;
// comparison		-> addition ( ( "<" | "<=" | ">" | ">=" ) addition )* ;
// addition			-> multiplication ( ( "+" | "-" ) multiplication )* ;
// multiplication 	-> unary ( ( "*" | "/" | "%" ) unary )* ;
// unary			-> ( "!" | "-" )? primary ;
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
		if r := recover(); r != nil {
			p.hadError = true
			p.synchronize()
		}
	}()

	switch {
	case p.match(TokenVar):
		return p.varDeclaration()
	default:
		return p.statement()
	}
	return nil
}

func (p *Parser) varDeclaration() Stmt {
	var (
		name        *Token
		initializer Expr
	)

	name = p.consume(TokenIdentifier, "expect variable name.")
	if p.match(TokenEqual) {
		initializer = p.expression()
	}

	p.consume(TokenSemi, "expect ';' after variable declaration.")
	return NewVar(name, initializer)
}

func (p *Parser) statement() Stmt {
	switch {
	case p.match(TokenPrint):
		return p.printStmt()
	default:
		return p.expressionStmt()
	}
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

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()

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
	// TODO: call p.call() when function call is implemented.
	return p.primary()
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
		ParsingError(p.peek(), "expect expression.")
		panic("expect expression.")
	}
}
