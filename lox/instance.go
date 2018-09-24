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
	// property
	if val, ok := o.props[name.Lexeme]; ok {
		return val
	}

	// method
	if meth := o.class.FindMethod(o, name.Lexeme); meth != nil {
		return meth
	}

	// getter
	if get := o.class.FindGetter(o, name.Lexeme); get != nil {
		return get.Call(interpreter, nil)
	}

	panic(NewRuntimeError(name, "undefined property."))
}

// Set sets a field to the given value.
func (o *LoxInstance) Set(interpreter *Interpreter, name *Token, value interface{}) interface{} {
	// setter.
	if set := o.class.FindSetter(o, name.Lexeme); set != nil {
		return set.Call(interpreter, value)
	}

	// property.
	o.props[name.Lexeme] = value
	return value
}

func (o *LoxInstance) String() string {
	stringfied := "[" + o.class.String() + " instance" + "] {\n"

	for name, value := range o.props {
		valueStr, _ := value.(string)
		stringfied += "\t" + name + ": " + valueStr + "\n"
	}

	stringfied += "}"
	return stringfied
}
