package evaluator

import "compiler/object"

type Environment struct {
	currentEnvironment  map[string]object.Object
	extendedEnvironment *Environment
}

func NewEnvironment() *Environment {
	return &Environment{currentEnvironment: make(map[string]object.Object)}
}

func FromMap(source map[string]object.Object) *Environment {
	return &Environment{currentEnvironment: source}
}

func FromEnvironment(env *Environment) *Environment {
	return &Environment{extendedEnvironment: env, currentEnvironment: make(map[string]object.Object)}
}

func (env *Environment) get(key string) object.Object {
	value, ok := env.currentEnvironment[key]
	if !ok && env.extendedEnvironment != nil {
		value = env.extendedEnvironment.get(key)
	}
	return value
}

func (env *Environment) put(key string, value object.Object) {
	env.currentEnvironment[key] = value
}
