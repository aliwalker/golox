package lox

// LoxClass is a runtime object for a lox class.
type LoxClass struct {
	Name    string
	Statics map[string]*LoxFunction // static props or functions.
	Methods map[string]*LoxFunction // class methods.
	Getters map[string]*LoxFunction // getters are functions in essence.
	Setters map[string]*LoxFunction // setters are functions in essence.
}

// NewLoxClass returns a runtime object for a class
func NewLoxClass(name string,
	statics map[string]*LoxFunction,
	methods, getters, setters map[string]*LoxFunction) *LoxClass {

	return &LoxClass{
		Name:    name,
		Statics: statics,
		Methods: methods,
		Getters: getters,
		Setters: setters,
	}
}

// Arity returns the number of args the initializer takes.
func (c *LoxClass) Arity() int {
	if initializer, ok := c.Methods["init"]; ok {
		return initializer.Arity()
	}
	return 0
}

// Call returns a LoxInstance. It's a factory.
func (c *LoxClass) Call(i *Interpreter, args ...interface{}) interface{} {
	instance := NewLoxInstance(c)

	if initializer, ok := c.Methods["init"]; ok {
		initializer.Bind(instance).Call(i, args...)
	}

	return instance
}

// FindStatic returns the requested static method of the class.
func (c *LoxClass) FindStatic(name string) *LoxFunction {
	if val, ok := c.Statics[name]; ok {
		return val
	}
	return nil
}

// FindMethod returns a binded method.
func (c *LoxClass) FindMethod(instance *LoxInstance, name string) *LoxFunction {
	return findFunction(c.Methods, instance, name)
}

// FindGetter returns a binded getter.
func (c *LoxClass) FindGetter(instance *LoxInstance, name string) *LoxFunction {
	return findFunction(c.Getters, instance, name)
}

// FindSetter returns a binded setter.
func (c *LoxClass) FindSetter(instance *LoxInstance, name string) *LoxFunction {
	return findFunction(c.Setters, instance, name)
}

// findFunction finds a specific function, binds it to `instance` & returns the newly binded function.
func findFunction(fields map[string]*LoxFunction, instance *LoxInstance, name string) *LoxFunction {
	if fields == nil {
		return nil
	}

	if fn, ok := fields[name]; ok {
		return fn.Bind(instance)
	}

	return nil
}

func (c *LoxClass) String() string {
	return c.Name
}
