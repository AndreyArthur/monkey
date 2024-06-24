package parsing_test

import (
	"monkey/lexing"
	"monkey/parsing"
	"testing"
)

func TestParseExpressionStatement(t *testing.T) {
	expectations := []struct {
		input  string
		output string
	}{
		{"-1", "(-1);"},
		{"2 + 2 * 4", "(2 + (2 * 4));"},
		{"2 - -2", "(2 - (-2));"},
		{"2 + 2 < 2 + 2", "((2 + 2) < (2 + 2));"},
		{"2 == 2 <= 2 + 2 * 2", "(2 == (2 <= (2 + (2 * 2))));"},
		{"2 < 3 == !false", "((2 < 3) == (!false));"},
		{"(2 + 2) * 6", "((2 + 2) * 6);"},
		{"2 - 2 + 2", "((2 - 2) + 2);"},
		{"2 - -my_variable + 2", "((2 - (-my_variable)) + 2);"},
		{"2 + add(1, 2 + 3)", "(2 + add(1, (2 + 3)));"},
	}

	for _, expectation := range expectations {
		lexer := lexing.NewLexer(expectation.input)
		parser := parsing.NewParser(lexer)
		ast := parser.Parse()
		output := ast.String()

		if output != expectation.output {
			t.Fatalf(
				"Expected: %q\nGot: %q",
				expectation.output,
				output,
			)
		}
	}
}
