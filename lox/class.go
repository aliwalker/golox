package lox

type LoxClass struct {
	Name string
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{Name: name}
}

func (c *LoxClass) Arity() int {
	return 0
}

func (c *LoxClass) Call(i *Interpreter, args ...interface{}) interface{} {
	return NewLoxInstance(c)
}

func (c *LoxClass) String() string {
	return c.Name
}
