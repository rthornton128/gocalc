package ast

import (
	"misc/calc/token"
)

type (
	Node interface {
		BegPos() token.Pos
		EndPos() token.Pos
	}
	Identifier struct {
		Pos token.Pos
		Lit string
	}
	Number struct {
		Pos token.Pos
		Lit string
		Val int
	}
	Operator struct {
		Pos token.Pos
		Val byte
	}
	Expression struct {
		LParen token.Pos
		RParen token.Pos
		Nodes  []Node
	}
	File struct {
		pos   token.Pos
		end   token.Pos
		Nodes []Node
		Scope *Scope
	}
	Scope struct {
		defs   map[string]Node
		parent *Scope
	}
)

func (i *Identifier) BegPos() token.Pos { return i.Pos }
func (i *Identifier) EndPos() token.Pos { return i.Pos + token.Pos(len(i.Lit)) }

func (n *Number) BegPos() token.Pos { return n.Pos }
func (n *Number) EndPos() token.Pos { return n.Pos + token.Pos(len(n.Lit)) }

func (o *Operator) BegPos() token.Pos { return o.Pos }
func (o *Operator) EndPos() token.Pos { return o.Pos + 1 }

func (e *Expression) BegPos() token.Pos { return e.LParen }
func (e *Expression) EndPos() token.Pos { return e.RParen }

func (f *File) BegPos() token.Pos { return f.pos }
func (f *File) EndPos() token.Pos { return f.end }

func NewScope(parent *Scope) *Scope {
	return &Scope{make(map[string]Node, 0), parent}
}

func (s *Scope) Insert(ident string, n Node) {
}

func (s *Scope) Lookup() Node { return nil }
