package lox

// LoxClass is a runtime object for a lox class.
type LoxClass struct {
	Name    string
	Methods map[string]*LoxFunction
}

// NewLoxClass returns a runtime object for a class
func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{Name: name, Methods: methods}
}

// Arity returns the number of args the initializer takes.
// TODO: add initializer.
func (c *LoxClass) Arity() int {
	return 0
}

// Call returns a LoxInstance. It's a factory.
func (c *LoxClass) Call(i *Interpreter, args ...interface{}) interface{} {
	return NewLoxInstance(c)
}

// FindMethod returns the requested method.
func (c *LoxClass) FindMethod(instance *LoxInstance, name string) interface{} {
	if c.Methods == nil {
		return nil
	}

	fn, ok := c.Methods[name]
	if ok {
		return fn.Bind(instance)
	}
	return nil
}

func (c *LoxClass) String() string {
	return c.Name
}
