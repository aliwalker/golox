package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: generate_ast <output_file>")
		os.Exit(1)
	}

	out, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	defineAst(out, "Expr", []string{
		"Assign		: Name *Token, Operator *Token, Value Expr",
		"Binary		: Left Expr, Operator *Token, Right Expr",
		"Call		: Callee Expr, Paren *Token, Arguments []Expr",
		"Grouping	: Expression Expr",
		"Literal	: Value interface{}",
		"Logical	: Left Expr, Operator *Token, Right Expr",
		"Unary		: Operator *Token, Right Expr",
		"Variable	: Name *Token",
	})

	defineAst(out, "Stmt", []string{
		"Block		: Stmts []Stmt",
		"Control	: CtrlType ControlType, Value Expr",
		"Function	: Name *Token, Params []*Token, Body []Stmt",
		"Expression	: Expression Expr",
		"If			: Condition Expr, ThenBranch Stmt, ElseBranch Stmt",
		"Print		: Expression Expr",
		"Var		: Name *Token, Initializer Expr",
		"While		: Condition Expr, Body Stmt",
	})
}

func defineAst(out, base string, types []string) {
	var src string

	src += fmt.Sprintln("")
	src += fmt.Sprintln("package lox")
	src += fmt.Sprintln("")

	src += defineVisitor(base, types)

	for _, t := range types {
		klass := strings.TrimRight(strings.Split(t, ":")[0], "\t")
		fields := strings.TrimRight(strings.Split(t, ":")[1], " ")
		src += defineType(base, klass, fields)
	}

	path := fmt.Sprintf("%s/%s.go", out, strings.ToLower(base))
	if err := saveFile(path, src); err != nil {
		panic(err)
	}
}

func defineVisitor(base string, fields []string) string {
	var src string

	src += fmt.Sprintln("")
	src += fmt.Sprintf("type %sVisitor interface {\n", base)
	for _, field := range fields {
		klass := strings.TrimRight(strings.Split(field, ":")[0], "\t")
		src += fmt.Sprintf("Visit%s%s(%s *%s) interface{}", klass, base, strings.ToLower(base), klass)
		src += fmt.Sprintln("")
	}
	src += fmt.Sprintln("}")

	src += fmt.Sprintln("")
	src += fmt.Sprintf("type %s interface {", base)
	src += fmt.Sprintln("")
	src += fmt.Sprintf("Accept(v %sVisitor) interface{}\n", base)
	src += fmt.Sprintln("}")

	return src
}

func defineType(base, klass, fields string) string {
	var src string

	src += fmt.Sprintln("")
	src += fmt.Sprintf("type %s struct {", klass)
	src += fmt.Sprintln("")

	flds := strings.Split(fields, ",")
	for _, fld := range flds {
		src += fmt.Sprintln(fld)
	}
	src += fmt.Sprintln("}")

	// NewXxx
	src += fmt.Sprintf("func New%s(", klass)
	params := []string{}
	for _, fld := range flds {
		t := (strings.Split(fld, " "))[2]
		name := strings.ToLower(strings.Split(fld, " ")[1])
		params = append(params, fmt.Sprintf("%s %s", name, t))
	}
	src += fmt.Sprintf(strings.Join(params, ", "))
	src += fmt.Sprintf(") %s {", base)
	src += fmt.Sprintln("")
	src += fmt.Sprintf("return &%s{", klass)

	args := []string{}
	for _, fld := range flds {
		t := strings.ToLower(strings.Split(fld, " ")[1])
		name := strings.Split(fld, " ")[1]
		args = append(args, fmt.Sprintf("%s: %s", name, t))
	}
	src += fmt.Sprintf(strings.Join(args, ","))
	src += fmt.Sprintln("}")
	src += fmt.Sprintln("}")

	src += fmt.Sprintf("func (expr *%s) Accept(v %sVisitor) interface{} {", klass, base)
	src += fmt.Sprintln("")
	src += fmt.Sprintf("return v.Visit%s%s(expr)", klass, base)
	src += fmt.Sprintln("}")
	src += fmt.Sprintln("")

	return src
}

func saveFile(path, src string) error {
	buf, err := format.Source([]byte(src))
	if err != nil {
		return err
	}

	ioutil.WriteFile(path, buf, 0644)
	return nil
}
