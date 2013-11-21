// Copyright (c) 2013, Rob Thornton
// All rights reserved.
// This software is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ast

import (
	"github.com/rthornton128/gocalc/token"
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
	String struct {
		Str token.Pos
		Lit string
	}
	Operator struct {
		Opr token.Pos
		Val string
	}
	Expression struct {
		LParen token.Pos
		RParen token.Pos
		Nodes  []Node
	}
	CompExpr struct {
		Expression
		CompLit string
	}
	DefineExpr struct {
		Expression
		Scope *Scope
		Name  string
		Args  []string // TODO: remove?
		Impl  []Node
	}
	IfExpr struct {
		Expression
	}
	ImportExpr struct {
		Expression
	}
	MathExpr struct {
		Expression
		OpLit    string
		ExprList []Node
	}
	PrintExpr struct {
		Expression
		Nodes []Node
	}
	SetExpr struct {
		Expression
		Name  string
		Value Node
	}
	UserExpr struct {
		Expression
		Name  string
		Nodes []Node
	}
	File struct {
		pos   token.Pos
		end   token.Pos
		Nodes []Node
		Scope *Scope
	}
	Scope struct {
		defs   map[string]interface{}
		Parent *Scope
	}
)

func (i *Identifier) Pos() token.Pos { return i.Id }
func (i *Identifier) End() token.Pos { return i.Id + token.Pos(len(i.Lit)) }

func (n *Number) Pos() token.Pos { return n.Num }
func (n *Number) End() token.Pos { return n.Num + token.Pos(len(n.Lit)) }

func (s *String) Pos() token.Pos { return s.Str }
func (s *String) End() token.Pos { return s.Str + token.Pos(len(s.Lit)) }

func (o *Operator) Pos() token.Pos { return o.Opr }
func (o *Operator) End() token.Pos { return o.Opr + 1 }

func (e *Expression) Pos() token.Pos { return e.LParen }
func (e *Expression) End() token.Pos { return e.RParen }

func NewFile(beg, end token.Pos) *File {
	return &File{beg, end, make([]Node, 0), NewScope(nil)}
}

func (f *File) Pos() token.Pos { return f.pos }
func (f *File) End() token.Pos { return f.end }

func NewScope(parent *Scope) *Scope {
	return &Scope{make(map[string]interface{}), parent}
}

func (s *Scope) Insert(ident string, n interface{}) {
	s.defs[ident] = n
}

func (s *Scope) Lookup(ident string) interface{} {
	m := s
	for m != nil {
		if n, ok := m.defs[ident]; ok {
			return n
		}
		m = m.Parent
	}
	return nil
}

func (s *Scope) String() string {
	tmp := s
	var str string
	for k, _ := range tmp.defs {
		str += k + " "
	}
	return "[ " + str + "]"
}
