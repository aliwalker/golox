package lox

// LoxClass is a runtime object for a lox class.
type LoxClass struct {
	Name    string
	Methods map[string]*LoxFunction
	Getters map[string]*LoxFunction // getters are function in essence.
}

// NewLoxClass returns a runtime object for a class
func NewLoxClass(name string, methods, getters map[string]*LoxFunction) *LoxClass {
	return &LoxClass{Name: name, Methods: methods, Getters: getters}
}

// Arity returns the number of args the initializer takes.
func (c *LoxClass) Arity() int {
	initializer, ok := c.Methods["init"]
	if ok {
		return initializer.Arity()
	}
	return 0
}

// Call returns a LoxInstance. It's a factory.
func (c *LoxClass) Call(i *Interpreter, args ...interface{}) interface{} {
	instance := NewLoxInstance(c)
	initializer, ok := c.Methods["init"]

	if ok {
		initializer.Bind(instance).Call(i, args...)
	}
	return instance
}

// FindMethod returns the requested method.
func (c *LoxClass) FindMethod(instance *LoxInstance, name string) *LoxFunction {
	return findFunction(c.Methods, instance, name)
}

// FindGetter returns the requested getter.
func (c *LoxClass) FindGetter(instance *LoxInstance, name string) *LoxFunction {
	return findFunction(c.Getters, instance, name)
}

func findFunction(fields map[string]*LoxFunction, instance *LoxInstance, name string) *LoxFunction {
	if fields == nil {
		return nil
	}

	fn, ok := fields[name]
	if ok {
		return fn.Bind(instance)
	}
	return nil
}

func (c *LoxClass) String() string {
	return c.Name
}
