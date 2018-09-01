package lox

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{enclosing: enclosing, values: make(map[string]interface{})}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name *Token) interface{} {
	if val, ok := e.values[name.Lexeme]; ok == true {
		return val
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	panic(NewRuntimeError(name, "undefined variable '"+name.Lexeme+"'."))
}

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
