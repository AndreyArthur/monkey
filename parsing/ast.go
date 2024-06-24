package parsing

import (
	"fmt"
	"monkey/lexing"
)

const (
	_ = iota
	AST_EXPRESSION_STATEMENT
	AST_INTEGER_LITERAL
	AST_BOOLEAN_LITERAL
	AST_PREFIX_EXPRESSION
	AST_INFIX_EXPRESSION
	AST_IDENTIFIER
)

type AstType int

type AstNode interface {
	Type() AstType
	TokenLiteral() string
	String() string
}

type AstStatement interface {
	AstNode
	statement()
}

type AstExpression interface {
	AstNode
	expression()
}

type AstExpressionStatement struct {
	Token      *lexing.Token
	Expression AstExpression
}

func (expressionStatement *AstExpressionStatement) statement() {}
func (expressionStatement *AstExpressionStatement) Type() AstType {
	return AST_EXPRESSION_STATEMENT
}
func (expressionStatement *AstExpressionStatement) TokenLiteral() string {
	return expressionStatement.Token.Literal
}
func (expressionStatement *AstExpressionStatement) String() string {
	return expressionStatement.Expression.String() + ";"
}

type AstIntegerLiteral struct {
	Token *lexing.Token
	Value int64
}

func (integerLiteral *AstIntegerLiteral) expression() {}
func (integerLiteral *AstIntegerLiteral) Type() AstType {
	return AST_INTEGER_LITERAL
}
func (integerLiteral *AstIntegerLiteral) TokenLiteral() string {
	return integerLiteral.Token.Literal
}
func (integerLiteral *AstIntegerLiteral) String() string {
	return fmt.Sprintf("%d", integerLiteral.Value)
}

type AstBooleanLiteral struct {
	Token *lexing.Token
	Value bool
}

func (booleanLiteral *AstBooleanLiteral) expression() {}
func (booleanLiteral *AstBooleanLiteral) Type() AstType {
	return AST_BOOLEAN_LITERAL
}
func (booleanLiteral *AstBooleanLiteral) TokenLiteral() string {
	return booleanLiteral.Token.Literal
}
func (booleanLiteral *AstBooleanLiteral) String() string {
	return fmt.Sprintf("%t", booleanLiteral.Value)
}

type AstPrefixExpression struct {
	Token    *lexing.Token
	Operator string
	Right    AstExpression
}

func (prefixExpression *AstPrefixExpression) expression() {}
func (prefixExpression *AstPrefixExpression) Type() AstType {
	return AST_PREFIX_EXPRESSION
}
func (prefixExpression *AstPrefixExpression) TokenLiteral() string {
	return prefixExpression.Token.Literal
}
func (prefixExpression *AstPrefixExpression) String() string {
	return "(" +
		prefixExpression.Operator +
		prefixExpression.Right.String() +
		")"
}

type AstInfixExpression struct {
	Token    *lexing.Token
	Left     AstExpression
	Operator string
	Right    AstExpression
}

func (infixExpression *AstInfixExpression) expression() {}
func (infixExpression *AstInfixExpression) Type() AstType {
	return AST_INFIX_EXPRESSION
}
func (infixExpression *AstInfixExpression) TokenLiteral() string {
	return infixExpression.Token.Literal
}
func (infixExpression *AstInfixExpression) String() string {
	return "(" +
		infixExpression.Left.String() +
		" " +
		infixExpression.Operator +
		" " +
		infixExpression.Right.String() +
		")"
}

type AstIdentifier struct {
	Token *lexing.Token
	Name  string
}

func (identifier *AstIdentifier) expression() {}
func (identifier *AstIdentifier) Type() AstType {
	return AST_IDENTIFIER
}
func (identifier *AstIdentifier) TokenLiteral() string {
	return identifier.Token.Literal
}
func (identifier *AstIdentifier) String() string {
	return identifier.Name
}
