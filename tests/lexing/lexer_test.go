package lexing_test

import (
	"monkey/lexing"
	"testing"
)

func TestLexer(t *testing.T) {
	content := `let a = 2 + 2 - 4;

let b = 2 * 2 / 4;

let function = fn () {
    return !true != false;
};

let c = if (function) true else false;

2 <= 2;
2 < 3;
2 >= 2;
2 > 1;
5 == 6;

"Hello, World";

let d = {"hello": "world"};

d["hello"];

let e = [1, 2, 3, 4, 5];

e[4] == 5;

@;
`

	expectations := []struct {
		tokenType lexing.TokenType
		value     string
	}{
		{lexing.TOKEN_LET, "let"},
		{lexing.TOKEN_IDENTIFIER, "a"},
		{lexing.TOKEN_ASSIGN, "="},
		{lexing.TOKEN_INTEGER, "2"},
		{lexing.TOKEN_PLUS, "+"},
		{lexing.TOKEN_INTEGER, "2"},
		{lexing.TOKEN_MINUS, "-"},
		{lexing.TOKEN_INTEGER, "4"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_LET, "let"},
		{lexing.TOKEN_IDENTIFIER, "b"},
		{lexing.TOKEN_ASSIGN, "="},
		{lexing.TOKEN_INTEGER, "2"},
		{lexing.TOKEN_ASTERISK, "*"},
		{lexing.TOKEN_INTEGER, "2"},
		{lexing.TOKEN_SLASH, "/"},
		{lexing.TOKEN_INTEGER, "4"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_LET, "let"},
		{lexing.TOKEN_IDENTIFIER, "function"},
		{lexing.TOKEN_ASSIGN, "="},
		{lexing.TOKEN_FUNCTION, "fn"},
		{lexing.TOKEN_OPEN_PAREN, "("},
		{lexing.TOKEN_CLOSE_PAREN, ")"},
		{lexing.TOKEN_OPEN_BRACE, "{"},
		{lexing.TOKEN_RETURN, "return"},
		{lexing.TOKEN_BANG, "!"},
		{lexing.TOKEN_IDENTIFIER, "true"},
		{lexing.TOKEN_NOT_EQUALS, "!="},
		{lexing.TOKEN_IDENTIFIER, "false"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_CLOSE_BRACE, "}"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_LET, "let"},
		{lexing.TOKEN_IDENTIFIER, "c"},
		{lexing.TOKEN_ASSIGN, "="},
		{lexing.TOKEN_IF, "if"},
		{lexing.TOKEN_OPEN_PAREN, "("},
		{lexing.TOKEN_IDENTIFIER, "function"},
		{lexing.TOKEN_CLOSE_PAREN, ")"},
		{lexing.TOKEN_IDENTIFIER, "true"},
		{lexing.TOKEN_ELSE, "else"},
		{lexing.TOKEN_IDENTIFIER, "false"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_INTEGER, "2"},
		{lexing.TOKEN_LESS_OR_EQUALS, "<="},
		{lexing.TOKEN_INTEGER, "2"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_INTEGER, "2"},
		{lexing.TOKEN_LESS, "<"},
		{lexing.TOKEN_INTEGER, "3"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_INTEGER, "2"},
		{lexing.TOKEN_GREATER_OR_EQUALS, ">="},
		{lexing.TOKEN_INTEGER, "2"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_INTEGER, "2"},
		{lexing.TOKEN_GREATER, ">"},
		{lexing.TOKEN_INTEGER, "1"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_INTEGER, "5"},
		{lexing.TOKEN_EQUALS, "=="},
		{lexing.TOKEN_INTEGER, "6"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_STRING, "Hello, World"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_LET, "let"},
		{lexing.TOKEN_IDENTIFIER, "d"},
		{lexing.TOKEN_ASSIGN, "="},
		{lexing.TOKEN_OPEN_BRACE, "{"},
		{lexing.TOKEN_STRING, "hello"},
		{lexing.TOKEN_COLON, ":"},
		{lexing.TOKEN_STRING, "world"},
		{lexing.TOKEN_CLOSE_BRACE, "}"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_IDENTIFIER, "d"},
		{lexing.TOKEN_OPEN_BRACKET, "["},
		{lexing.TOKEN_STRING, "hello"},
		{lexing.TOKEN_CLOSE_BRACKET, "]"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_LET, "let"},
		{lexing.TOKEN_IDENTIFIER, "e"},
		{lexing.TOKEN_ASSIGN, "="},
		{lexing.TOKEN_OPEN_BRACKET, "["},
		{lexing.TOKEN_INTEGER, "1"},
		{lexing.TOKEN_COMMA, ","},
		{lexing.TOKEN_INTEGER, "2"},
		{lexing.TOKEN_COMMA, ","},
		{lexing.TOKEN_INTEGER, "3"},
		{lexing.TOKEN_COMMA, ","},
		{lexing.TOKEN_INTEGER, "4"},
		{lexing.TOKEN_COMMA, ","},
		{lexing.TOKEN_INTEGER, "5"},
		{lexing.TOKEN_CLOSE_BRACKET, "]"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_IDENTIFIER, "e"},
		{lexing.TOKEN_OPEN_BRACKET, "["},
		{lexing.TOKEN_INTEGER, "4"},
		{lexing.TOKEN_CLOSE_BRACKET, "]"},
		{lexing.TOKEN_EQUALS, "=="},
		{lexing.TOKEN_INTEGER, "5"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_ILLEGAL, "@"},
		{lexing.TOKEN_SEMICOLON, ";"},
		{lexing.TOKEN_EOF, "\x00"},
	}

	lexer := lexing.NewLexer(content)

	for index, expectation := range expectations {
		token := lexer.Next()

		if token.Type != expectation.tokenType ||
			token.Literal != expectation.value {
			t.Fatalf(
				"[%d] Expected token %q of type %d, found token %q of type %d.",
				index,
				expectation.value,
				expectation.tokenType,
				token.Literal,
				token.Type,
			)
		}
	}
}
