package parsing

import (
	"monkey/lexing"
	"strconv"
)

const (
	_ = iota
	PRECEDENCE_LOWEST
	PRECEDENCE_EQUALS
	PRECEDENCE_LESS_GREATER
	PRECEDENCE_SUM
	PRECEDENCE_PRODUCT
	PRECEDENCE_PREFIX
)

func getPrecedence(tokenType lexing.TokenType) int {
	tokenTypeToPrecedence := map[lexing.TokenType]int{
		lexing.TOKEN_EQUALS:            PRECEDENCE_EQUALS,
		lexing.TOKEN_NOT_EQUALS:        PRECEDENCE_EQUALS,
		lexing.TOKEN_GREATER:           PRECEDENCE_LESS_GREATER,
		lexing.TOKEN_GREATER_OR_EQUALS: PRECEDENCE_LESS_GREATER,
		lexing.TOKEN_LESS:              PRECEDENCE_LESS_GREATER,
		lexing.TOKEN_LESS_OR_EQUALS:    PRECEDENCE_LESS_GREATER,
		lexing.TOKEN_PLUS:              PRECEDENCE_SUM,
		lexing.TOKEN_MINUS:             PRECEDENCE_SUM,
		lexing.TOKEN_ASTERISK:          PRECEDENCE_PRODUCT,
		lexing.TOKEN_SLASH:             PRECEDENCE_PRODUCT,
	}
	return tokenTypeToPrecedence[tokenType]
}

type Parser struct {
	tokens   []*lexing.Token
	position int
	current  *lexing.Token
}

func NewParser(lexer *lexing.Lexer) *Parser {
	var tokens []*lexing.Token

	for {
		token := lexer.Next()
		tokens = append(tokens, token)
		if token.Type == lexing.TOKEN_EOF {
			break
		}
	}

	return &Parser{
		tokens:   tokens,
		position: 0,
		current:  tokens[0],
	}
}

func (parser *Parser) advance() {
	if parser.position+1 >= len(parser.tokens) {
		parser.current = nil
		return
	}
	parser.position += 1
	parser.current = parser.tokens[parser.position]
}

func (parser *Parser) parseIntegerLiteral() *AstIntegerLiteral {
	value, _ := strconv.ParseInt(parser.current.Literal, 10, 64)
	integerLiteral := &AstIntegerLiteral{
		Token: parser.current,
		Value: value,
	}
	parser.advance()
	return integerLiteral
}

func (parser *Parser) parsePrefixExpression() *AstPrefixExpression {
	prefixExpresion := &AstPrefixExpression{
		Token:    parser.current,
		Operator: parser.current.Literal,
	}
	parser.advance()
	prefixExpresion.Right = parser.parseExpression(PRECEDENCE_PREFIX)
	return prefixExpresion
}

func (parser *Parser) parseInfixExpression(left AstExpression) *AstInfixExpression {
	infixExpression := &AstInfixExpression{
		Token:    parser.current,
		Left:     left,
		Operator: parser.current.Literal,
	}
	precedence := getPrecedence(parser.current.Type)
	parser.advance()
	infixExpression.Right = parser.parseExpression(precedence)
	return infixExpression
}

func (parser *Parser) parseBooleanLiteral() *AstBooleanLiteral {
	booleanLiteral := &AstBooleanLiteral{
		Token: parser.current,
		Value: parser.current.Type == lexing.TOKEN_TRUE,
	}
	parser.advance()
	return booleanLiteral
}

func (parser *Parser) parseEnforcedPrecedenceExpression() AstExpression {
	parser.advance()
	expression := parser.parseExpression(PRECEDENCE_LOWEST)
	parser.advance()
	return expression
}

func (parser *Parser) parseExpression(precedence int) AstExpression {
	var left AstExpression

	switch parser.current.Type {
	case lexing.TOKEN_INTEGER:
		left = parser.parseIntegerLiteral()
	case lexing.TOKEN_TRUE, lexing.TOKEN_FALSE:
		left = parser.parseBooleanLiteral()
	case lexing.TOKEN_BANG, lexing.TOKEN_MINUS:
		left = parser.parsePrefixExpression()
	case lexing.TOKEN_OPEN_PAREN:
		left = parser.parseEnforcedPrecedenceExpression()
	}

	for precedence < getPrecedence(parser.current.Type) {
		left = parser.parseInfixExpression(left)
	}

	return left
}

func (parser *Parser) parseExpressionStatement() *AstExpressionStatement {
	return &AstExpressionStatement{
		Token:      parser.current,
		Expression: parser.parseExpression(PRECEDENCE_LOWEST),
	}
}

func (parser *Parser) parseStatement() AstStatement {
	switch parser.current.Type {
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) Parse() AstNode {
	return parser.parseStatement()
}
