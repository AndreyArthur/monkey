package evaluating

import (
	"fmt"
	"monkey/parsing"
)

const (
	_ = iota
	OBJECT_ERROR
	OBJECT_INTEGER
	OBJECT_BOOLEAN
	OBJECT_NULL
	OBJECT_ARRAY
	OBJECT_FUNCTION
	OBJECT_HASH
	OBJECT_STRING
	OBJECT_RETURN_VALUE
	OBJECT_BUILTIN
)

type ObjectType int

func ObjectTypeToString(objectType ObjectType) string {
	switch objectType {
	case OBJECT_ERROR:
		return "error"
	case OBJECT_INTEGER:
		return "integer"
	case OBJECT_BOOLEAN:
		return "boolean"
	case OBJECT_NULL:
		return "null"
	case OBJECT_ARRAY:
		return "array"
	case OBJECT_FUNCTION:
		return "function"
	case OBJECT_HASH:
		return "hash"
	case OBJECT_STRING:
		return "string"
	case OBJECT_RETURN_VALUE:
		return "return value"
	case OBJECT_BUILTIN:
		return "builtin"
	default:
		return "unknown"
	}
}

type Object interface {
	Type() ObjectType
	Inspect() string
	Truthiness() bool
}

type ObjectError struct {
	Message string
}

func (error *ObjectError) Type() ObjectType {
	return OBJECT_ERROR
}
func (error *ObjectError) Inspect() string {
	return error.Message
}
func (error *ObjectError) Truthiness() bool {
	return true
}

type ObjectNull struct{}

func (null *ObjectNull) Type() ObjectType {
	return OBJECT_NULL
}
func (null *ObjectNull) Inspect() string {
	return "null"
}
func (null *ObjectNull) Truthiness() bool {
	return false
}

type ObjectInteger struct {
	Value int64
}

func (integer *ObjectInteger) Type() ObjectType {
	return OBJECT_INTEGER
}
func (integer *ObjectInteger) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}
func (integer *ObjectInteger) Truthiness() bool {
	if integer.Value == 0 {
		return false
	}
	return true
}

type ObjectBoolean struct {
	Value bool
}

func (boolean *ObjectBoolean) Type() ObjectType {
	return OBJECT_BOOLEAN
}
func (boolean *ObjectBoolean) Inspect() string {
	return fmt.Sprintf("%t", boolean.Value)
}
func (boolean *ObjectBoolean) Truthiness() bool {
	return boolean.Value
}

type ObjectFunction struct {
	Parameters  []*parsing.AstIdentifier
	Body        *parsing.AstCompound
	Environment *Environment
}

func (function *ObjectFunction) Type() ObjectType {
	return OBJECT_FUNCTION
}
func (function *ObjectFunction) Inspect() string {
	text := "fn ("

	for index, parameter := range function.Parameters {
		text += parameter.String()
		if index < len(function.Parameters)-1 {
			text += ", "
		}
	}

	text += ")"

	return text
}
func (function *ObjectFunction) Truthiness() bool {
	return true
}

type ObjectArray struct {
	Items []Object
}

func (array *ObjectArray) Type() ObjectType {
	return OBJECT_ARRAY
}
func (array *ObjectArray) Inspect() string {
	text := "["
	for index, element := range array.Items {
		text += element.Inspect()
		if index < len(array.Items)-1 {
			text += ", "
		}
	}
	text += "]"
	return text
}
func (array *ObjectArray) Truthiness() bool {
	return true
}

type ObjectString struct {
	Value string
}

func (string *ObjectString) Type() ObjectType {
	return OBJECT_STRING
}
func (string *ObjectString) Inspect() string {
	return fmt.Sprintf("%q", string.Value)
}
func (string *ObjectString) Truthiness() bool {
	if string.Value == "" {
		return false
	}
	return true
}

type ObjectHash struct {
	Keys   []Object
	Values []Object
}

func (hash *ObjectHash) Type() ObjectType {
	return OBJECT_HASH
}
func (hash *ObjectHash) Inspect() string {
	text := "{"
	for index, key := range hash.Keys {
		text += key.Inspect() + ": " + hash.Values[index].Inspect()
		if index < len(hash.Keys)-1 {
			text += ", "
		}
	}
	text += "}"
	return text
}
func (hash *ObjectHash) Truthiness() bool {
	return true
}
func (hash *ObjectHash) Get(key Object) Object {
	for index, hashKey := range hash.Keys {
		if hashKey.Type() != key.Type() {
			continue
		}
		switch key.Type() {
		case OBJECT_STRING:
			if hashKey.(*ObjectString).Value == key.(*ObjectString).Value {
				return hash.Values[index]
			}
		case OBJECT_INTEGER:
			if hashKey.(*ObjectInteger).Value == key.(*ObjectInteger).Value {
				return hash.Values[index]
			}
		case OBJECT_BOOLEAN:
			if hashKey.(*ObjectBoolean).Value == key.(*ObjectBoolean).Value {
				return hash.Values[index]
			}
		default:
			continue
		}
	}
	return &ObjectNull{}
}

type ObjectReturnValue struct {
	Value Object
}

func (returnValue *ObjectReturnValue) Type() ObjectType {
	return OBJECT_RETURN_VALUE
}
func (returnValue *ObjectReturnValue) Inspect() string {
	return "return " + returnValue.Value.Inspect()
}
func (returnValue *ObjectReturnValue) Truthiness() bool {
	return returnValue.Value.Truthiness()
}

type BuiltinFunction func(arguments ...Object) Object

type ObjectBuiltin struct {
	Function BuiltinFunction
}

func (builtin *ObjectBuiltin) Type() ObjectType {
	return OBJECT_BUILTIN
}
func (builtin *ObjectBuiltin) Inspect() string {
	return "fn (...)"
}
func (builtin *ObjectBuiltin) Truthiness() bool {
	return true
}
