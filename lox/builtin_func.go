package lox

// BuiltInFunc is the runtime representation of builtin functions
type BuiltInFunc struct {
	arity    int
	call     func(*LoxInstance, ...interface{}) interface{} // internal go function
	instance *LoxInstance
}

// NewBuiltinFunc returns a new built-in functions.
func NewBuiltinFunc(arity int, call func(*LoxInstance, ...interface{}) interface{}) *BuiltInFunc {
	return &BuiltInFunc{arity: arity, call: call, instance: nil}
}

// Arity returns the number of arguments it takes.
func (bf *BuiltInFunc) Arity() int {
	return bf.arity
}

// Call implements the Callable interface.
func (bf *BuiltInFunc) Call(interpreter *Interpreter, args ...interface{}) interface{} {
	// we actually omit the `interpreter` because we already known how
	// to convert it into a go function.
	return bf.call(bf.instance, args...)
}

// Bind is called when interpreting `Get` expression.
// It is rather
func (bf *BuiltInFunc) Bind(instance *LoxInstance) Callable {
	bf.instance = instance
	// return itself.
	return bf
}

func (bf *BuiltInFunc) String() string {
	return "<native function>"
}
