package parser_test

import (
	"github.com/rthornton128/gocalc/ast"
	"github.com/rthornton128/gocalc/parser"
	"testing"
)

func TestParserBasic(t *testing.T) {
	/* Number */
	n := parser.ParseExpr("123")
	if f, ok := n.(*ast.File); !ok {
		t.Log("File not received")
		t.FailNow()
	} else {
		if len(f.Nodes) != 1 {
			t.Log("expected 1 node, got:", len(f.Nodes))
			t.FailNow()
		}
		if i, ok := f.Nodes[0].(*ast.Number); !ok || i.Val != 123 {
			t.Log("Expected: true, 123")
			t.Fatal("Got:", ok, ",", i.Val)
		}
	}

	/* Identifier */
	n = parser.ParseExpr("abc")
	if f, ok := n.(*ast.File); !ok {
		t.Log("File not received")
		t.FailNow()
	} else {
		if len(f.Nodes) != 1 {
			t.Log("expected 1 node, got:", len(f.Nodes))
			t.FailNow()
		}
		if _, ok := f.Nodes[0].(*ast.Identifier); !ok {
			t.Log("Expected: true")
			t.Fatal("Got:", ok)
			t.FailNow()
		}
	}

	/* Operator */
	n = parser.ParseExpr("+")
	if f, ok := n.(*ast.File); !ok {
		t.Log("File not received")
		t.FailNow()
	} else {
		if len(f.Nodes) != 0 {
			t.Log("expected 0 nodes, got:", len(f.Nodes))
			t.FailNow()
		}
	}

	/* Expression */
	n = parser.ParseExpr("(+ 2 4)")
	if f, ok := n.(*ast.File); !ok {
		t.Log("File not received")
		t.FailNow()
	} else {
		if len(f.Nodes) != 1 {
			t.Log("expected 1 node, got:", len(f.Nodes))
			t.FailNow()
		}
		if _, ok := f.Nodes[0].(*ast.MathExpr); !ok {
			t.Log("Expected: true")
			t.Fatal("Got:", ok)
			t.FailNow()
		}
	}

	/* String */
	n = parser.ParseExpr("\"a string\"")
	if f, ok := n.(*ast.File); !ok {
		t.Log("File not received")
		t.FailNow()
	} else {
		if len(f.Nodes) != 1 {
			t.Log("expected 1 node, got:", len(f.Nodes))
			t.FailNow()
		}
		if i, ok := f.Nodes[0].(*ast.String); !ok || i.Lit != "\"a string\"" {
			t.Log("Expected: true, \"a string\"")
			t.Fatal("Got:", ok, ",", i.Lit)
			t.FailNow()
		}
	}
}
