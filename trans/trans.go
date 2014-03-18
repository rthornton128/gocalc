// Copyright (c) 2013, Rob Thornton
// All rights reserved.
// This software is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package trans

import (
	"fmt"
	"io"

	"github.com/rthornton128/gocalc/ast"
	"github.com/rthornton128/gocalc/parser"
	"github.com/rthornton128/gocalc/token"
)

type translator struct {
	out   io.Writer
	file  *token.File
	scope *ast.Scope
}

/* TransExpr is really only for initial testing and will probably be removed
 * in the near future */
func TransExpr(w io.Writer, expr string) {
	TransFile(w, "", expr)
}

func TransFile(w io.Writer, fname, expr string) {
	f := token.NewFile(fname, expr, 1)
	n := parser.ParseFile(f, expr)

	if f.NumErrors() > 0 {
		f.PrintErrors()
		return
	}

	t := &translator{out: w, file: f, scope: n.Scope}
	t.topComment()
	/* includes will/might eventually reflect the imports from Calc. It's
	 * possible that stdio might be an auto-include if print remains a
	 * built-in function, which is likely won't */
	t.includes() /* temporary */
	t.transFuncSigs(n)
	t.transpile(n, false)

	if f.NumErrors() > 0 {
		f.PrintErrors()
	}

	if t.scope.Lookup("main") == nil {
		fmt.Println("No function \"main\" found!")
	}
	return
}

func (t *translator) nodeType(n ast.Node) string {
	switch node := n.(type) {
	case *ast.CompExpr, *ast.Number, *ast.MathExpr:
		return "int"
	case *ast.String, *ast.ConcatExpr:
		return "char *"
	case *ast.DefineExpr:
		return t.nodeType(node.Nodes[len(node.Nodes)-1])
	case *ast.Identifier:
		if x := t.scope.Lookup(node.Lit); x != nil {
			//fmt.Println(x.(ast.Node))
			return t.nodeType(x.(ast.Node))
		}
		return "void *"
	case *ast.IfExpr:
		return t.nodeType(node.Nodes[2])
	case *ast.UserExpr:
		return t.nodeType(node.Nodes[len(node.Nodes)-1])
	default:
		return "void *"
	}
}

/* Scope */
func (t *translator) openScope() {
	t.scope = ast.NewScope(t.scope)
}

func (t *translator) closeScope() {
	t.scope = t.scope.Parent
}

/* Transpiler */
func (t *translator) transpile(n ast.Node, semi bool) {
	switch node := n.(type) {
	case *ast.CompExpr:
		t.transCompExpr(node)
	case *ast.DefineExpr:
		semi = false
		t.transDefineExpr(node)
	case *ast.File:
		for _, n := range node.Nodes {
			t.transpile(n, true)
		}
	case *ast.Identifier:
		t.write(node.Lit)
	case *ast.IfExpr:
		t.transIfExpr(node)
	case *ast.MathExpr:
		t.transMathExpr(node)
	case *ast.Number:
		t.write(node.Lit)
	case *ast.PrintExpr:
		t.transPrintExpr(node)
	case *ast.SetExpr:
		t.transSetExpr(node)
	case *ast.String:
		t.write(node.Lit)
	case *ast.UserExpr:
		t.transUserExpr(node)
	}
	if semi {
		t.write(";\n")
	}
}

func (t *translator) write(s string) {
	if _, err := t.out.Write([]byte(s)); err != nil {
		panic(err)
	}
}

func (t *translator) writeln(s string) {
	if _, err := t.out.Write([]byte(s + "\n")); err != nil {
		panic(err)
	}
}

func (t *translator) topComment() {
	t.writeln("/* This program was created by Translitorator 2000 v0.1 */")
}

func (t *translator) includes() {
	t.writeln("#include <stdio.h>")
}

func (t *translator) openBlock() {
	t.writeln("{")
}

func (t *translator) closeBlock() {
	t.writeln("}")
}

func (t *translator) returnStatement(n ast.Node) {
	t.write("return ")
	t.transpile(n, true)
}

func (t *translator) transCompExpr(ce *ast.CompExpr) {
	t.transpile(ce.Nodes[0], false)
	if ce.CompLit == "=" {
		t.write(" == ")
	} else {
		t.write(" " + ce.CompLit + " ")
	}
	t.transpile(ce.Nodes[1], false)
}

func (t *translator) transDefineExpr(de *ast.DefineExpr) {
	t.scope.Insert(de.Name, de)
	t.openScope()
	t.transFuncDecl(de)
	t.openBlock()
	for i := 0; i < len(de.Nodes)-1; i++ {
		t.transpile(de.Nodes[i], true)
	}
	last := de.Nodes[len(de.Nodes)-1]
	if n, ok := last.(*ast.IfExpr); ok {
		t.transpile(n, false)
	} else {
		t.returnStatement(last)
	}
	t.closeBlock()
	t.closeScope()
	t.write("\n")
}

func (t *translator) transFuncDecl(de *ast.DefineExpr) {
	t.write(t.nodeType(de) + " ")
	t.write(de.Name + "(")
	for i, a := range de.Args {
		t.write("int ") /* again, so bad... */
		t.write(a)
		t.scope.Insert(a, 0)
		if i < len(de.Args)-1 {
			t.write(",")
		}
	}
	if len(de.Args) == 0 {
		t.write("void")
	}
	t.write(")")
}

func (t *translator) transFuncSigs(n ast.Node) {
	/* walk the AST and create function signatures for every function found */
	switch node := n.(type) {
	case *ast.DefineExpr:
		t.transFuncDecl(node)
		t.write(";\n")
	case *ast.File:
		for _, v := range node.Nodes {
			t.transFuncSigs(v)
		}
	case *ast.Expression:
		for _, v := range node.Nodes {
			t.transFuncSigs(v)
		}
	}
	return
}

func (t *translator) transIfExpr(ie *ast.IfExpr) {
	t.write("if (")
	t.transpile(ie.Nodes[0], false)
	t.write(")")
	t.openBlock()
	switch ie.Nodes[1].(type) {
	case *ast.Identifier, *ast.Number, *ast.String, *ast.MathExpr, *ast.UserExpr:
		t.write("return ")
	}
	t.transpile(ie.Nodes[1], true)
	t.closeBlock()
	if ie.Nodes[2] != nil {
		t.write("else")
		t.openBlock()
		switch ie.Nodes[2].(type) {
		case *ast.Number, *ast.String, *ast.MathExpr, *ast.UserExpr:
			t.write("return ")
		}
		t.transpile(ie.Nodes[2], true)
		t.closeBlock()
	}
}

func (t *translator) transMathExpr(me *ast.MathExpr) {
	t.write("(")
	for i, n := range me.Nodes {
		t.transpile(n, false)
		if i < len(me.Nodes)-1 {
			switch me.OpLit {
			case "or":
				t.write("||")
			case "and":
				t.write("&&")
			default:
				t.write(me.OpLit)
			}
		}
	}
	t.write(")")
}

func (t *translator) transPrintExpr(pe *ast.PrintExpr) {
	t.write("printf(\"")
	for i, n := range pe.Nodes {
		switch t.nodeType(n) {
		case "int":
			t.write("%d")
		case "char *":
			t.write("%s")
		case "void *":
			t.write("%p")
		}
		if i < len(pe.Nodes)-1 {
			t.write(" ")
		}
	}
	if len(pe.Nodes) > 0 {
		t.write("\\n\",")
		for _, n := range pe.Nodes {
			t.transpile(n, false)
		}
	} else {
		t.write("\"")
	}
	t.write(")")
}

func (t *translator) transSetExpr(se *ast.SetExpr) {
	t.scope.Insert(se.Name, se.Value)
	t.write(t.nodeType(se.Value) + " ")
	t.write(se.Name + " = ")
	t.transpile(se.Value, false)
}

func (t *translator) transUserExpr(ue *ast.UserExpr) {
	t.write(ue.Name + "(")
	for i, v := range ue.Nodes {
		t.transpile(v, false)
		if i < len(ue.Nodes)-1 {
			t.write(",")
		}
	}
	t.write(")")
}
