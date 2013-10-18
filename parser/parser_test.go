package parser_test

import (
	"misc/calc/ast"
	"misc/calc/parser"
	"testing"
)

func TestParserBasic(t *testing.T) {
	/* Number */
	n := parser.ParseExpr("123")
	if f, ok := n.(*ast.File); !ok {
		t.FailNow()
	} else {
		if len(f.Nodes) != 1 {
			t.FailNow()
		}
		if i, ok := f.Nodes[0].(*ast.Number); !ok || i.Val != 123 {
			t.FailNow()
		}
	}

	/* Identifier */
	n = parser.ParseExpr("abc")
	if f, ok := n.(*ast.File); !ok {
		t.FailNow()
	} else {
		if len(f.Nodes) != 1 {
			t.FailNow()
		}
		if _, ok := f.Nodes[0].(*ast.Identifier); !ok {
			t.FailNow()
		}
	}

	/* Operator */
	n = parser.ParseExpr("+")
	if f, ok := n.(*ast.File); !ok {
		t.FailNow()
	} else {
		if len(f.Nodes) != 1 {
			t.FailNow()
		}
		if _, ok := f.Nodes[0].(*ast.Operator); !ok {
			t.FailNow()
		}
	}

	/* Expression */
	n = parser.ParseExpr("(+ 2 4)")
	if f, ok := n.(*ast.File); !ok {
		t.FailNow()
	} else {
		if len(f.Nodes) != 1 {
			t.FailNow()
		}
		if _, ok := f.Nodes[0].(*ast.Expression); !ok {
			t.FailNow()
		}
	}
}
