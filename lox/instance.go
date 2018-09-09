package lox

// LoxInstance represents a runtime object for lox instance.
type LoxInstance struct {
	class *LoxClass
}

// NewLoxInstance returns a runtime object.
func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class: class}
}

func (c *LoxInstance) String() string {
	return c.class.String() + " instance"
}
