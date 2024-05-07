package evaluator

type Environment struct {
	currentEnvironment  map[string]Object
	extendedEnvironment *Environment
}

func NewEnvironment() *Environment {
	return &Environment{currentEnvironment: make(map[string]Object)}
}

func FromEnvironment(env *Environment) *Environment {
	return &Environment{extendedEnvironment: env, currentEnvironment: make(map[string]Object)}
}

func (env *Environment) get(key string) Object {
	value, ok := env.currentEnvironment[key]
	if !ok && env.extendedEnvironment != nil {
		value = env.extendedEnvironment.get(key)
	}
	return value
}

func (env *Environment) put(key string, value Object) {
	env.currentEnvironment[key] = value
}
