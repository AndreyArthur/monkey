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
		{"-1;", "(-1);"},
		{"2 + 2 * 4;", "(2 + (2 * 4));"},
		{"2 - -2;", "(2 - (-2));"},
		{"2 + 2 < 2 + 2;", "((2 + 2) < (2 + 2));"},
		{"2 == 2 <= 2 + 2 * 2;", "(2 == (2 <= (2 + (2 * 2))));"},
		{"2 < 3 == !false;", "((2 < 3) == (!false));"},
		{"(2 + 2) * 6;", "((2 + 2) * 6);"},
		{"2 - 2 + 2;", "((2 - 2) + 2);"},
		{"2 - -my_variable + 2;", "((2 - (-my_variable)) + 2);"},
		{"2 + add(1, 2 + 3);", "(2 + add(1, (2 + 3)));"},
		{"array[1];", "array[1];"},
		{"hashmap[3 + -2];", "hashmap[(3 + (-2))];"},
		{"hashmap[!false];", "hashmap[(!false)];"},
		{"-array[1];", "(-array[1]);"},
		{"\"Hello, \" + \"World!\";", "(\"Hello, \" + \"World!\");"},
		{"[\"name\", 2, true, 5 < 4];", "[\"name\", 2, true, (5 < 4)];"},
		{"[1, 2, 3][0];", "[1, 2, 3][0];"},
		{"{\"name\": \"christian\", true: 2 + 2};", "{\"name\": \"christian\", true: (2 + 2)};"},
		{"fn (a, b) { a + b; a * b; };", "fn (a, b) { (a + b); (a * b); };"},
		{"fn () {};", "fn () {  };"},
	}

	for _, expectation := range expectations {
		lexer := lexing.NewLexer(expectation.input)
		parser := parsing.NewParser(lexer)
		ast := parser.Parse()

		if parser.HasErrors() {
			for _, error := range parser.GetErrors() {
				t.Log(error)
			}
			t.Fail()
		}

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

func TestParseLetStatement(t *testing.T) {
	expectations := []struct {
		input  string
		output string
	}{
		{"let a = -1;", "let a = (-1);"},
		{"let b = true;", "let b = true;"},
		{"let c = \"Hello, World\";", "let c = \"Hello, World\";"},
	}

	for _, expectation := range expectations {
		lexer := lexing.NewLexer(expectation.input)
		parser := parsing.NewParser(lexer)
		ast := parser.Parse()

		if parser.HasErrors() {
			for _, error := range parser.GetErrors() {
				t.Log(error)
			}
			t.Fail()
		}

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
func TestParserErrors(t *testing.T) {
	expectations := []struct {
		input string
		error string
	}{
		{"-1", `Expected token of type semicolon. Found token "\x00" of type eof.`},
		{"myfunction(4; 5);", `Expected token of type comma. Found token ";" of type semicolon.`},
		{"myfunction(4, 5};", `Expected token of type comma. Found token "}" of type close brace.`},
		{"array[4};", `Expected token of type close bracket. Found token "}" of type close brace.`},
		{"[4; 5];", `Expected token of type comma. Found token ";" of type semicolon.`},
		{"[4, 5};", `Expected token of type comma. Found token "}" of type close brace.`},
		{"{4; 5};", `Expected token of type colon. Found token ";" of type semicolon.`},
		{"{4: 5];", `Expected token of type comma. Found token "]" of type close bracket.`},
		{"fn (2) {};", `Expected token of type identifier. Found token "2" of type integer.`},
	}

	for _, expectation := range expectations {
		lexer := lexing.NewLexer(expectation.input)
		parser := parsing.NewParser(lexer)
		_ = parser.Parse()

		if parser.HasErrors() {
			errors := parser.GetErrors()
			if expectation.error != errors[0] {
				t.Fatalf("Expected: %q\nGot: %q", expectation.error, errors[0])
			}
			t.Fail()
		} else {
			t.Fatalf("Expected errors.")
		}

	}
}
