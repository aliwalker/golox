package golox

// Parser parses the tokens into an AST.
type Parser struct {
	tokens  []*Token
	current int
}

// NewParser creates a parser.
func NewParser(tokens []*Token) *Parser {
	return &Parser{tokens, 0}
}

// when there's parsing error that cannot be continue, parser should panic.
func parserPanic(token *Token, message string) error {
	ParsingError(token, message)
	panic(message)
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
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
	parserPanic(p.peek(), message)
	return nil
}

func (p *Parser) isAtEnd() bool {
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

// program			-> statement* EOF ;
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
func (p *Parser) Parse() []Stmt {
	var stmts []Stmt
	for !p.isAtEnd() {
		stmts = append(stmts, p.statement())
	}

	return stmts
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
		parserPanic(p.peek(), "expect expression.")
		return nil
	}
}
