package evaluating

import "fmt"

const (
	_ = iota
	OBJECT_ERROR
	OBJECT_INTEGER
	OBJECT_BOOLEAN
)

type ObjectType int

func ObjectTypeToString(objectType ObjectType) string {
	switch objectType {
	case OBJECT_INTEGER:
		return "integer"
	case OBJECT_BOOLEAN:
		return "boolean"
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
