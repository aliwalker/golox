package lox

// LoxInstance represents a runtime object for lox instance.
type LoxInstance struct {
	class *LoxClass
	props map[string]interface{}
}

// NewLoxInstance returns a runtime object.
func NewLoxInstance(class *LoxClass) *LoxInstance {
	var props = map[string]interface{}{}
	return &LoxInstance{class: class, props: props}
}

// Get returns the requested field.
// If the requested field is a getter, find the method, execute it and return the value.
func (o *LoxInstance) Get(interpreter *Interpreter, name *Token) interface{} {
	val, ok := o.props[name.Lexeme]
	if ok {
		return val
	}

	meth := o.class.FindMethod(o, name.Lexeme)
	if meth != nil {
		return meth
	}

	get := o.class.FindGetter(o, name.Lexeme)
	if get != nil {
		return get.Call(interpreter, nil)
	}
	panic(NewRuntimeError(name, "undefined property."))
}

// Set sets a field to the given value.
func (o *LoxInstance) Set(name *Token, value interface{}) interface{} {
	o.props[name.Lexeme] = value
	return value
}

func (o *LoxInstance) String() string {
	return o.class.String() + " instance"
}
