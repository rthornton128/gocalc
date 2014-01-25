package parser_test

import (
	"github.com/rthornton128/gocalc/ast"
	"github.com/rthornton128/gocalc/parser"
	"testing"
)

func TestParserBasic(t *testing.T) {
	var tests = []struct {
		expr string
		res  int
	}{
		{"42", 0},
		{"\"string\"", 0},
		{"()", 0},
		{"(+)", 0},
		{"(a)", 0},
		{"(42)", 0},
		{"(\"string\")", 0},
		{"(+ 2)", 0},
		{"(+ 42 32)", 1},
	}
	for i, test := range tests {
		n := parser.ParseExpr(test.expr)
		_, ok := n.(*ast.File)
		if !ok {
			t.Log(i, ") File not received")
			t.FailNow()
		}
	}
}
