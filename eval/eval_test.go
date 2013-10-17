package eval_test

import (
	"misc/calc/eval"
	"testing"
)

func TestEvalAddition(t *testing.T) {
	var tests = []struct {
		expr string
		res  int
	}{
		{"(+ 2)", 2},
		{"(+ 1 2)", 3},
		{"(+ 1 2 3)", 6},
		{"(+ 1 2 -3)", 0},
		{"(+ 1 (+ 2 3))", 6},
		{"(+ (+ 1 2) (+ 3 4))", 10},
	}
	for x, test := range tests {
		res := eval.EvalExpr(test.expr)
		i, ok := res.(int)
		if !ok || i != test.res {
			t.Log(x, "- Expected:", test.res)
			t.Fatal(x, "- Got:", i)
		}
	}
}

func TestEvalSubtraction(t *testing.T) {
	var tests = []struct {
		expr string
		res  int
	}{
		{"(- 1)", 1}, // funny behaviour...maybe should error on 1 argument?
		{"(- 2 4)", -2},
		{"(- -2 4)", -6},
		{"(- 4 (- 2 1))", 3},
		{"(- (- 8 2) (- 3 4))", 7},
	}
	for x, test := range tests {
		res := eval.EvalExpr(test.expr)
		i, ok := res.(int)
		if !ok || i != test.res {
			t.Log(x, "- Expected:", test.res)
			t.Fatal(x, "- Got:", i)
		}
	}
}

func TestEvalMultiplication(t *testing.T) {
	var tests = []struct {
		expr string
		res  int
	}{
		{"(* 1)", 1},
		{"(* 1 2)", 2},
		{"(* -1 -2)", 2},
		{"(* 2 (* 3 4))", 24},
		{"(* (* 2 3) (* 3 3))", 54},
	}
	for x, test := range tests {
		res := eval.EvalExpr(test.expr)
		i, ok := res.(int)
		if !ok || i != test.res {
			t.Log(x, "- Expected:", test.res)
			t.Fatal(x, "- Got:", i)
		}
	}
}

func TestEvalDivision(t *testing.T) {
	var tests = []struct {
		expr string
		res  int
	}{
		{"(/ 1)", 1},
		{"(/ 4 2)", 2},
		{"(/ 8 (/ 4 2))", 4},
		{"(/ (/ 1 1) (/ 3 3))", 1},
	}
	for x, test := range tests {
		res := eval.EvalExpr(test.expr)
		i, ok := res.(int)
		if !ok || i != test.res {
			t.Log(x, "- Expected:", test.res)
			t.Fatal(x, "- Got:", i)
		}
	}
}

func TestEvalSet(t *testing.T) {
	var tests = []struct {
		expr string
		res  interface{}
	}{
		{"(set a 1)", nil},
		{"(set expr (+ 2 3))", nil},
		{"a", nil},
		{"expr", nil},
		{"(+ a expr)", nil},
	}
	for x, test := range tests {
		res := eval.EvalExpr(test.expr)
		if res != test.res {
			t.Log(x, "- Expected:", test.res)
			t.Fatal(x, "- Got:", res)
		}
	}
}
