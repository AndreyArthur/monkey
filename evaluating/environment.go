package evaluating

type Environment struct {
	Store  map[string]Object
	Parent *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		Store:  map[string]Object{},
		Parent: parent,
	}
}
