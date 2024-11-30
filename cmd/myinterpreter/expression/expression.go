package expression

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"

type Visitor interface {
	VisitBinary(b *BinaryExpression)
	VisitGrouping(g *GroupingExpression)
	VisitUnary(u *UnaryExpression)
	VisitLiteral(u *LiteralExpression)
}

type Expression interface {
	Accept(v Visitor)
}

type BinaryExpression struct {
	Lhs Expression
	Op  *token.Token
	Rhs Expression
}

type GroupingExpression struct {
	Exp Expression
}

type LiteralExpression struct {
	Val *token.Token
}

type UnaryExpression struct {
	Op  *token.Token
	Rhs Expression
}

func (this *BinaryExpression) Accept(v Visitor) {
	v.VisitBinary(this)
}

func (this *GroupingExpression) Accept(v Visitor) {
	v.VisitGrouping(this)
}

func (this *LiteralExpression) Accept(v Visitor) {
	v.VisitLiteral(this)
}

func (this *UnaryExpression) Accept(v Visitor) {
	v.VisitUnary(this)
}

func NewBinaryExpression(
	lhs Expression,
	op *token.Token,
	rhs Expression,
) *BinaryExpression {
	return &BinaryExpression{
		Lhs: lhs,
		Op:  op,
		Rhs: rhs,
	}
}

func NewGroupingExpression(exp Expression) *GroupingExpression {
	return &GroupingExpression{exp}
}

func NewLiteralExpression(val *token.Token) *LiteralExpression {
	return &LiteralExpression{val}
}

func NewUnaryExpression(
	op *token.Token,
	exp Expression,
) *UnaryExpression {
	return &UnaryExpression{
		Op:  op,
		Rhs: exp,
	}
}