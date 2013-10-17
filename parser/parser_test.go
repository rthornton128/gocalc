package parser_test

import (
	"misc/calc/ast"
	"misc/calc/parser"
	"testing"
)

func TestParserBasic(t *testing.T) {
	n := parser.ParseExpr("123")
	if _, ok := n.(*ast.Number); !ok {
		t.FailNow()
	}
	n = parser.ParseExpr("abc")
	if _, ok := n.(*ast.Identifier); !ok {
		t.FailNow()
	}
	n = parser.ParseExpr("+")
	if _, ok := n.(*ast.Identifier); ok {
		t.FailNow()
	}
	n = parser.ParseExpr("(+ 2 4)")
	if _, ok := n.(*ast.Expression); !ok {
		t.FailNow()
	}
}
