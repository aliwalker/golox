package lox

// ControlType identifies the Control Stmt, which can be return stmt or break stmt.
type ControlType int

const (
	_             ControlType = iota
	ControlReturn             // return control
	ControlBreak              // break control
)
