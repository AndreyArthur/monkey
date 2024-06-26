package parsing

import (
	"fmt"
	"monkey/lexing"
	"slices"
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
	PRECEDENCE_INDEX
	PRECEDENCE_CALL
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
		lexing.TOKEN_OPEN_BRACKET:      PRECEDENCE_INDEX,
		lexing.TOKEN_OPEN_PAREN:        PRECEDENCE_CALL,
	}
	return tokenTypeToPrecedence[tokenType]
}

type Parser struct {
	tokens       []*lexing.Token
	position     int
	current      *lexing.Token
	errors       []string
	currentError string
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

func (parser *Parser) commitError() {
	if parser.currentError == "" {
		return
	}
	parser.errors = append(parser.errors, parser.currentError)
	parser.currentError = ""
	return
}

func (parser *Parser) error(message string) {
	if parser.currentError == "" {
		parser.currentError = message
	}
}

func (parser *Parser) expect(tokenTypes ...lexing.TokenType) {
	if slices.Contains(tokenTypes, parser.current.Type) {
		return
	}
	tokenTypesString := ""
	for index, tokenType := range tokenTypes {
		tokenTypesString += lexing.TokenTypeToString(tokenType)
		if index < len(tokenTypes)-1 {
			tokenTypesString += ", "
		}
	}
	parser.error(fmt.Sprintf(
		"Expected token of type %s. Found token %q of type %s.",
		tokenTypesString,
		parser.current.Literal,
		lexing.TokenTypeToString(parser.current.Type),
	))
}

func (parser *Parser) HasErrors() bool {
	return len(parser.errors) >= 1
}

func (parser *Parser) GetErrors() []string {
	return parser.errors
}

func (parser *Parser) advance() {
	if parser.position+1 >= len(parser.tokens) {
		parser.current = parser.tokens[len(parser.tokens)-1]
		return
	}
	parser.position += 1
	parser.current = parser.tokens[parser.position]
}

func (parser *Parser) peek() *lexing.Token {
	if parser.position+1 >= len(parser.tokens) {
		return nil
	}
	return parser.tokens[parser.position+1]
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

func (parser *Parser) parseIdentifier() *AstIdentifier {
	identifier := &AstIdentifier{
		Token: parser.current,
		Name:  parser.current.Literal,
	}
	parser.advance()
	return identifier
}

func (parser *Parser) parseFunctionCall(left AstExpression) *AstFunctionCall {
	functionCall := &AstFunctionCall{
		Token:     parser.current,
		Left:      left,
		Arguments: []AstExpression{},
	}
	parser.advance()

	for parser.current.Type != lexing.TOKEN_CLOSE_PAREN {
		expression := parser.parseExpression(PRECEDENCE_LOWEST)
		functionCall.Arguments = append(functionCall.Arguments, expression)
		if parser.current.Type != lexing.TOKEN_CLOSE_PAREN {
			parser.expect(lexing.TOKEN_COMMA)
			if parser.current.Type == lexing.TOKEN_COMMA {
				parser.advance()
			} else {
				parser.advance()
				break
			}
		}
	}

	parser.expect(lexing.TOKEN_CLOSE_PAREN)
	parser.advance()

	parser.commitError()
	return functionCall
}

func (parser *Parser) parseIndex(left AstExpression) *AstIndex {
	index := &AstIndex{
		Token: parser.current,
		Left:  left,
	}
	parser.advance()

	index.Index = parser.parseExpression(PRECEDENCE_LOWEST)

	parser.expect(lexing.TOKEN_CLOSE_BRACKET)
	parser.advance()

	parser.commitError()
	return index
}

func (parser *Parser) parseStringLiteral() *AstStringLiteral {
	stringLiteral := &AstStringLiteral{
		Token: parser.current,
		Value: parser.current.Literal,
	}
	parser.advance()
	return stringLiteral
}

func (parser *Parser) parseArrayLiteral() *AstArrayLiteral {
	arrayLiteral := &AstArrayLiteral{
		Token: parser.current,
		Items: []AstExpression{},
	}
	parser.advance()
	for parser.current.Type != lexing.TOKEN_CLOSE_BRACKET {
		expression := parser.parseExpression(PRECEDENCE_LOWEST)
		arrayLiteral.Items = append(arrayLiteral.Items, expression)

		if parser.current.Type != lexing.TOKEN_CLOSE_BRACKET {
			parser.expect(lexing.TOKEN_COMMA)
			if parser.current.Type == lexing.TOKEN_COMMA {
				parser.advance()
			} else {
				parser.advance()
				break
			}
		}
	}

	parser.expect(lexing.TOKEN_CLOSE_BRACKET)
	parser.advance()

	parser.commitError()
	return arrayLiteral
}

func (parser *Parser) parseHashLiteral() *AstHashLiteral {
	hashLiteral := &AstHashLiteral{
		Token: parser.current,
		Pairs: []*AstHashLiteralPair{},
	}
	parser.advance()
	for parser.current.Type != lexing.TOKEN_CLOSE_BRACE {
		key := parser.parseExpression(PRECEDENCE_LOWEST)

		parser.expect(lexing.TOKEN_COLON)
		parser.advance()

		value := parser.parseExpression(PRECEDENCE_LOWEST)

		hashLiteral.Pairs = append(hashLiteral.Pairs, &AstHashLiteralPair{
			Key:   key,
			Value: value,
		})

		if parser.current.Type != lexing.TOKEN_CLOSE_BRACE {
			parser.expect(lexing.TOKEN_COMMA)
			if parser.current.Type == lexing.TOKEN_COMMA {
				parser.advance()
			} else {
				parser.advance()
				break
			}
		}
	}

	parser.expect(lexing.TOKEN_CLOSE_BRACE)
	parser.advance()

	parser.commitError()
	return hashLiteral
}

func (parser *Parser) parseCompound() *AstCompound {
	compound := &AstCompound{
		Token:      parser.current,
		Statements: []AstStatement{},
	}

	for parser.current.Type != lexing.TOKEN_CLOSE_BRACE &&
		parser.current.Type != lexing.TOKEN_EOF {
		compound.Statements = append(
			compound.Statements,
			parser.parseStatement(),
		)
	}

	return compound
}

func (parser *Parser) parseFunctionDefinition() *AstFunctionDefinition {
	functionDefinition := &AstFunctionDefinition{
		Token:      parser.current,
		Parameters: []*AstIdentifier{},
	}

	parser.advance()
	parser.advance()

	for parser.current.Type != lexing.TOKEN_CLOSE_PAREN {
		parser.expect(lexing.TOKEN_IDENTIFIER)
		functionDefinition.Parameters = append(
			functionDefinition.Parameters,
			parser.parseIdentifier(),
		)

		if parser.current.Type != lexing.TOKEN_CLOSE_PAREN {
			parser.expect(lexing.TOKEN_COMMA)
			if parser.current.Type == lexing.TOKEN_COMMA {
				parser.advance()
			} else {
				parser.advance()
				break
			}
		}
	}

	parser.expect(lexing.TOKEN_CLOSE_PAREN)
	parser.advance()

	parser.expect(lexing.TOKEN_OPEN_BRACE)
	parser.advance()

	functionDefinition.Body = parser.parseCompound()

	parser.expect(lexing.TOKEN_CLOSE_BRACE)
	parser.advance()

	parser.commitError()
	return functionDefinition
}

func (parser *Parser) parseIfElse() *AstIfElse {
	ifElse := &AstIfElse{
		Token: parser.current,
	}
	parser.advance()

	parser.expect(lexing.TOKEN_OPEN_PAREN)
	parser.advance()

	ifElse.Condition = parser.parseExpression(PRECEDENCE_LOWEST)

	parser.expect(lexing.TOKEN_CLOSE_PAREN)
	parser.advance()

	parser.expect(lexing.TOKEN_OPEN_BRACE)
	parser.advance()

	ifElse.Then = parser.parseCompound()

	parser.expect(lexing.TOKEN_CLOSE_BRACE)
	parser.advance()

	if parser.current.Type != lexing.TOKEN_ELSE {
		ifElse.Else = &AstCompound{
			Token:      ifElse.Token,
			Statements: []AstStatement{},
		}

		parser.commitError()
		return ifElse
	}

	parser.advance()

	parser.expect(lexing.TOKEN_OPEN_BRACE)
	parser.advance()

	ifElse.Else = parser.parseCompound()

	parser.expect(lexing.TOKEN_CLOSE_BRACE)
	parser.advance()

	parser.commitError()
	return ifElse
}

func (parser *Parser) parseExpression(precedence int) AstExpression {
	var left AstExpression

	switch parser.current.Type {
	case lexing.TOKEN_INTEGER:
		left = parser.parseIntegerLiteral()
	case lexing.TOKEN_TRUE, lexing.TOKEN_FALSE:
		left = parser.parseBooleanLiteral()
	case lexing.TOKEN_IDENTIFIER:
		left = parser.parseIdentifier()
	case lexing.TOKEN_STRING:
		left = parser.parseStringLiteral()
	case lexing.TOKEN_OPEN_PAREN:
		left = parser.parseEnforcedPrecedenceExpression()
	case lexing.TOKEN_OPEN_BRACE:
		left = parser.parseHashLiteral()
	case lexing.TOKEN_OPEN_BRACKET:
		left = parser.parseArrayLiteral()
	case lexing.TOKEN_FUNCTION:
		left = parser.parseFunctionDefinition()
	case lexing.TOKEN_IF:
		left = parser.parseIfElse()
	case lexing.TOKEN_BANG, lexing.TOKEN_MINUS:
		left = parser.parsePrefixExpression()
	}

	for precedence < getPrecedence(parser.current.Type) {
		switch parser.current.Type {
		case lexing.TOKEN_OPEN_PAREN:
			left = parser.parseFunctionCall(left)
		case lexing.TOKEN_OPEN_BRACKET:
			left = parser.parseIndex(left)
		default:
			left = parser.parseInfixExpression(left)
		}
	}

	return left
}

func (parser *Parser) parseExpressionStatement() *AstExpressionStatement {
	expressionStatement := &AstExpressionStatement{
		Token:      parser.current,
		Expression: parser.parseExpression(PRECEDENCE_LOWEST),
	}

	parser.expect(lexing.TOKEN_SEMICOLON)
	parser.advance()

	parser.commitError()
	return expressionStatement
}

func (parser *Parser) parseLetStatement() *AstLetStatement {
	letStatement := &AstLetStatement{
		Token: parser.current,
	}
	parser.advance()

	letStatement.Identifier = parser.parseIdentifier()

	parser.expect(lexing.TOKEN_ASSIGN)
	parser.advance()

	letStatement.Value = parser.parseExpression(PRECEDENCE_LOWEST)

	parser.expect(lexing.TOKEN_SEMICOLON)
	parser.advance()

	parser.commitError()
	return letStatement
}

func (parser *Parser) parseReturnStatement() *AstReturnStatement {
	returnStatement := &AstReturnStatement{
		Token: parser.current,
	}
	parser.advance()

	returnStatement.Value = parser.parseExpression(PRECEDENCE_LOWEST)

	parser.expect(lexing.TOKEN_SEMICOLON)
	parser.advance()

	parser.commitError()
	return returnStatement
}

func (parser *Parser) parseStatement() AstStatement {
	switch parser.current.Type {
	case lexing.TOKEN_LET:
		return parser.parseLetStatement()
	case lexing.TOKEN_RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) Parse() *AstCompound {
	compound := parser.parseCompound()

	parser.expect(lexing.TOKEN_EOF)

	parser.commitError()
	return compound
}
