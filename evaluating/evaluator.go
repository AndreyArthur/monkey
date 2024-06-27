package evaluating

import (
	"fmt"
	"monkey/parsing"
)

func objectError(format string, arguments ...interface{}) Object {
	return &ObjectError{
		Message: fmt.Sprintf(format, arguments...),
	}
}

func objectErrorInfixTypeMismatch(
	leftType ObjectType,
	operator string,
	rightType ObjectType,
) Object {
	return objectError(
		"Type mismatch: %s %s %s.",
		ObjectTypeToString(leftType),
		operator,
		ObjectTypeToString(rightType),
	)
}

func objectErrorUnknownInfixOperator(
	leftType ObjectType,
	operator string,
	rightType ObjectType,
) Object {
	return objectError(
		"Unknown operator: %s %s %s.",
		ObjectTypeToString(leftType),
		operator,
		ObjectTypeToString(rightType),
	)
}

func objectErrorPrefixTypeMismatch(
	operator string,
	right ObjectType,
) Object {
	return objectError(
		"Type mismatch: %s%s.",
		operator,
		ObjectTypeToString(right),
	)
}

func objectErrorUnknownPrefixOperator(
	operator string,
	right ObjectType,
) Object {
	return objectError(
		"Unknown operator: %s%s.",
		operator,
		ObjectTypeToString(right),
	)
}

func objectErrorIdentifierNotFound(name string) Object {
	return objectError("Identifier not found: %q.", name)
}

func objectErrorIdentifierAlreadyDeclared(name string) Object {
	return objectError("Identifier already declared in this scope: %q.", name)
}

func objectErrorNotCallable(expression parsing.AstExpression) Object {
	return objectError("Expression %q is not a callable.", expression.String())
}

func objectErrorUnsupportedIndex(indexType ObjectType) Object {
	return objectError(
		"Unsupported index, must be of type %s, %s, or %s, got type %s.",
		ObjectTypeToString(OBJECT_INTEGER),
		ObjectTypeToString(OBJECT_STRING),
		ObjectTypeToString(OBJECT_BOOLEAN),
		ObjectTypeToString(indexType),
	)
}

func objectErrornotIndexable(expression parsing.AstExpression) Object {
	return objectError("Expression %q is not a indexable.", expression.String())
}

func objectReturnValue(value Object) Object {
	return &ObjectReturnValue{Value: value}
}

func objectErrorWrongNumberOfArguments(
	expected int,
	got int,
) Object {
	return objectError(
		"Wrong number of arguments. Expected %d, got %d.",
		expected,
		got,
	)
}

func InjectBuiltinFunctions(environment *Environment) {
	environment.Set("len", &ObjectBuiltin{
		Function: func(arguments ...Object) Object {
			if len(arguments) != 1 {
				return objectErrorWrongNumberOfArguments(1, len(arguments))
			}

			object := arguments[0]

			switch arguments[0].Type() {
			case OBJECT_STRING:
				return &ObjectInteger{Value: int64(len(object.(*ObjectString).Value))}
			case OBJECT_ARRAY:
				return &ObjectInteger{Value: int64(len(object.(*ObjectArray).Items))}
			default:
				return objectError(
					"Type builtin function %q expects a string or array, got %s.",
					"len",
					ObjectTypeToString(object.Type()),
				)
			}
		},
	})
}

func evalCompound(
	environment *Environment,
	compound *parsing.AstCompound,
) Object {
	var last Object
	for _, statement := range compound.Statements {
		last = Eval(environment, statement)
		if last.Type() == OBJECT_RETURN_VALUE {
			return last
		}
	}
	return last
}

func evalExpressionStatement(
	environment *Environment,
	expressionStatement *parsing.AstExpressionStatement,
) Object {
	return Eval(environment, expressionStatement.Expression)
}

func evalEquality(left Object, operator string, right Object) Object {
	if left.Type() != right.Type() {
		return objectErrorInfixTypeMismatch(left.Type(), operator, right.Type())
	}

	if left.Type() == OBJECT_INTEGER && right.Type() == OBJECT_INTEGER {
		leftInteger := left.(*ObjectInteger).Value
		rightInteger := right.(*ObjectInteger).Value
		if operator == "==" {
			return &ObjectBoolean{Value: leftInteger == rightInteger}
		} else {
			return &ObjectBoolean{Value: leftInteger != rightInteger}
		}
	} else if left.Type() == OBJECT_BOOLEAN && right.Type() == OBJECT_BOOLEAN {
		leftBoolean := left.(*ObjectBoolean).Value
		rightBoolean := right.(*ObjectBoolean).Value
		if operator == "==" {
			return &ObjectBoolean{Value: leftBoolean == rightBoolean}
		} else {
			return &ObjectBoolean{Value: leftBoolean != rightBoolean}
		}
	} else {
		return objectErrorUnknownInfixOperator(left.Type(), operator, right.Type())
	}
}

func evalIntegerOperation(left Object, operator string, right Object) Object {
	if left.Type() != OBJECT_INTEGER || right.Type() != OBJECT_INTEGER {
		return objectErrorInfixTypeMismatch(left.Type(), operator, right.Type())
	}
	leftInteger := left.(*ObjectInteger).Value
	rightInteger := right.(*ObjectInteger).Value
	switch operator {
	case "+":
		return &ObjectInteger{Value: leftInteger + rightInteger}
	case "-":
		return &ObjectInteger{Value: leftInteger - rightInteger}
	case "*":
		return &ObjectInteger{Value: leftInteger * rightInteger}
	case "/":
		return &ObjectInteger{Value: leftInteger / rightInteger}
	case ">":
		return &ObjectBoolean{Value: leftInteger > rightInteger}
	case ">=":
		return &ObjectBoolean{Value: leftInteger >= rightInteger}
	case "<":
		return &ObjectBoolean{Value: leftInteger < rightInteger}
	case "<=":
		return &ObjectBoolean{Value: leftInteger <= rightInteger}
	default:
		return objectErrorUnknownInfixOperator(left.Type(), operator, right.Type())
	}
}

func evalStringConcatenation(left Object, right Object) Object {
	leftString := left.(*ObjectString).Value
	rightString := right.(*ObjectString).Value
	return &ObjectString{Value: leftString + rightString}
}

func evalInfixOperation(left Object, operator string, right Object) Object {
	switch operator {
	case "+":
		if left.Type() != right.Type() {
			return objectErrorInfixTypeMismatch(left.Type(), operator, right.Type())
		}
		if left.Type() == OBJECT_STRING {
			return evalStringConcatenation(left, right)
		}
		return evalIntegerOperation(left, operator, right)
	case "-", "*", "/", ">", ">=", "<", "<=":
		return evalIntegerOperation(left, operator, right)
	case "==", "!=":
		return evalEquality(left, operator, right)
	default:
		return objectErrorUnknownInfixOperator(left.Type(), operator, right.Type())
	}
}

func evalInfixExpression(
	environment *Environment,
	infixExpression *parsing.AstInfixExpression,
) Object {
	left := Eval(environment, infixExpression.Left)
	right := Eval(environment, infixExpression.Right)
	return evalInfixOperation(left, infixExpression.Operator, right)
}

func evalPrefixOperation(operator string, right Object) Object {
	switch operator {
	case "!":
		return &ObjectBoolean{
			Value: !right.Truthiness(),
		}
	case "-":
		if right.Type() != OBJECT_INTEGER {
			return objectErrorPrefixTypeMismatch(operator, right.Type())
		}
		return &ObjectInteger{
			Value: -right.(*ObjectInteger).Value,
		}
	default:
		return objectErrorUnknownPrefixOperator(operator, right.Type())
	}
}

func evalPrefixExpression(
	environment *Environment,
	prefixExpression *parsing.AstPrefixExpression,
) Object {
	right := Eval(environment, prefixExpression.Right)
	return evalPrefixOperation(prefixExpression.Operator, right)
}

func evalLetStatement(
	environment *Environment,
	letStatement *parsing.AstLetStatement,
) Object {
	if environment.Get(letStatement.Identifier.Name) != nil {
		return objectErrorIdentifierAlreadyDeclared(letStatement.Identifier.Name)
	}

	var value Object
	if letStatement.Value == nil {
		value = &ObjectNull{}
	} else {
		value = Eval(environment, letStatement.Value)
	}
	environment.Set(
		letStatement.Identifier.Name,
		value,
	)
	return &ObjectNull{}
}

func evalReturnStatement(
	environment *Environment,
	returnStatement *parsing.AstReturnStatement,
) Object {
	if returnStatement.Value == nil {
		return objectReturnValue(&ObjectNull{})
	}
	return objectReturnValue(Eval(environment, returnStatement.Value))
}

func evalIdentifier(
	environment *Environment,
	identifier *parsing.AstIdentifier,
) Object {
	object := environment.Get(identifier.Name)
	if object == nil {
		return objectErrorIdentifierNotFound(identifier.Name)
	}
	return object
}

func evalFunctionDefinition(
	environment *Environment,
	functionDefinition *parsing.AstFunctionDefinition,
) Object {
	for _, parameter := range functionDefinition.Parameters {
		if environment.Get(parameter.Name) != nil {
			return objectErrorIdentifierAlreadyDeclared(parameter.Name)
		}
	}

	return &ObjectFunction{
		Parameters:  functionDefinition.Parameters,
		Body:        functionDefinition.Body,
		Environment: environment,
	}
}

func evalExpressions(
	environment *Environment,
	expressions []parsing.AstExpression,
) []Object {
	objects := make([]Object, len(expressions))
	for i, expression := range expressions {
		objects[i] = Eval(environment, expression)
	}
	return objects
}

func extendFunctionEnvironment(
	function *ObjectFunction,
	arguments []Object,
) *Environment {
	environment := NewEnvironment(function.Environment)
	for i, parameter := range function.Parameters {
		environment.Set(parameter.Name, arguments[i])
	}
	return environment
}

func applyFunction(
	function *ObjectFunction,
	arguments []Object,
) Object {
	extendedEnvironment := extendFunctionEnvironment(
		function,
		arguments,
	)
	evaluated := Eval(extendedEnvironment, function.Body)
	if evaluated.Type() == OBJECT_RETURN_VALUE {
		return evaluated.(*ObjectReturnValue).Value
	}
	return evaluated
}

func evalFunctionCall(
	environment *Environment,
	functionCall *parsing.AstFunctionCall,
) Object {
	function := Eval(environment, functionCall.Left)
	if function.Type() != OBJECT_FUNCTION && function.Type() != OBJECT_BUILTIN {
		return objectErrorNotCallable(functionCall.Left)
	}
	arguments := evalExpressions(
		environment,
		functionCall.Arguments,
	)
	if function.Type() == OBJECT_BUILTIN {
		return function.(*ObjectBuiltin).Function(arguments...)
	}
	argumentsLen := len(arguments)
	parametersLen := len(function.(*ObjectFunction).Parameters)
	if argumentsLen != parametersLen {
		return objectErrorWrongNumberOfArguments(parametersLen, argumentsLen)
	}
	return applyFunction(function.(*ObjectFunction), arguments)
}

func evalArrayLiteral(
	environment *Environment,
	arrayLiteral *parsing.AstArrayLiteral,
) Object {
	items := evalExpressions(
		environment,
		arrayLiteral.Items,
	)
	return &ObjectArray{Items: items}
}

func evalHashLiteral(
	environment *Environment,
	hashLiteral *parsing.AstHashLiteral,
) Object {
	object := &ObjectHash{
		Keys:   []Object{},
		Values: []Object{},
	}
	for _, pair := range hashLiteral.Pairs {
		key := Eval(environment, pair.Key)
		if key.Type() != OBJECT_STRING &&
			key.Type() != OBJECT_INTEGER &&
			key.Type() != OBJECT_BOOLEAN {
			return objectErrorUnsupportedIndex(key.Type())
		}
		object.Keys = append(object.Keys, Eval(environment, pair.Key))
		object.Values = append(object.Values, Eval(environment, pair.Value))
	}
	return object
}

func evalArrayIndex(array *ObjectArray, indexObject Object) Object {
	index := indexObject.(*ObjectInteger).Value
	if index < 0 || index >= int64(len(array.Items)) {
		return &ObjectNull{}
	}
	return array.Items[index]
}

func evalHashIndex(hash *ObjectHash, key Object) Object {
	return hash.Get(key)
}

func evalIndex(
	environment *Environment,
	index *parsing.AstIndex,
) Object {
	left := Eval(environment, index.Left)
	key := Eval(environment, index.Index)
	if key.Type() != OBJECT_STRING &&
		key.Type() != OBJECT_INTEGER &&
		key.Type() != OBJECT_BOOLEAN {
		return objectErrorUnsupportedIndex(key.Type())
	}

	switch left.Type() {
	case OBJECT_ARRAY:
		return evalArrayIndex(left.(*ObjectArray), key)
	case OBJECT_HASH:
		return evalHashIndex(left.(*ObjectHash), key)
	default:
		return objectErrornotIndexable(index.Left)
	}
}

func evalIfElse(
	environment *Environment,
	ifElse *parsing.AstIfElse,
) Object {
	condition := Eval(environment, ifElse.Condition)
	if condition.Truthiness() == true {
		return Eval(environment, ifElse.Then)
	}
	if ifElse.Else == nil {
		return &ObjectNull{}
	}
	return Eval(environment, ifElse.Else)
}

func Eval(environment *Environment, ast parsing.AstNode) Object {
	switch ast.Type() {
	case parsing.AST_COMPOUND:
		return evalCompound(environment, ast.(*parsing.AstCompound))
	case parsing.AST_EXPRESSION_STATEMENT:
		return evalExpressionStatement(
			environment,
			ast.(*parsing.AstExpressionStatement),
		)
	case parsing.AST_LET_STATEMENT:
		return evalLetStatement(
			environment,
			ast.(*parsing.AstLetStatement),
		)
	case parsing.AST_RETURN_STATEMENT:
		return evalReturnStatement(
			environment,
			ast.(*parsing.AstReturnStatement),
		)
	case parsing.AST_INTEGER_LITERAL:
		return &ObjectInteger{
			Value: ast.(*parsing.AstIntegerLiteral).Value,
		}
	case parsing.AST_BOOLEAN_LITERAL:
		return &ObjectBoolean{
			Value: ast.(*parsing.AstBooleanLiteral).Value,
		}
	case parsing.AST_PREFIX_EXPRESSION:
		return evalPrefixExpression(
			environment,
			ast.(*parsing.AstPrefixExpression),
		)
	case parsing.AST_INFIX_EXPRESSION:
		return evalInfixExpression(
			environment,
			ast.(*parsing.AstInfixExpression),
		)
	case parsing.AST_IDENTIFIER:
		return evalIdentifier(
			environment,
			ast.(*parsing.AstIdentifier),
		)
	case parsing.AST_FUNCTION_DEFINITION:
		return evalFunctionDefinition(
			environment,
			ast.(*parsing.AstFunctionDefinition),
		)
	case parsing.AST_FUNCTION_CALL:
		return evalFunctionCall(
			environment,
			ast.(*parsing.AstFunctionCall),
		)
	case parsing.AST_ARRAY_LITERAL:
		return evalArrayLiteral(
			environment,
			ast.(*parsing.AstArrayLiteral),
		)
	case parsing.AST_STRING_LITERAL:
		return &ObjectString{
			Value: ast.(*parsing.AstStringLiteral).Value,
		}
	case parsing.AST_HASH_LITERAL:
		return evalHashLiteral(
			environment,
			ast.(*parsing.AstHashLiteral),
		)
	case parsing.AST_INDEX:
		return evalIndex(
			environment,
			ast.(*parsing.AstIndex),
		)
	case parsing.AST_IF_ELSE:
		return evalIfElse(
			environment,
			ast.(*parsing.AstIfElse),
		)
	default:
		// the switch will be exaustive so this should never happen
		return nil
	}
}
