package scanner_test

import (
	"github.com/rthornton128/gocalc/scanner"
	"github.com/rthornton128/gocalc/token"
	"testing"
)

func TestScannerInit(t *testing.T) {
	var tests = []struct {
		expr string
		res  bool
	}{
		{"", true},
		{"(+ 1 2)", true},
	}
	for _, t := range tests {
		s := new(scanner.Scanner)
		f := token.NewFile("", t.expr)
		s.Init(f, t.expr)
	}
}

func TestScannerScan(t *testing.T) {
	var tests = []struct {
		expr string
		toks []token.Token
		poss []token.Pos
		lits []string
	}{
		{"123", []token.Token{token.NUMBER}, []token.Pos{0}, []string{"123"}},
		{"-123", []token.Token{token.NUMBER}, []token.Pos{0}, []string{"-123"}},
		{"a", []token.Token{token.IDENT}, []token.Pos{0}, []string{"a"}},
		{
			"123 456",
			[]token.Token{token.NUMBER, token.NUMBER},
			[]token.Pos{0, 4}, []string{"123", "456"},
		},
		{
			"(+123 -456)",
			[]token.Token{token.LPAREN, token.ADD, token.NUMBER, token.NUMBER,
				token.RPAREN},
			[]token.Pos{0, 1, 2, 6, 10}, []string{"(", "+", "123", "-456", ")"},
		},
		{
			"(set A_Variable 123)",
			[]token.Token{token.LPAREN, token.SET, token.IDENT, token.NUMBER,
				token.RPAREN},
			[]token.Pos{0, 1, 5, 16, 19}, []string{"(", "set", "A_Variable", "123",
				")"},
		},
		{
			"; A comment",
			[]token.Token{token.COMMENT},
			[]token.Pos{0},
			[]string{"; A comment"},
		},
		{
			"234;comment",
			[]token.Token{token.NUMBER, token.COMMENT},
			[]token.Pos{0, 3},
			[]string{"234", ";comment"},
		},
		{
			";comment\n234",
			[]token.Token{token.COMMENT, token.NUMBER},
			[]token.Pos{0, 9},
			[]string{";comment", "234"},
		},
		{
			"\"a string\"",
			[]token.Token{token.STRING},
			[]token.Pos{0},
			[]string{"\"a string\""},
		},
		{
			"\"a string\"\"a string\"",
			[]token.Token{token.STRING, token.STRING},
			[]token.Pos{0, 10},
			[]string{"\"a string\"", "\"a string\""},
		},
	}
	for x, test := range tests {
		s := new(scanner.Scanner)
		f := token.NewFile("", test.expr)
		s.Init(f, test.expr)
		for i := 0; i < len(test.toks); i++ {
			tok, pos, lit := s.Scan()
			if tok != test.toks[i] || pos != test.poss[i] || lit != test.lits[i] {
				t.Log("Test:", x)
				t.Log("Expected:", test.toks[i], test.poss[i], test.lits[i])
				t.Fatal("Got:", tok, pos, lit)
			}
		}
	}
}
