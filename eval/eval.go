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
		var n ast.Node
		//fmt.Println("Looking up:", node.Lit)
		n = e.scope.Lookup(node.Lit)
		if n != nil {
			return e.eval(n)
		}
		return nil
	case *ast.IfExpr:
		return e.evalIfExpr(node)
	case *ast.MathExpr:
		return e.evalMathExpr(node)
	case *ast.Number:
		return node.Val
	case *ast.DefineExpr:
		e.evalDefineExpr(node)
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
		return convBool(a < b)
	case "<=":
		return convBool(a <= b)
	case "<>":
		return convBool(a != b)
	case ">":
		return convBool(a > b)
	case ">=":
		return convBool(a >= b)
	case "=":
		return convBool(a == b)
	}
	return 0
}

func (e *evaluator) evalDefineExpr(d *ast.DefineExpr) {
	// TODO: replace with proper scoping code
	//functions[d.Name] = d.Impl
	fmt.Print("define ", d.Name, " with args:")
	e.scope.Insert(d.Name, d.Impl)
	for _, arg := range d.Args {
		//variables[arg] = nil
		fmt.Print(arg)
		e.scope.Insert(arg, nil)
	}
	fmt.Println()
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

func convBool(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (e *evaluator) evalUserExpr(u *ast.UserExpr) interface{} {
	n := e.scope.Lookup(u.Name)
	return e.eval(n)
}
