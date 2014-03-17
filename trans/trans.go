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
	fmt.Println("transexpr")
	TransFile(w, "", expr)
}

func TransFile(w io.Writer, fname, expr string) {
	f := token.NewFile(fname, expr, 1)
	n := parser.ParseFile(f, expr)
	
  if f.NumErrors() > 0 {
		f.PrintErrors()
		return
	}

	t := &translator{out: w, file: f} //, scope: n.Scope}
	t.topComment()
	/* includes will/might eventually reflect the imports from Calc. It's
	 * possible that stdio might be an auto-include if print remains a
	 * built-in function, which is likely won't */
	t.includes() /* temporary */
	t.transpile(n, false)

  if f.NumErrors() > 0 {
		f.PrintErrors()
	}
	return
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
	case *ast.DefineExpr:
		semi = false
		t.transDefineExpr(node)
	case *ast.File:
		for _, n := range node.Nodes {
			t.transpile(n, true)
		}
	case *ast.Identifier:
		t.write(node.Lit)
	case *ast.MathExpr:
		t.transMathExpr(node)
	case *ast.Number:
		t.write(node.Lit)
	case *ast.PrintExpr:
		t.transPrintExpr(node)
	case *ast.SetExpr:
		t.transSetExpr(node)
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

func (t *translator) returnStatement() {
	t.writeln("return 0;")
}

func (t *translator) transDefineExpr(de *ast.DefineExpr) {
	t.write("int ") /* temporary, type should be inferred */
	t.write(de.Name + "(")
	for i, a := range de.Args {
		t.write(a)
		if i < len(de.Args)-1 {
			t.write(",")
		}
	}
	if len(de.Args) == 0 {
		t.write("void")
	}
	t.write(")")
	t.openBlock()
	for _, n := range de.Nodes {
		t.transpile(n, true)
	}
	t.returnStatement()
	t.closeBlock()
	t.write("\n")
}

func (t *translator) transMathExpr(me *ast.MathExpr) {
	t.write("(")
	for i, n := range me.Nodes {
		t.transpile(n, false)
		if i < len(me.Nodes)-1 {
			t.write(me.OpLit)
		}
	}
	t.write(")")
}

func (t *translator) transPrintExpr(pe *ast.PrintExpr) {
	t.write("printf(\"%d\",")
	t.transpile(pe.Nodes[0], false)
	t.write(")")
}

func (t *translator) transSetExpr(se *ast.SetExpr) {
	t.write("int ") /* temp */
	t.write(se.Name + " = ")
	t.transpile(se.Value, false)
}
