package parsing

import (
	"fmt"
	"monkey/lexing"
)

const (
	_ = iota
	AST_COMPOUND
	AST_EXPRESSION_STATEMENT
	AST_LET_STATEMENT
	AST_RETURN_STATEMENT
	AST_INTEGER_LITERAL
	AST_BOOLEAN_LITERAL
	AST_PREFIX_EXPRESSION
	AST_INFIX_EXPRESSION
	AST_IDENTIFIER
	AST_FUNCTION_DEFINITION
	AST_FUNCTION_CALL
	AST_ARRAY_LITERAL
	AST_STRING_LITERAL
	AST_HASH_LITERAL
	AST_INDEX
	AST_IF_ELSE
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

type AstCompound struct {
	Token      *lexing.Token
	Statements []AstStatement
}

func (compound *AstCompound) Type() AstType {
	return AST_COMPOUND
}
func (compound *AstCompound) TokenLiteral() string {
	return compound.Token.Literal
}
func (compound *AstCompound) String() string {
	text := ""
	for index, statement := range compound.Statements {
		text += statement.String()
		if index < len(compound.Statements)-1 {
			text += " "
		}
	}
	return text
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

type AstLetStatement struct {
	Token      *lexing.Token
	Identifier *AstIdentifier
	Value      AstExpression
}

func (letStatement *AstLetStatement) statement() {}
func (letStatement *AstLetStatement) Type() AstType {
	return AST_LET_STATEMENT
}
func (letStatement *AstLetStatement) TokenLiteral() string {
	return letStatement.Token.Literal
}
func (letStatement *AstLetStatement) String() string {
	text := letStatement.TokenLiteral() +
		" " +
		letStatement.Identifier.String()
	if letStatement.Value != nil {
		text += " = " + letStatement.Value.String()
	}
	text += ";"

	return text
}

type AstReturnStatement struct {
	Token *lexing.Token
	Value AstExpression
}

func (returnStatement *AstReturnStatement) statement() {}
func (returnStatement *AstReturnStatement) Type() AstType {
	return AST_RETURN_STATEMENT
}
func (returnStatement *AstReturnStatement) TokenLiteral() string {
	return returnStatement.Token.Literal
}
func (returnStatement *AstReturnStatement) String() string {
	text := returnStatement.TokenLiteral()
	if returnStatement.Value != nil {
		text += " " + returnStatement.Value.String()
	}
	text += ";"

	return text
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

type AstFunctionCall struct {
	Token     *lexing.Token
	Left      AstExpression
	Arguments []AstExpression
}

func (functionCall *AstFunctionCall) expression() {}
func (functionCall *AstFunctionCall) Type() AstType {
	return AST_FUNCTION_CALL
}
func (functionCall *AstFunctionCall) TokenLiteral() string {
	return functionCall.Token.Literal
}
func (functionCall *AstFunctionCall) String() string {
	text := functionCall.Left.String() + "("

	for index, argument := range functionCall.Arguments {
		text += argument.String()
		if index < len(functionCall.Arguments)-1 {
			text += ", "
		}
	}
	text += ")"

	return text
}

type AstIndex struct {
	Token *lexing.Token
	Left  AstExpression
	Index AstExpression
}

func (index *AstIndex) expression() {}
func (index *AstIndex) Type() AstType {
	return AST_INDEX
}
func (index *AstIndex) TokenLiteral() string {
	return index.Token.Literal
}
func (index *AstIndex) String() string {
	return index.Left.String() + "[" + index.Index.String() + "]"
}

type AstStringLiteral struct {
	Token *lexing.Token
	Value string
}

func (stringLiteral *AstStringLiteral) expression() {}
func (stringLiteral *AstStringLiteral) Type() AstType {
	return AST_STRING_LITERAL
}
func (stringLiteral *AstStringLiteral) TokenLiteral() string {
	return stringLiteral.Token.Literal
}
func (stringLiteral *AstStringLiteral) String() string {
	return "\"" + stringLiteral.Value + "\""
}

type AstArrayLiteral struct {
	Token *lexing.Token
	Items []AstExpression
}

func (arrayLiteral *AstArrayLiteral) expression() {}
func (arrayLiteral *AstArrayLiteral) Type() AstType {
	return AST_ARRAY_LITERAL
}
func (arrayLiteral *AstArrayLiteral) TokenLiteral() string {
	return arrayLiteral.Token.Literal
}
func (arrayLiteral *AstArrayLiteral) String() string {
	text := "["
	for index, item := range arrayLiteral.Items {
		text += item.String()
		if index < len(arrayLiteral.Items)-1 {
			text += ", "
		}
	}
	text += "]"

	return text
}

type AstHashLiteralPair struct {
	Key   AstExpression
	Value AstExpression
}

type AstHashLiteral struct {
	Token *lexing.Token
	Pairs []*AstHashLiteralPair
}

func (hashLiteral *AstHashLiteral) expression() {}
func (hashLiteral *AstHashLiteral) Type() AstType {
	return AST_HASH_LITERAL
}
func (hashLiteral *AstHashLiteral) TokenLiteral() string {
	return hashLiteral.Token.Literal
}
func (hashLiteral *AstHashLiteral) String() string {
	text := "{"
	for index, pair := range hashLiteral.Pairs {
		text += pair.Key.String() + ": " + pair.Value.String()
		if index < len(hashLiteral.Pairs)-1 {
			text += ", "
		}
	}
	text += "}"

	return text

}

type AstFunctionDefinition struct {
	Token      *lexing.Token
	Parameters []*AstIdentifier
	Body       *AstCompound
}

func (functionDefinition *AstFunctionDefinition) expression() {}
func (functionDefinition *AstFunctionDefinition) Type() AstType {
	return AST_FUNCTION_DEFINITION
}
func (functionDefinition *AstFunctionDefinition) TokenLiteral() string {
	return functionDefinition.Token.Literal
}
func (functionDefinition *AstFunctionDefinition) String() string {
	text := functionDefinition.TokenLiteral() + " ("

	for index, parameter := range functionDefinition.Parameters {
		text += parameter.String()
		if index < len(functionDefinition.Parameters)-1 {
			text += ", "
		}
	}

	text += ") { " + functionDefinition.Body.String() + " }"

	return text
}

type AstIfElse struct {
	Token     *lexing.Token
	Condition AstExpression
	Then      *AstCompound
	Else      *AstCompound
}

func (ifElse *AstIfElse) expression() {}
func (ifElse *AstIfElse) Type() AstType {
	return AST_IF_ELSE
}
func (ifElse *AstIfElse) TokenLiteral() string {
	return ifElse.Token.Literal
}
func (ifElse *AstIfElse) String() string {
	text := ifElse.TokenLiteral() +
		" (" +
		ifElse.Condition.String() +
		") { " +
		ifElse.Then.String() +
		" }"

	if ifElse.Else != nil {
		text += " else { " + ifElse.Else.String() + " }"
	}
	return text
}
