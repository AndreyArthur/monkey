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
	if environment.Get(name) == nil {
		environment.Store[name] = value
	}

	currentEnvironment := environment
	_, ok := currentEnvironment.Store[name]
	for !ok {
		currentEnvironment = currentEnvironment.Parent
		_, ok = currentEnvironment.Store[name]
	}

	currentEnvironment.Store[name] = value
}
