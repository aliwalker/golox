package lox

// Callable represents Lox callable object at runtime.
type Callable interface {
	Call(interpreter *Interpreter, args ...interface{}) interface{}
	Arity() int
	Bind(instance *LoxInstance) Callable
}
