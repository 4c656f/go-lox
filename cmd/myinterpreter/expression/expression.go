package expression

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"

type Visitor interface {
	VisitBinary(b *BinaryExpression)
	VisitGrouping(g *GroupingExpression)
	VisitUnary(u *UnaryExpression)
	VisitLiteral(u *LiteralExpression)
	VisitVarExpression(u *VarExpression)
	VisitAssignmentExpression(u *AssignmentExpression)
	VisitLogicalExpression(u *LogicalExpression)
	VisitFunctionCallExpression(u *FunctionCallExpression)
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

type VarExpression struct {
	Name *token.Token
}

type AssignmentExpression struct {
	Name *token.Token
	Val  Expression
}

type LogicalExpression struct {
	Lhs Expression
	Op  *token.Token
	Rhs Expression
}

type FunctionCallExpression struct {
	Callee      Expression
	Args       []Expression
	RightParan *token.Token
}

func (this *FunctionCallExpression) Accept(v Visitor) {
    v.VisitFunctionCallExpression(this)
}

func (this *BinaryExpression) Accept(v Visitor) {
	v.VisitBinary(this)
}

func (this *AssignmentExpression) Accept(v Visitor) {
	v.VisitAssignmentExpression(this)
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

func (this *VarExpression) Accept(v Visitor) {
	v.VisitVarExpression(this)
}

func (this *LogicalExpression) Accept(v Visitor) {
	v.VisitLogicalExpression(this)
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

func NewVarExpression(
	name *token.Token,
) *VarExpression {
	return &VarExpression{
		Name: name,
	}
}

func NewAssignmentExprExpression(
	name *token.Token,
	value Expression,
) *AssignmentExpression {
	return &AssignmentExpression{
		Name: name,
		Val:  value,
	}
}

func NewLogicalExpression(
	lhs Expression,
	op *token.Token,
	rhs Expression,
) *LogicalExpression {
	return &LogicalExpression{
		Lhs: lhs,
		Op:  op,
		Rhs: rhs,
	}
}

func NewFunctionCallExpression(
    callee Expression,
    args []Expression,
    rightParan *token.Token,
) *FunctionCallExpression {
    return &FunctionCallExpression{
        Callee:      callee,
        Args:       args,
        RightParan: rightParan,
    }
}
