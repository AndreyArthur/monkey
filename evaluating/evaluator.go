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

func evalCompound(
	environment *Environment,
	compound *parsing.AstCompound,
) Object {
	var last Object
	for _, statement := range compound.Statements {
		last = Eval(environment, statement)
		if statement.Type() == parsing.AST_RETURN_STATEMENT {
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

func evalInfixOperation(left Object, operator string, right Object) Object {
	switch operator {
	case "+", "-", "*", "/", ">", ">=", "<", "<=":
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
		return &ObjectNull{}
	}
	return Eval(environment, returnStatement.Value)
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
	default:
		// the switch will be exaustive so this should never happen
		return nil
	}
}
