package eval

import (
	"fmt"
	"misc/calc/ast"
	"misc/calc/parser"
	"misc/calc/token"
)

func EvalExpr(expr string) interface{} {
	return EvalFile("", expr)
}

func EvalFile(fname, expr string) interface{} {
	f := token.NewFile(fname, expr)
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

type evaluator struct {
	file  *token.File
	scope *ast.Scope // current scope
}

func (e *evaluator) eval(n ast.Node) interface{} {
	if n == nil {
		return nil
	}
	switch node := n.(type) {
	case *ast.CompExpr:
		return e.evalCompExpr(node)
	case *ast.DefineExpr:
		e.evalDefineExpr(node)
	case *ast.File:
		var x interface{}
		for _, n := range node.Nodes {
			x = e.eval(n) // scoping seems like it should come into play here
			switch t := x.(type) {
			case *ast.Identifier:
				e.file.AddError(t.Pos(), "Unknown identifier: ", t.Lit)
				return nil
			}
		}
		return x
	case *ast.Identifier:
		//var n ast.Node
		//fmt.Println("Looking up:", node.Lit)
		//fmt.Println(e.scope)
		n := e.scope.Lookup(node.Lit)
		//fmt.Println(n)
		if n != nil {
			//fmt.Println("eval:", node.Lit, "is nil")
			return e.eval(n)
		}
		return nil
	case *ast.IfExpr:
		return e.evalIfExpr(node)
	case *ast.MathExpr:
		return e.evalMathExpr(node)
	case *ast.Number:
		return node.Val
		return nil
	case *ast.PrintExpr:
		e.evalPrintExpr(node)
		return nil
	case *ast.SetExpr:
		e.evalSetExpr(node)
		return nil
	case *ast.UserExpr:
		return e.evalUserExpr(node)
	}
	return nil
}

func (e *evaluator) evalCompExpr(ce *ast.CompExpr) interface{} {
	a, aok := e.eval(ce.A).(int)
	b, bok := e.eval(ce.B).(int)
	if !aok || !bok {
		return 0
	}
	switch ce.CompLit {
	case "<":
		return BtoI(a < b)
	case "<=":
		return BtoI(a <= b)
	case "<>":
		return BtoI(a != b)
	case ">":
		return BtoI(a > b)
	case ">=":
		return BtoI(a >= b)
	case "=":
		return BtoI(a == b)
	}
	return 0
}

func (e *evaluator) evalDefineExpr(d *ast.DefineExpr) {
	e.scope.Insert(d.Name, d)
}

func (e *evaluator) evalIfExpr(i *ast.IfExpr) interface{} {
	x := 0 // default to false
	x, _ = e.eval(i.Comp).(int)
	if x >= 1 {
		return e.eval(i.Then)
	}
	return e.eval(i.Else) // returns nil if no else clause
}

func (e *evaluator) evalMathExpr(m *ast.MathExpr) interface{} {
	switch m.OpLit {
	case "+":
		return e.evalMathFunc(m.ExprList, func(a, b int) int { return a + b })
	case "-":
		return e.evalMathFunc(m.ExprList, func(a, b int) int { return a - b })
	case "*":
		return e.evalMathFunc(m.ExprList, func(a, b int) int { return a * b })
	case "/":
		return e.evalMathFunc(m.ExprList, func(a, b int) int { return a / b })
	case "%":
		return e.evalMathFunc(m.ExprList, func(a, b int) int { return a % b })
	case "and":
		return e.evalMathFunc(m.ExprList,
			func(a, b int) int { return BtoI(ItoB(a) && ItoB(b)) })
	case "or":
		return e.evalMathFunc(m.ExprList,
			func(a, b int) int { return BtoI(ItoB(a) || ItoB(b)) })
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

func BtoI(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ItoB(i int) bool {
	return i != 0
}

func (e *evaluator) evalUserExpr(u *ast.UserExpr) interface{} {
	n := e.scope.Lookup(u.Name)
	d, _ := n.(*ast.DefineExpr)
	tmp := e.scope
	e.scope = d.Scope
	for i, a := range d.Args {
		if len(u.Nodes) <= i {
			break
		}
		e.scope.Insert(a, u.Nodes[i])
	}
	// TODO: this is kind of hokey...
	var r interface{}
	for _, i := range d.Impl {
		r = e.eval(i)
	}
	e.scope = tmp
	return r
}
