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

func (f *LoxFunction) Call(interpreter *Interpreter, arguments ...interface{}) (returnVal interface{}) {
	enclosingEnv := interpreter.environment // for return usage.
	env := NewEnvironment(f.Enclosing)

	defer func() {
		if val := recover(); val != nil {
			returnControl, ok := val.(*Control)
			// If it didn't catch a Control, repanic.
			if ok != true {
				panic(val)
			}

			// If it is not a return control, repanic.
			if returnControl.CtrlType != ControlReturn {
				panic(val)
			}

			returnVal = interpreter.evaluate(returnControl.Value)
			interpreter.environment = enclosingEnv
		}
	}()

	for i, param := range f.Declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.Declaration.Body, env)
	return
}

func (f *LoxFunction) String() string {
	if f.Declaration.Name == nil {
		return "<fn " + "lambda" + ">"
	}
	return "<fn " + f.Declaration.Name.Lexeme + ">"
}
