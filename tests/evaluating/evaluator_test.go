package evaluating_test

import (
	"monkey/evaluating"
	"monkey/lexing"
	"monkey/parsing"
	"testing"
)

func TestEval(t *testing.T) {
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
	}

	for _, expectation := range expectations {
		lexer := lexing.NewLexer(expectation.input)
		parser := parsing.NewParser(lexer)
		ast := parser.Parse()
		environment := evaluating.NewEnvironment(nil)
		object := evaluating.Eval(environment, ast)

		if object.Type() != expectation.objectType {
			t.Fatalf(
				"Expected %d, got %d.",
				expectation.objectType,
				object.Type(),
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
	}

	for _, expectation := range expectations {
		lexer := lexing.NewLexer(expectation.input)
		parser := parsing.NewParser(lexer)
		ast := parser.Parse()
		environment := evaluating.NewEnvironment(nil)
		object := evaluating.Eval(environment, ast)

		if object.Type() != evaluating.OBJECT_ERROR {
			t.Fatalf(
				"Expected %d, got %d.",
				evaluating.OBJECT_ERROR,
				object.Type(),
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
