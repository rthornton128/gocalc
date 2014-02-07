// Copyright (c) 2013, Rob Thornton
// All rights reserved.
// This software is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package eval

import (
	"fmt"
	"github.com/rthornton128/gocalc/ast"
	"github.com/rthornton128/gocalc/parser"
	"github.com/rthornton128/gocalc/token"
	"strconv"
)

func EvalExpr(expr string) interface{} {
	return EvalFile("", expr)
}

func EvalFile(fname, expr string) interface{} {
	f := token.NewFile(fname, expr, 1)
	n := parser.ParseFile(f, expr)
	if f.NumErrors() > 0 {
		f.PrintErrors()
		return nil
	}
	e := &evaluator{file: f, scope: n.Scope}
	res := e.eval(n)
	if f.NumErrors() > 0 {
		f.PrintErrors()
		return nil
	}
	return res
}

func EvalDirectory(path string) {
}

type evaluator struct {
	file  *token.File
	scope *ast.Scope // current scope
}

/* Scope */
func (e *evaluator) openScope() {
	e.scope = ast.NewScope(e.scope)
}

func (e *evaluator) closeScope() {
	e.scope = e.scope.Parent
}

/* Evaluation */
func (e *evaluator) eval(n interface{}) interface{} {
	if n == nil {
		return nil
	}
	switch node := n.(type) {
	case *ast.CaseExpr:
		return e.evalCaseExpr(node)
	case *ast.CompExpr:
		return e.evalCompExpr(node)
	case *ast.ConcatExpr:
		return e.evalConcatExpr(node)
	case *ast.DefineExpr:
		e.evalDefineExpr(node)
	case *ast.File:
		var x interface{}
		for _, n := range node.Nodes {
			x = e.eval(n)
			switch t := x.(type) {
			case *ast.Identifier:
				e.file.AddError(t.Pos(), "Unknown identifier: ", t.Lit)
				return nil
			}
		}
		return x
	case *ast.Identifier:
		return e.eval(e.scope.Lookup(node.Lit))
	case *ast.IfExpr:
		return e.evalIfExpr(node)
	case *ast.MathExpr:
		return e.evalMathExpr(node)
	case *ast.Number:
		return node.Val
	case *ast.PrintExpr:
		e.evalPrintExpr(node)
		return nil
	case *ast.SetExpr:
		e.evalSetExpr(node)
		return nil
	case *ast.String:
		return node.Lit[1 : len(node.Lit)-1]
	case *ast.SwitchExpr:
		e.evalSwitchExpr(node)
	case *ast.UserExpr:
		return e.evalUserExpr(node)
	default:
		return node
	}
	return nil // unreachable
}

func (e *evaluator) evalCaseExpr(ce *ast.CaseExpr) interface{} {
	if e.eval(ce.Nodes[0]) == 1 {
		for _, n := range ce.Nodes[1:] {
			e.eval(n)
		}
		return 1
	}
	return nil
}

func (e *evaluator) evalCompExpr(ce *ast.CompExpr) interface{} {
	a, aok := e.eval(ce.Nodes[0]).(int)
	b, bok := e.eval(ce.Nodes[1]).(int)
	if !aok || !bok {
		return 0
	}
	switch ce.CompLit {
	case "<":
		return btoi(a < b)
	case "<=":
		return btoi(a <= b)
	case "<>":
		return btoi(a != b)
	case ">":
		return btoi(a > b)
	case ">=":
		return btoi(a >= b)
	case "=":
		return btoi(a == b)
	}
	return 0
}

func (e *evaluator) evalConcatExpr(ce *ast.ConcatExpr) interface{} {
	s := ""
	for _, node := range ce.Nodes {
		switch t := node.(type) {
		case *ast.Number:
			s += t.Lit
		default:
			r := e.eval(t)
			switch t := r.(type) {
			case string:
				s += t
			case int:
				s += strconv.Itoa(t)
			}
		}
	}
	return s
}

func (e *evaluator) evalDefineExpr(d *ast.DefineExpr) {
	e.scope.Insert(d.Name, d)
}

func (e *evaluator) evalIfExpr(i *ast.IfExpr) interface{} {
	x, _ := e.eval(i.Nodes[0]).(int)
	if x >= 1 {
		return e.eval(i.Nodes[1])
	}
	return e.eval(i.Nodes[2]) // returns nil if no else clause
}

func (e *evaluator) evalMathExpr(m *ast.MathExpr) interface{} {
	switch m.OpLit {
	case "+":
		return e.evalMathFunc(m.Nodes, func(a, b int) int { return a + b })
	case "-":
		return e.evalMathFunc(m.Nodes, func(a, b int) int { return a - b })
	case "*":
		return e.evalMathFunc(m.Nodes, func(a, b int) int { return a * b })
	case "/":
		return e.evalMathFunc(m.Nodes, func(a, b int) int { return a / b })
	case "%":
		return e.evalMathFunc(m.Nodes, func(a, b int) int { return a % b })
	case "and":
		return e.evalMathFunc(m.Nodes,
			func(a, b int) int { return btoi(itob(a) && itob(b)) })
	case "or":
		return e.evalMathFunc(m.Nodes,
			func(a, b int) int { return btoi(itob(a) || itob(b)) })
	default:
		return nil // not reachable (fingers crossed!)
	}
}

func (e *evaluator) evalMathFunc(list []ast.Node, fn func(int, int) int) int {
	a, ok := e.eval(list[0]).(int)
	if !ok {
		return 0 // or should this return an error?
	}
	for _, n := range list[1:] {
		b, ok := e.eval(n).(int)
		if !ok {
			return 0
		}
		a = fn(a, b)
	}
	return a
}

func (e *evaluator) evalPrintExpr(p *ast.PrintExpr) {
	args := make([]interface{}, len(p.Nodes))
	for i, n := range p.Nodes {
		args[i] = e.eval(n)
	}
	fmt.Println(args...)
}

func (e *evaluator) evalSetExpr(s *ast.SetExpr) {
	e.scope.Insert(s.Name, s.Value)
}

func (e *evaluator) evalSwitchExpr(s *ast.SwitchExpr) {
	if s.Pred == nil {
		for _, n := range s.Nodes {
			if e.eval(n.(*ast.CaseExpr)) != nil {
				break
			}
		}
	} else {
		p := e.eval(s.Pred)
		for _, n := range s.Nodes {
			ce := n.(*ast.CaseExpr)
			if e.eval(ce.Nodes[0]) == p {
				for _, y := range ce.Nodes[1:] {
					e.eval(y)
				}
				break
			}
		}
	}
}

func (e *evaluator) evalUserExpr(u *ast.UserExpr) interface{} {
	n := e.scope.Lookup(u.Name)
	d, _ := n.(*ast.DefineExpr)
	e.openScope()
	args := make([]interface{}, len(d.Args))
	for i, _ := range args {
		if len(u.Nodes) <= i {
			break
		}
		args[i] = e.eval(u.Nodes[i])
	}
	for i, v := range args {
		e.scope.Insert(d.Args[i], v)
	}
	var r interface{}
	for _, v := range d.Nodes {
		r = e.eval(v)
		if r != nil {
			break
		}
	}
	e.closeScope()
	return r
}
