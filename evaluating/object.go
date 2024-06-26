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
)

type ObjectType int

func ObjectTypeToString(objectType ObjectType) string {
	switch objectType {
	case OBJECT_INTEGER:
		return "integer"
	case OBJECT_BOOLEAN:
		return "boolean"
	case OBJECT_ERROR:
		return "error"
	case OBJECT_FUNCTION:
		return "function"
	case OBJECT_ARRAY:
		return "array"
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
