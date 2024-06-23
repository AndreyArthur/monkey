package lexing

const (
	_ = iota
	TOKEN_EOF
	TOKEN_ILLEGAL

	TOKEN_LET
	TOKEN_FUNCTION
	TOKEN_RETURN
	TOKEN_IF
	TOKEN_ELSE

	TOKEN_IDENTIFIER

	TOKEN_INTEGER
	TOKEN_STRING

	TOKEN_ASSIGN
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_ASTERISK
	TOKEN_SLASH
	TOKEN_BANG
	TOKEN_EQUALS
	TOKEN_NOT_EQUALS
	TOKEN_GREATER
	TOKEN_GREATER_OR_EQUALS
	TOKEN_LESS
	TOKEN_LESS_OR_EQUALS

	TOKEN_OPEN_PAREN
	TOKEN_CLOSE_PAREN
	TOKEN_OPEN_BRACE
	TOKEN_CLOSE_BRACE
	TOKEN_OPEN_BRACKET
	TOKEN_CLOSE_BRACKET
	TOKEN_COMMA
	TOKEN_COLON
	TOKEN_SEMICOLON
)

type TokenType int

type Token struct {
	Type  TokenType
	Value string
}

func NewToken(tokenType TokenType, value string) *Token {
	return &Token{
		Type:  tokenType,
		Value: value,
	}
}
