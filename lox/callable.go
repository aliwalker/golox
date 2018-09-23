package lox

type Callable interface {
	Call(interpreter *Interpreter, args ...interface{}) interface{}
	Arity() int
	Bind(instance *LoxInstance) Callable
}
