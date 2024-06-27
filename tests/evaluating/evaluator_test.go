package evaluating_test

import (
	"monkey/evaluating"
	"monkey/lexing"
	"monkey/parsing"
	"testing"
)

func TestEvalExpressions(t *testing.T) {
	expectations := []struct {
		input      string
		objectType evaluating.ObjectType
		output     any
	}{
		{"2 + 2;", evaluating.OBJECT_INTEGER, 4},
		{"2 * 2 == 2 + 2;", evaluating.OBJECT_BOOLEAN, true},
		{"!0 == true;", evaluating.OBJECT_BOOLEAN, true},
		{"!2 == true;", evaluating.OBJECT_BOOLEAN, false},
		{"5 / 2;", evaluating.OBJECT_INTEGER, 2},
		{"1 + 7 * 2;", evaluating.OBJECT_INTEGER, 15},
		{"let a = 2;", evaluating.OBJECT_NULL, nil},
		{"let a = true; a;", evaluating.OBJECT_BOOLEAN, true},
		{"fn (a, b) { return a + b; };", evaluating.OBJECT_FUNCTION, "fn (a, b)"},
		{"fn (a, b) { return a + b; }(1, 2);", evaluating.OBJECT_INTEGER, 3},
		{"fn (a) { return fn (b) { return a + b; }; }(2)(1);", evaluating.OBJECT_INTEGER, 3},
		{"[!2, 4 + 8, true, false];", evaluating.OBJECT_ARRAY, "[false, 12, true, false]"},
		{"\"Hello, \" + \"World!\";", evaluating.OBJECT_STRING, "\"Hello, World!\""},
		{"{4 - 2: false, !0: true, \"hello\": \"world\"}", evaluating.OBJECT_HASH, "{2: false, true: true, \"hello\": \"world\"}"},
		{"{4 - 2: false, !0: true, \"hello\": \"world\"}[\"hello\"]", evaluating.OBJECT_STRING, "\"world\""},
		{"if (1 > 0) { true; } else { false; };", evaluating.OBJECT_BOOLEAN, true},
		{"if (1 < 0) { true; } else { false; };", evaluating.OBJECT_BOOLEAN, false},
	}

	for _, expectation := range expectations {
		lexer := lexing.NewLexer(expectation.input)
		parser := parsing.NewParser(lexer)
		ast := parser.Parse()
		environment := evaluating.NewEnvironment(nil)
		object := evaluating.Eval(environment, ast)

		if object.Type() != expectation.objectType {
			t.Fatalf(
				"Expected object type to be %s, got %s.",
				evaluating.ObjectTypeToString(expectation.objectType),
				evaluating.ObjectTypeToString(object.Type()),
			)
		}

		switch object.(type) {
		case *evaluating.ObjectInteger:
			if object.(*evaluating.ObjectInteger).Value != int64(expectation.output.(int)) {
				t.Fatalf(
					"Expected %v, got %v.",
					expectation.output,
					object.(*evaluating.ObjectInteger).Value,
				)
			}
		case *evaluating.ObjectBoolean:
			if object.(*evaluating.ObjectBoolean).Value != expectation.output {
				t.Fatalf(
					"Expected %v, got %v.",
					expectation.output,
					object.(*evaluating.ObjectBoolean).Value,
				)
			}
		case *evaluating.ObjectNull:
			if expectation.output != nil {
				t.Fatalf(
					"Expected %v, got %v.",
					expectation.output,
					object.(*evaluating.ObjectNull),
				)
			}
		default:
			if object.Inspect() != string(expectation.output.(string)) {
				t.Fatalf(
					"Expected %v, got %v.",
					expectation.output,
					object.Inspect(),
				)
			}
		}
	}
}

func TestBuiltins(t *testing.T) {
	expectations := []struct {
		input      string
		objectType evaluating.ObjectType
		output     any
	}{
		{"len(\"Hello, World!\");", evaluating.OBJECT_INTEGER, 13},
		{"len(\"\");", evaluating.OBJECT_INTEGER, 0},
		{"len([1, true, fn () { return \"hello\"; }]);", evaluating.OBJECT_INTEGER, 3},
		{"len(2);", evaluating.OBJECT_ERROR, "Type builtin function \"len\" expects a string or array, got integer."},
	}

	for _, expectation := range expectations {
		lexer := lexing.NewLexer(expectation.input)
		parser := parsing.NewParser(lexer)
		ast := parser.Parse()
		environment := evaluating.NewEnvironment(nil)
		evaluating.InjectBuiltinFunctions(environment)
		object := evaluating.Eval(environment, ast)

		if object.Type() != expectation.objectType {
			t.Fatalf(
				"Expected object type to be %s, got %s.",
				evaluating.ObjectTypeToString(expectation.objectType),
				evaluating.ObjectTypeToString(object.Type()),
			)
		}

		switch object.(type) {
		case *evaluating.ObjectInteger:
			if object.(*evaluating.ObjectInteger).Value != int64(expectation.output.(int)) {
				t.Fatalf(
					"Expected %v, got %v.",
					expectation.output,
					object.(*evaluating.ObjectInteger).Value,
				)
			}
		case *evaluating.ObjectBoolean:
			if object.(*evaluating.ObjectBoolean).Value != expectation.output {
				t.Fatalf(
					"Expected %v, got %v.",
					expectation.output,
					object.(*evaluating.ObjectBoolean).Value,
				)
			}
		case *evaluating.ObjectNull:
			if expectation.output != nil {
				t.Fatalf(
					"Expected %v, got %v.",
					expectation.output,
					object.(*evaluating.ObjectNull),
				)
			}
		default:
			if object.Inspect() != string(expectation.output.(string)) {
				t.Fatalf(
					"Expected %v, got %v.",
					expectation.output,
					object.Inspect(),
				)
			}
		}
	}
}

func TestEvalError(t *testing.T) {
	expectations := []struct {
		input        string
		errorMessage string
	}{
		{"-false;", "Type mismatch: -boolean."},
		{"true + true;", "Type mismatch: boolean + boolean."},
		{"true + 2;", "Type mismatch: boolean + integer."},
		{"2 * false;", "Type mismatch: integer * boolean."},
		{"a;", "Identifier not found: \"a\"."},
		{"let a = 2; let a = 3;", "Identifier already declared in this scope: \"a\"."},
		{"let a = 2; fn (a) { a; };", "Identifier already declared in this scope: \"a\"."},
		{"fn (a) { return a; }(2, 3);", "Wrong number of arguments. Expected 1, got 2."},
		{"fn (a) { return a; }();", "Wrong number of arguments. Expected 1, got 0."},
	}

	for _, expectation := range expectations {
		lexer := lexing.NewLexer(expectation.input)
		parser := parsing.NewParser(lexer)
		ast := parser.Parse()
		environment := evaluating.NewEnvironment(nil)
		object := evaluating.Eval(environment, ast)

		if object.Type() != evaluating.OBJECT_ERROR {
			t.Fatalf(
				"Expected object type to be %s, got %s.",
				evaluating.ObjectTypeToString(evaluating.OBJECT_ERROR),
				evaluating.ObjectTypeToString(object.Type()),
			)
		}

		message := object.(*evaluating.ObjectError).Message

		if message != expectation.errorMessage {
			t.Fatalf(
				"Expected %q, got %q.",
				expectation.errorMessage,
				message,
			)
		}
	}
}
