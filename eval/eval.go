package eval

import (
	"fmt"
	"misc/calc/ast"
	"misc/calc/parser"
)

var builtins = map[string]func([]interface{}) interface{}{
	"+":     funcAdd,
	"-":     funcSub,
	"*":     funcMul,
	"/":     funcDiv,
	"%":     funcMod,
	"print": funcPrint,
	"set":   funcSet,
}

var variables = map[string]interface{}{}

func EvalExpr(expr string) interface{} {
	return eval(parser.ParseExpr(expr))
}

//func EvalFile(str string) interface{}

func eval(n ast.Node) interface{} {
	switch node := n.(type) {
	case *ast.File:
		var x interface{}
		fmt.Println("File type; any nodes?")
		for _, n := range node.Nodes {
			fmt.Println("evaluating nodes...")
			x = eval(n) // scoping seems like it should come into play here
			switch t := x.(type) {
			case *ast.Identifier:
				fmt.Println("Error - Unknown identifier:", t.Lit)
				return nil
			default:
				fmt.Printf("TYPE: %T\n", n)
			}
		}
		return x
	case *ast.Identifier:
		if fn, ok := builtins[node.Lit]; ok {
			return fn
		}
		if n, ok := variables[node.Lit]; ok {
			return n
		}
		return node
	case *ast.Number:
		return node.Val
	case *ast.Operator:
		// need to produce error if no built-in
		if fn, ok := builtins[string(node.Val)]; ok {
			return fn
		}
		return nil
	case *ast.Expression:
		if len(node.Nodes) < 2 {
			return nil
		}
		fn, _ := eval(node.Nodes[0]).(func([]interface{}) interface{})
		args := make([]interface{}, len(node.Nodes[1:]))
		for i, node := range node.Nodes[1:] {
			args[i] = eval(node)
		}
		//fmt.Println("calling fn with", len(args), "args")

		// the following line could be changed to:
		// res := fn(args)
		// if err, ok := res.(error); ok {
		//   file.AddError(node.BegPos(), err)
		// }
		// return res
		return fn(args)
	}
	return nil
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

func funcPrint(args []interface{}) interface{} {
	// some checks should be done on the args. For example, this current
	// implementation will return the address of a built-in function if
	// given as an argument.
	fmt.Println(args...)
	return nil
}

func funcSet(args []interface{}) interface{} {
	if len(args) != 2 {
		return nil // really feel like this should be an error...not just nil
	}
	if i, ok := args[0].(*ast.Identifier); ok {
		switch args[1].(type) {
		case *ast.Operator:
			return nil // this REALLY should produce an error...
		default:
			variables[i.Lit] = args[1]
		}
	}
	return nil
}
