package eval

import (
	"fmt"
	"misc/calc/ast"
	"misc/calc/parser"
	"misc/calc/token"
)

var builtins = map[string]func([]interface{}) interface{}{
	"+":  funcAdd,
	"-":  funcSub,
	"*":  funcMul,
	"/":  funcDiv,
	"%":  funcMod,
	"=":  funcEq,
	"<":  funcLess,
	"<=": funcLessEq,
	">":  funcGreater,
	">=": funcGreaterEq,
	"<>": funcNotEq,
	"if": funcIf,
}

var variables = map[string]interface{}{}
var functions = map[string]ast.Node{}

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
	switch node := n.(type) {
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
		if fn, ok := builtins[node.Lit]; ok {
			return fn
		}
		if fn, ok := functions[node.Lit]; ok {
			//fmt.Println("found something for:", node.Lit)
			return fn
		}
		if n, ok := variables[node.Lit]; ok {
			return n
		}
		return node
	case *ast.Number:
		return node.Val
	case *ast.Operator:
		return builtins[node.Val]
	case *ast.DefineExpr:
		e.evalDefineExpr(node)
		return nil
	case *ast.PrintExpr:
		e.evalPrintExpr(node)
		return nil
	case *ast.SetExpr:
		e.evalSetExpr(node)
		return nil
	case *ast.Expression:
		//fmt.Println(node.Nodes)
		// Ya...this section is an utter mess but it's an attempt to get a
		// callable function working without scoping. It works but it's ugly
		// as hell
		x := node.Nodes[0]
		if i, ok := x.(*ast.Identifier); ok {
			//fmt.Println("ident:", i)
			var ok bool
			var fn ast.Node
			if fn, ok = functions[i.Lit]; ok {
				//fmt.Println("function")
				// and...here's the problem. What variables belong to the function?
				if len(node.Nodes) > 1 {
					variables["x"] = e.eval(node.Nodes[1])
				}
				return e.eval(fn)
			}
			//fmt.Println("not a function")
		}
		fn, ok := e.eval(node.Nodes[0]).(func([]interface{}) interface{})
		if !ok {
			e.file.AddError(node.Nodes[0].Pos(), "First element of an expression "+
				"must be a function.")
			return nil
		}
		//fmt.Println("building args list")
		args := make([]interface{}, 0) //len(node.Nodes[1:]))
		if len(node.Nodes) > 1 {
			for _, node := range node.Nodes[1:] {
				args = append(args, e.eval(node))
			}
		}
		//fmt.Println("calling fn with", len(args), "args")

		res := fn(args)
		if err, ok := res.(error); ok {
			e.file.AddError(node.Pos(), err)
		}
		//fmt.Println("res:", res)
		return res
	}
	return nil
}

func (e *evaluator) evalDefineExpr(d *ast.DefineExpr) {
	functions[d.Name] = d.Impl
	for _, arg := range d.Args {
		variables[arg] = nil
	}
}

func (e *evaluator) evalPrintExpr(p *ast.PrintExpr) {
	args := make([]interface{}, len(p.Nodes))
	for i, n := range p.Nodes {
		args[i] = e.eval(n)
	}
	fmt.Println(args...)
}

func (e *evaluator) evalSetExpr(s *ast.SetExpr) {
	variables[s.Name] = e.eval(s.Value)
}

func genFunc(fn func(a, b int) int, args []interface{}) interface{} {
	if len(args) < 1 {
		return nil
	}
	if len(args) < 2 {
		if i, ok := args[0].(int); ok {
			return i
		}
		return nil
	}
	var res int
	if i, ok := args[0].(int); ok {
		res = i
	}
	for _, x := range args[1:] {
		switch v := x.(type) {
		case int:
			res = fn(res, v)
		default:
			// maybe return something like:
			// errors.New("Function accepts numerical types only, got:", v)
			return nil
		}
	}
	return res
}

func convBool(b bool) int {
	if b {
		return 1
	}
	return 0
}

func funcAdd(args []interface{}) interface{} {
	return genFunc(func(a, b int) int { return a + b }, args)
}

func funcSub(args []interface{}) interface{} {
	return genFunc(func(a, b int) int { return a - b }, args)
}

func funcMul(args []interface{}) interface{} {
	return genFunc(func(a, b int) int { return a * b }, args)
}

func funcDiv(args []interface{}) interface{} {
	return genFunc(func(a, b int) int { return a / b }, args)
}

func funcMod(args []interface{}) interface{} {
	return genFunc(func(a, b int) int { return a % b }, args)
}

func funcEq(args []interface{}) interface{} {
	return genFunc(func(a, b int) int { return convBool(a == b) }, args)
}

func funcLess(args []interface{}) interface{} {
	return genFunc(func(a, b int) int { return convBool(a < b) }, args)
}

func funcLessEq(args []interface{}) interface{} {
	return genFunc(func(a, b int) int { return convBool(a <= b) }, args)
}

func funcGreater(args []interface{}) interface{} {
	return genFunc(func(a, b int) int { return convBool(a > b) }, args)
}

func funcGreaterEq(args []interface{}) interface{} {
	return genFunc(func(a, b int) int { return convBool(a >= b) }, args)
}

func funcNotEq(args []interface{}) interface{} {
	return genFunc(func(a, b int) int { return convBool(a != b) }, args)
}

func funcIf(args []interface{}) interface{} {
	if len(args) != 3 {
		return nil //should produce error
	}
	if eq, ok := args[0].(int); ok {
		if eq == 0 {
			return args[2]
		}
		return args[1]
	}
	return nil // also an error
}
