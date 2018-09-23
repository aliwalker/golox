package lox

// LoxClass is a runtime object for a lox class.
type LoxClass struct {
	Name    string
	Super   *LoxClass
	Statics map[string]Callable // static props or functions.
	Methods map[string]Callable // class methods.
	Getters map[string]Callable // getters are functions in essence.
	Setters map[string]Callable // setters are functions in essence.
}

// NewLoxClass returns a runtime object for a class
func NewLoxClass(name string,
	super *LoxClass,
	statics,
	methods, getters, setters map[string]Callable) *LoxClass {

	return &LoxClass{
		Name:    name,
		Super:   super,
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

func (c *LoxClass) Bind(instance *LoxInstance) Callable {
	panic(NewRuntimeError(nil, "unable to bind a Lox class!"))
}

// Call returns a LoxInstance. It's a factory.
func (c *LoxClass) Call(i *Interpreter, args ...interface{}) interface{} {
	instance := NewLoxInstance(c)

	if init, ok := c.Methods["init"]; ok {
		initializer, _ := init.(Callable)
		initializer.Bind(instance).Call(i, args...)
	}

	return instance
}

// FindStatic returns the requested static method of the class.
func (c *LoxClass) FindStatic(name string) Callable {
	if val, ok := c.Statics[name]; ok {
		return val
	}
	return nil
}

// TODO: fix Find...

// FindMethod returns a binded method.
func (c *LoxClass) FindMethod(instance *LoxInstance, name string) Callable {
	if fn := findFunction(c.Methods, instance, name); fn != nil {
		return fn
	}

	if c.Super != nil {
		return c.Super.FindMethod(instance, name)
	}

	return nil
}

// FindGetter returns a binded getter.
func (c *LoxClass) FindGetter(instance *LoxInstance, name string) Callable {
	if fn := findFunction(c.Getters, instance, name); fn != nil {
		return fn
	}

	if c.Super != nil {
		return c.Super.FindGetter(instance, name)
	}

	return nil
}

// FindSetter returns a binded setter.
func (c *LoxClass) FindSetter(instance *LoxInstance, name string) Callable {
	if fn := findFunction(c.Setters, instance, name); fn != nil {
		return fn
	}

	if c.Super != nil {
		return c.Super.FindSetter(instance, name)
	}

	return nil
}

// findFunction finds a specific function, binds it to `instance` & returns the newly binded function.
func findFunction(fields map[string]Callable, instance *LoxInstance, name string) Callable {
	if fields == nil {
		return nil
	}

	if f, ok := fields[name]; ok {
		fn, _ := f.(Callable)
		return fn.Bind(instance)
	}

	return nil
}

func (c *LoxClass) String() string {
	return c.Name
}
