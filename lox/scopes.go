package lox

// Scopes & Scope are used by Resolver.

// variable status.
type varStatus int

const (
	_             varStatus = iota
	varUndeclared           // variable is not declared.
	varDeclared             // variable is declared but not available for reference.
	varDefined              // variable is available for reference.
)

// Scope is a map with key = name of a variable, indicates the status of the variable.
type Scope map[string]varStatus

// HasName checks whether `name` is declared in `scope`.
func (scope Scope) HasName(name string) bool {
	val := scope[name]

	if val == varUndeclared {
		return false
	}
	return true
}

// Scopes is a scope stack.
// All methods of Scopes might panic.
type Scopes []Scope

// NewScopes creates and returns a scope stack.
func NewScopes() *Scopes {
	return &Scopes{}
}

// Len returns the number of scopes in the scope stack.
func (scopes *Scopes) Len() int {
	return len(*scopes)
}

// Peek returns the top most Scope from a Scope stack.
func (scopes *Scopes) Peek() Scope {
	return (*scopes)[len(*scopes)-1]
}

// Push appends a new scope to scopes.
func (scopes *Scopes) Push() {
	*scopes = append(*scopes, Scope{})
}

// Pop removes the current topmost scope.
func (scopes *Scopes) Pop() {
	*scopes = (*scopes)[:len(*scopes)-1]
}

// Get return the i-th scope from scope stack.
func (scopes *Scopes) Get(i int) Scope {
	return (*scopes)[i]
}

// Empty checks whether the scope stack is empty.
func (scopes *Scopes) Empty() bool {
	return len(*scopes) == 0
}
