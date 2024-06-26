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

func (environment *Environment) Get(name string) Object {
	currentEnvironment := environment
	for currentEnvironment != nil {
		if object, ok := currentEnvironment.Store[name]; ok {
			return object
		}
		currentEnvironment = currentEnvironment.Parent
	}
	return nil
}

func (environment *Environment) Set(name string, value Object) {
	environment.Store[name] = value
}
