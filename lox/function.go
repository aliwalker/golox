package lox

type LoxFunction struct {
	Declaration *Function
	Enclosing   *Environment
}

func NewLoxFunction(declaration *Function, enclosing *Environment) *LoxFunction {
	return &LoxFunction{declaration, enclosing}
}

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments ...interface{}) interface{} {
	env := NewEnvironment(f.Enclosing)

	for i, param := range f.Declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.Declaration.Body, env)
	// TODO: add return value when return statement is implemented.
	return nil
}

func (f *LoxFunction) toString() string {
	return "<fn " + f.Declaration.Name.Lexeme + ">"
}
