package lox

// ObjectType is a general object type for lox. Every
// lox object will need to implement this interface.
type ObjectType interface {
	Get(*Interpreter, *Token) interface{}
	Set(*Interpreter, *Token, interface{}) interface{}
	String() string
}
