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
	Type    TokenType
	Literal string
}

func TokenTypeToString(tokenType TokenType) string {
	tokenTypeToString := map[TokenType]string{
		TOKEN_EOF:               "eof",
		TOKEN_ILLEGAL:           "illegal",
		TOKEN_LET:               "let",
		TOKEN_FUNCTION:          "function",
		TOKEN_RETURN:            "return",
		TOKEN_IF:                "if",
		TOKEN_ELSE:              "else",
		TOKEN_IDENTIFIER:        "identifier",
		TOKEN_INTEGER:           "integer",
		TOKEN_STRING:            "string",
		TOKEN_ASSIGN:            "assign",
		TOKEN_PLUS:              "plus",
		TOKEN_MINUS:             "minus",
		TOKEN_ASTERISK:          "asterisk",
		TOKEN_SLASH:             "slash",
		TOKEN_BANG:              "bang",
		TOKEN_EQUALS:            "equals",
		TOKEN_NOT_EQUALS:        "not equals",
		TOKEN_GREATER:           "greater",
		TOKEN_GREATER_OR_EQUALS: "greater or equals",
		TOKEN_LESS:              "less",
		TOKEN_LESS_OR_EQUALS:    "less or equals",
		TOKEN_OPEN_PAREN:        "open paren",
		TOKEN_CLOSE_PAREN:       "close paren",
		TOKEN_OPEN_BRACE:        "open brace",
		TOKEN_CLOSE_BRACE:       "close brace",
		TOKEN_OPEN_BRACKET:      "open bracket",
		TOKEN_CLOSE_BRACKET:     "close bracket",
		TOKEN_COMMA:             "comma",
		TOKEN_COLON:             "colon",
		TOKEN_SEMICOLON:         "semicolon",
	}
	return tokenTypeToString[tokenType]
}

func NewToken(tokenType TokenType, literal string) *Token {
	return &Token{
		Type:    tokenType,
		Literal: literal,
	}
}
