package lox

// Environment represents the runtime environment.
type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

// NewEnvironment returns an environment on top of `enclosing`.
func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{enclosing: enclosing, values: make(map[string]interface{})}
}

// Define defines a new name(variable, function, class) in the calling Environment.
// The caller need to make sure the name isn't defined twice.
func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

// Get gets the value of `name` in the calling Environment.
// The method panics if there is no such `name`.
func (e *Environment) Get(name *Token) interface{} {
	if val, ok := e.values[name.Lexeme]; ok == true {
		return val
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	panic(NewRuntimeError(name, "undefined variable '"+name.Lexeme+"'."))
}

// GetAt gets the value of `name` from the the calling env's `distance` ancestor.
func (e *Environment) GetAt(distance int, name *Token) interface{} {
	return e.ancestor(distance).Get(name)
}

func (e *Environment) ancestor(distance int) *Environment {
	var env = e

	for i := 0; i < distance; i++ {
		env = env.enclosing
	}
	return env
}

// Assign assigns `value` to `name`.
// This method panics if the name isn't defined yet.
func (e *Environment) Assign(name *Token, value interface{}) {
	if _, ok := e.values[name.Lexeme]; ok == true {
		e.values[name.Lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}
	panic(NewRuntimeError(name, "undefined variable '"+name.Lexeme+"'."))
}

// AssignAt assigns `value` to `name` to the calling env's `distance` ancestor.
func (e *Environment) AssignAt(distance int, name *Token, value interface{}) {
	e.ancestor(distance).Assign(name, value)
}
