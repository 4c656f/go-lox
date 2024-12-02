package parser

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/expression"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/stmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type ASTPrinter struct {
	outString string
}

func (a *ASTPrinter) PrintProgram(st []stmt.Stmt) string {
	var b strings.Builder
	if st == nil {
		return ""
	}

	for _, s := range st {
		s.Accept(a)
		b.WriteString(a.Out())
	}
	return b.String()
}

func (a *ASTPrinter) Print(exp expression.Expression) string {
	if exp == nil {
		return ""
	}
	exp.Accept(a)
	return a.Out()
}

func (a *ASTPrinter) Out() string {
	return a.outString
}

func (a *ASTPrinter) VisitExpressionStmt(s *stmt.ExpressionStmt) {
	a.outString = a.parenthesize("stmt", s.Exp)
}

func (a *ASTPrinter) VisitBlockStmt(s *stmt.BlockStmt) {
	var inside strings.Builder
	for _, st := range s.Statements {
		st.Accept(a)
		inside.WriteString(a.outString)
	}
	a.outString = fmt.Sprintf("{ %s }", inside.String())
}

func (a *ASTPrinter) VisitPrintStmt(s *stmt.PrintStmt) {
	a.outString = a.parenthesize("print", s.Exp)
}

func (a *ASTPrinter) VisitIfStmt(s *stmt.IfStmt) {
	s.ThenBranch.Accept(a)
	a.outString = fmt.Sprintf("%s, {\n%s\n}", a.parenthesize("if", s.Condition), a.Out())
}

func (a *ASTPrinter) VisitVarStmt(s *stmt.VarStmt) {
	a.outString = a.parenthesize("var = ", s.Init)
}

func (a *ASTPrinter) VisitBinary(b *expression.BinaryExpression) {
	a.outString = a.parenthesize(b.Op.Text, b.Lhs, b.Rhs)
}

func (a *ASTPrinter) VisitGrouping(g *expression.GroupingExpression) {
	a.outString = a.parenthesize("group", g.Exp)
}

func (a *ASTPrinter) VisitLiteral(l *expression.LiteralExpression) {
	if l.Val.TokenValue.Type == token.BoolValue || l.Val.TokenValue.Type == token.NullValue {
		a.outString = l.Val.Text
		return
	}
	a.outString = l.Val.TokenValue.String()
}

func (a *ASTPrinter) VisitUnary(u *expression.UnaryExpression) {
	a.outString = a.parenthesize(u.Op.Text, u.Rhs)
}

func (a *ASTPrinter) VisitVarExpression(u *expression.VarExpression) {
	a.outString = fmt.Sprintf("var %s", u.Name.Text)
}

func (a *ASTPrinter) VisitAssignmentExpression(u *expression.AssignmentExpression) {
	a.outString = fmt.Sprintf("ass %s", a.parenthesize(u.Name.Text, u.Val))
}

// Helper function to create parenthesized expressions
func (a *ASTPrinter) parenthesize(name string, exprs ...expression.Expression) string {
	var result strings.Builder
	result.WriteString("(" + name)
	for _, expr := range exprs {
		result.WriteString(" ")
		expr.Accept(a)
		result.WriteString(a.Out())
	}
	result.WriteString(")")
	return result.String()
}

func NewAstPrinter() *ASTPrinter {
	return &ASTPrinter{}
}

type Parser struct {
	tokens []*token.Token
	errors []error
	cur    int
}

func (p *Parser) Parse() (expression.Expression, []error) {
	return p.expression(), p.errors
}

func (p *Parser) ParseProgram() ([]stmt.Stmt, []error) {
	statements := []stmt.Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements, p.errors
}

func (p *Parser) declaration() stmt.Stmt {
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) statement() stmt.Stmt {
	if p.match(token.IF) {
		return p.ifStmt()
	}
	if p.match(token.PRINT) {
		return p.printStmt()
	}
	if p.match(token.LEFT_BRACE) {
		return p.blockStmt()
	}
	return p.expStmt()
}

func (p *Parser) ifStmt() stmt.Stmt {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		p.onError(err)
		return nil
	}
	condition := p.expression()
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		p.onError(err)
		return nil
	}
	flow := p.statement()
	var elseStmt stmt.Stmt
	if p.match(token.ELSE) {
		elseStmt = p.statement()
	}
	return stmt.NewIfStmt(condition, flow, elseStmt)
}

func (p *Parser) varDeclaration() stmt.Stmt {
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		p.onError(err)
		return nil
	}
	var initializer expression.Expression

	if p.match(token.EQUAL) {
		initializer = p.expression()
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")

	if err != nil {
		p.onError(err)
		return nil
	}
	return stmt.NewVarStmt(name, initializer)
}

func (p *Parser) printStmt() stmt.Stmt {
	value := p.expression()
	_, err := p.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		p.onError(err)
		return nil
	}
	return stmt.NewPrintStmt(value)
}

func (p *Parser) blockStmt() stmt.Stmt {
	statemnts := []stmt.Stmt{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statemnts = append(statemnts, p.declaration())
	}
	_, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		p.onError(err)
		return nil
	}
	return stmt.NewBlockStmt(statemnts)
}

func (p *Parser) expStmt() stmt.Stmt {
	exp := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return stmt.NewExpressionStmt(exp)
}

func (p *Parser) expression() expression.Expression {
	return p.assignment()
}

func (p *Parser) assignment() expression.Expression {
	exp := p.equality()
	if p.match(token.EQUAL) {
		equals := p.prev()
		value := p.assignment()
		name, ok := exp.(*expression.VarExpression)
		if !ok {
			p.onError(NewParserError(equals, "Invalid assignment target."))
			return exp
		}
		return expression.NewAssignmentExprExpression(name.Name, value)
	}
	return exp
}

func (p *Parser) equality() expression.Expression {
	exp := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		op := p.prev()
		rhs := p.comparison()
		exp = expression.NewBinaryExpression(exp, op, rhs)
	}

	return exp
}

func (p *Parser) comparison() expression.Expression {
	exp := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS_EQUAL, token.LESS) {
		op := p.prev()
		rhs := p.term()
		exp = expression.NewBinaryExpression(exp, op, rhs)
	}

	return exp
}

func (p *Parser) term() expression.Expression {
	exp := p.factor()
	for p.match(token.MINUS, token.PLUS) {
		op := p.prev()
		rhs := p.factor()
		exp = expression.NewBinaryExpression(exp, op, rhs)
	}
	return exp
}

func (p *Parser) factor() expression.Expression {
	exp := p.unary()
	for p.match(token.SLASH, token.STAR) {
		op := p.prev()
		rhs := p.unary()
		exp = expression.NewBinaryExpression(exp, op, rhs)
	}
	return exp
}

func (p *Parser) unary() expression.Expression {
	if p.match(token.BANG, token.MINUS) {
		op := p.prev()
		rhs := p.unary()
		return expression.NewUnaryExpression(op, rhs)
	}
	return p.primary()
}

func (p *Parser) primary() expression.Expression {
	if p.match(token.TRUE, token.FALSE, token.NIL, token.NUMBER, token.STRING) {
		return expression.NewLiteralExpression(p.prev())
	}
	if p.match(token.IDENTIFIER) {
		return expression.NewVarExpression(p.prev())
	}

	if p.match(token.LEFT_PAREN) {
		exp := p.expression()
		_, err := p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			p.onError(err)
			return nil
		}
		return expression.NewGroupingExpression(exp)
	}
	p.onError(NewParserError(p.peek(), "Expect expression."))
	return nil
}

func (p *Parser) consume(t token.TokenType, message string) (*token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return nil, NewParserError(p.peek(), message)
}

func (p *Parser) onError(err error) {
	p.errors = append(p.errors, err)
	p.sync()
}

func (p *Parser) match(tokens ...token.TokenType) bool {
	for _, t := range tokens {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) sync() {
	p.advance()

	for !p.isAtEnd() {
		if p.prev().Type == token.SEMICOLON {
			return
		}
		switch p.peek().Type {
		case token.CLASS:
			return
		case token.FUN:
			return
		case token.VAR:
			return
		case token.FOR:
			return
		case token.IF:
			return
		case token.WHILE:
			return
		case token.PRINT:
			return
		case token.RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.cur++
	}
	return p.prev()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.cur]
}

func (p *Parser) prev() *token.Token {
	if p.cur == 0 {
		return nil
	}
	return p.tokens[p.cur-1]
}

func New(tokens []*token.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}
