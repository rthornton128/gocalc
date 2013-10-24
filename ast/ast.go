package ast

import (
	"misc/calc/token"
)

type (
	Node interface {
		Pos() token.Pos
		End() token.Pos
	}
	Identifier struct {
		Id  token.Pos
		Lit string
	}
	Number struct {
		Num token.Pos
		Lit string
		Val int
	}
	Operator struct {
		Opr token.Pos
		Val string
	}
	Expression struct {
		LParen token.Pos
		RParen token.Pos
		Nodes  []Node
		Scope  *Scope
	}
	File struct {
		pos   token.Pos
		end   token.Pos
		Nodes []Node
		Scope *Scope
	}
	PrintExpr struct {
		*Expression
	}
	Scope struct {
		defs map[string]Node
		//Nodes  []Node // temporary
		parent *Scope
	}
)

func (i *Identifier) Pos() token.Pos { return i.Id }
func (i *Identifier) End() token.Pos { return i.Id + token.Pos(len(i.Lit)) }

func (n *Number) Pos() token.Pos { return n.Num }
func (n *Number) End() token.Pos { return n.Num + token.Pos(len(n.Lit)) }

func (o *Operator) Pos() token.Pos { return o.Opr }
func (o *Operator) End() token.Pos { return o.Opr + 1 }

func (e *Expression) Pos() token.Pos { return e.LParen }
func (e *Expression) End() token.Pos { return e.RParen }

//func (p *PrintExpr) Pos() token.Pos { return p.Prnt }
//func (p *PrintExpr) End() token.Pos { return p.Prnt + 5 } // derp

func NewFile(beg, end token.Pos) *File {
	return &File{beg, end, make([]Node, 0), NewScope(nil)}
}

func (f *File) Pos() token.Pos { return f.pos }
func (f *File) End() token.Pos { return f.end }

func NewScope(parent *Scope) *Scope {
	return &Scope{make(map[string]Node), parent}
}

func (s *Scope) Insert(ident string, n Node) {
	s.defs[ident] = n
}

func (s *Scope) Lookup(ident string) Node {
	m := s
	for m != nil {
		if n, ok := m.defs[ident]; ok {
			return n
		}
		m = m.parent
	}
	return nil
}
