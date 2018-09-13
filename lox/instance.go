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
func (o *LoxInstance) Get(name *Token) interface{} {
	val, ok := o.props[name.Lexeme]

	if ok != true {
		panic(NewRuntimeError(name, "undefined property."))
	}
	return val
}

// Set sets a field to the given value.
func (o *LoxInstance) Set(name *Token, value interface{}) interface{} {
	o.props[name.Lexeme] = value
	return value
}

func (o *LoxInstance) String() string {
	return o.class.String() + " instance"
}
