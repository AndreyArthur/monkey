package lexing

import "bytes"

func isDigit(character byte) bool {
	return character >= '0' && character <= '9'
}

func isAlphabetic(character byte) bool {
	return (character >= 'A' && character <= 'Z') ||
		(character >= 'a' && character <= 'z') ||
		character == '_'
}

func isAlphanumeric(character byte) bool {
	return isDigit(character) || isAlphabetic(character)
}

type Lexer struct {
	content  string
	position int
	current  byte
}

func NewLexer(content string) *Lexer {
	return &Lexer{
		content:  content,
		position: 0,
		current:  content[0],
	}
}

func (lexer *Lexer) advance() {
	if lexer.position+1 >= len(lexer.content) {
		lexer.current = '\x00'
		return
	}
	lexer.position += 1
	lexer.current = lexer.content[lexer.position]
}

func (lexer *Lexer) peek() byte {
	if lexer.position+1 >= len(lexer.content) {
		return '\x00'
	}
	return lexer.content[lexer.position+1]
}

func (lexer *Lexer) skipWhitespaces() {
	for lexer.current == ' ' ||
		lexer.current == '\n' ||
		lexer.current == '\r' ||
		lexer.current == '\t' {
		lexer.advance()
	}
}

func (lexer *Lexer) collectCurrent(tokenType TokenType) *Token {
	token := NewToken(tokenType, string(lexer.current))
	lexer.advance()
	return token
}

func (lexer *Lexer) collectWithNext(tokenType TokenType) *Token {
	token := NewToken(tokenType, string(lexer.current)+string(lexer.peek()))
	lexer.advance()
	lexer.advance()
	return token
}

func (lexer *Lexer) collectIntegerLiteral() *Token {
	var buffer bytes.Buffer

	for isDigit(lexer.current) {
		buffer.WriteByte(lexer.current)
		lexer.advance()
	}

	return NewToken(TOKEN_INTEGER, buffer.String())
}

func (lexer *Lexer) collectStringLiteral() *Token {
	var buffer bytes.Buffer

	lexer.advance()
	// TODO: handle escaope characters
	for lexer.current != '"' {
		buffer.WriteByte(lexer.current)
		lexer.advance()
	}
	lexer.advance()

	return NewToken(TOKEN_STRING, buffer.String())
}

func (lexer *Lexer) collectIdentifierOrKeyword() *Token {
	var buffer bytes.Buffer

	for isAlphanumeric(lexer.current) {
		buffer.WriteByte(lexer.current)
		lexer.advance()
	}

	text := buffer.String()

	var tokenType TokenType

	switch text {
	case "let":
		tokenType = TOKEN_LET
	case "fn":
		tokenType = TOKEN_FUNCTION
	case "return":
		tokenType = TOKEN_RETURN
	case "if":
		tokenType = TOKEN_IF
	case "else":
		tokenType = TOKEN_ELSE
	case "true":
		tokenType = TOKEN_TRUE
	case "false":
		tokenType = TOKEN_FALSE
	default:
		tokenType = TOKEN_IDENTIFIER
	}

	return NewToken(tokenType, text)
}

func (lexer *Lexer) Next() *Token {
	lexer.skipWhitespaces()

	switch lexer.current {
	case '\x00':
		return lexer.collectCurrent(TOKEN_EOF)
	case '=':
		if lexer.peek() == '=' {
			return lexer.collectWithNext(TOKEN_EQUALS)
		}
		return lexer.collectCurrent(TOKEN_ASSIGN)
	case '+':
		return lexer.collectCurrent(TOKEN_PLUS)
	case '-':
		return lexer.collectCurrent(TOKEN_MINUS)
	case '*':
		return lexer.collectCurrent(TOKEN_ASTERISK)
	case '/':
		return lexer.collectCurrent(TOKEN_SLASH)
	case '!':
		if lexer.peek() == '=' {
			return lexer.collectWithNext(TOKEN_NOT_EQUALS)
		}
		return lexer.collectCurrent(TOKEN_BANG)
	case '>':
		if lexer.peek() == '=' {
			return lexer.collectWithNext(TOKEN_GREATER_OR_EQUALS)
		}
		return lexer.collectCurrent(TOKEN_GREATER)
	case '<':
		if lexer.peek() == '=' {
			return lexer.collectWithNext(TOKEN_LESS_OR_EQUALS)
		}
		return lexer.collectCurrent(TOKEN_LESS)
	case '(':
		return lexer.collectCurrent(TOKEN_OPEN_PAREN)
	case ')':
		return lexer.collectCurrent(TOKEN_CLOSE_PAREN)
	case '{':
		return lexer.collectCurrent(TOKEN_OPEN_BRACE)
	case '}':
		return lexer.collectCurrent(TOKEN_CLOSE_BRACE)
	case '[':
		return lexer.collectCurrent(TOKEN_OPEN_BRACKET)
	case ']':
		return lexer.collectCurrent(TOKEN_CLOSE_BRACKET)
	case ',':
		return lexer.collectCurrent(TOKEN_COMMA)
	case ':':
		return lexer.collectCurrent(TOKEN_COLON)
	case ';':
		return lexer.collectCurrent(TOKEN_SEMICOLON)
	case '"':
		return lexer.collectStringLiteral()
	default:
		if isDigit(lexer.current) {
			return lexer.collectIntegerLiteral()
		}
		if isAlphabetic(lexer.current) {
			return lexer.collectIdentifierOrKeyword()
		}
		return lexer.collectCurrent(TOKEN_ILLEGAL)
	}
}
