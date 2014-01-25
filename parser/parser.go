// Copyright (c) 2013, Rob Thornton
// All rights reserved.
// This software is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package parser

import (
	"fmt"
	"github.com/rthornton128/gocalc/ast"
	"github.com/rthornton128/gocalc/scanner"
	"github.com/rthornton128/gocalc/token"
	"strconv"
)

func ParseExpr(expr string) ast.Node {
	f := token.NewFile("", expr, 1)
	return ParseFile(f, expr)
}

func ParseFile(f *token.File, str string) *ast.File {
	if f.Size() != len(str) {
		fmt.Println("File size does not match string length.")
		return nil
	}

	root := ast.NewFile(token.Pos(1), token.Pos(len(str)+1))
	p := new(parser)
	p.init(f, str)
	p.topScope = root.Scope
	p.curScope = root.Scope
	n := p.parse()
	for p.file.NumErrors() < 10 && p.tok != token.EOF {
		root.Nodes = append(root.Nodes, n)
		p.next()
		n = p.parse()
	}
	if p.topScope != p.curScope {
		panic("Imbalanced scope!")
	}
	return root
}

type parser struct {
	file     *token.File
	scan     *scanner.Scanner
	topScope *ast.Scope
	curScope *ast.Scope
	tok      token.Token
	pos      token.Pos
	lit      string
}

func (p *parser) addError(args ...interface{}) {
	p.file.AddError(p.pos, args...)
}

func (p *parser) init(file *token.File, expr string) {
	p.file = file
	p.scan = new(scanner.Scanner)
	p.scan.Init(file, expr)
	p.next()
}

type perror struct {
	pos token.Pos
	msg error
}

func (p *parser) next() {
	p.tok, p.pos, p.lit = p.scan.Scan()
	p.pos += p.file.Base()
	//fmt.Println("tok:", p.tok)
	//fmt.Println("pos:", p.pos)
	//fmt.Println("lit:", p.lit)
}

func (p *parser) parse() ast.Node {
	switch p.tok {
	case token.LPAREN:
		return p.parseExpression()
	case token.COMMENT:
		// consume comment and move on
		p.next()
		return p.parse()
	case token.EOF:
		return nil
	default:
		str := "token"
		switch p.tok {
		case token.IDENT:
			str = "identifier"
		case token.NUMBER:
			str = "number"
		case token.STRING:
			str = "string"
		}
		p.addError("Unexpected ", str, " outside of expression: ", p.lit)
		return nil
	}
	return nil
}

func (p *parser) parseComparisonExpression(lp token.Pos) *ast.CompExpr {
	ce := new(ast.CompExpr)
	ce.Nodes = make([]ast.Node, 2)
	ce.LParen = lp
	ce.CompLit = p.lit
	p.next()
	ce.Nodes[0] = p.parseSubExpression()
	ce.Nodes[1] = p.parseSubExpression()
	if ce.Nodes[0] == nil || ce.Nodes[1] == nil {
		p.addError("Conditional must have at least two valid arguments")
	}
	//p.expect(token.RPAREN)
	if p.tok != token.RPAREN {
		p.addError("Expected ')', got:", p.lit)
		return nil
	}
	ce.RParen = p.pos
	return ce
}

func (p *parser) parseDefineExpression(lparen token.Pos) *ast.DefineExpr {
	d := new(ast.DefineExpr)
	d.LParen = lparen
	d.Args = make([]string, 0)
	tmp := p.curScope
	d.Scope = ast.NewScope(p.curScope)
	p.curScope = d.Scope
	d.Nodes = make([]ast.Node, 0)
	p.next()
	switch p.tok {
	case token.LPAREN:
		e := p.parseIdentifierList()
		l := e.Nodes
		d.Name = l[0].(*ast.Identifier).Lit
		l = l[1:]
		for _, v := range l {
			d.Args = append(d.Args, v.(*ast.Identifier).Lit)
			d.Scope.Insert(v.(*ast.Identifier).Lit, nil)
			p.curScope.Insert(v.(*ast.Identifier).Lit, d)
		}
	case token.IDENT:
		d.Name = p.parseIdentifier().Lit
		p.next()
	default:
		p.addError("Expected identifier(s) but got: ", p.lit)
		return nil
	}
	tmp.Insert(d.Name, d)
	for p.tok != token.RPAREN {
		switch p.tok {
		case token.COMMENT: // skip comments
		case token.LPAREN:
			d.Nodes = append(d.Nodes, p.parseExpression())
		case token.IDENT:
			d.Nodes = append(d.Nodes, p.parseIdentifier())
		default:
			p.addError("Expected expression or identifier but got: ", p.lit)
			break
		}
		p.next()
	}
	if len(d.Nodes) < 1 {
		p.addError("Expected list of expressions but got: ", p.lit)
		d = nil // don't exit without reverting scope
	}
	if p.tok != token.RPAREN {
		p.addError("Expected closing paren but got: ", p.lit)
		d = nil // don't exit without reverting scope
	}
	p.curScope = tmp
	return d
}

func (p *parser) parseExpression() ast.Node {
	lparen := p.pos
	p.next()
	switch p.tok {
	case token.LPAREN:
		p.addError("First element of an expression may not " +
			"be another expression!")
		return nil
	case token.RPAREN:
		p.addError("Empty expression not allowed.")
		return nil
	case token.LT, token.LTE, token.GT, token.GTE, token.EQ, token.NEQ:
		return p.parseComparisonExpression(lparen)
	case token.ADD, token.SUB, token.MUL, token.DIV, token.MOD, token.AND,
		token.OR:
		return p.parseMathExpression(lparen)
	case token.DEFINE:
		return p.parseDefineExpression(lparen)
	case token.IDENT:
		return p.parseUserExpression(lparen)
	case token.IF:
		return p.parseIfExpression(lparen)
	case token.PRINT:
		return p.parsePrintExpression(lparen)
	case token.SET:
		return p.parseSetExpression(lparen)
	}
	return nil
}

func (p *parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{p.pos, p.lit}
}

func (p *parser) parseIdentifierList() *ast.Expression {
	e := new(ast.Expression)
	e.LParen = p.pos
	e.Nodes = make([]ast.Node, 0)
	p.next()
	for p.tok == token.IDENT {
		e.Nodes = append(e.Nodes, p.parseIdentifier())
		p.next()
	}
	if p.tok != token.RPAREN {
		p.addError("Expected identifier or rparen, got: ", p.lit)
		return nil
	}
	e.RParen = p.pos
	p.next()
	return e
}

func (p *parser) parseIfExpression(lparen token.Pos) *ast.IfExpr {
	ie := new(ast.IfExpr)
	ie.Nodes = make([]ast.Node, 3)
	ie.LParen = lparen
	p.next()
	ie.Nodes[0] = p.parseSubExpression()
	ie.Nodes[1] = p.parseSubExpression()
	if p.tok == token.RPAREN {
		ie.Nodes[2] = nil
		ie.RParen = p.pos
		return ie
	}
	ie.Nodes[2] = p.parseSubExpression()
	if p.tok != token.RPAREN { // TODO: p.expect(token.RPAREN)
		p.addError("Expected closing paren, got: ", p.lit)
		return nil
	}
	ie.RParen = p.pos
	if ie.Nodes[0] == nil || ie.Nodes[1] == nil {
		p.addError("'if' expression must at least two elements")
		return nil
	}
	return ie
}

func (p *parser) parseImportExpression(lp token.Pos) *ast.ImportExpr {
	ie := new(ast.ImportExpr)
	ie.LParen = lp
	if p.tok != token.STRING { // p.expect(token.STRING)
		p.addError("Expected string, got:", p.lit)
		return nil
	}
	ie.Import = p.lit
	if p.tok != token.RPAREN { // p.expect(token.STRING)
		p.addError("Expected closing paren, got:", p.lit)
		return nil
	}
	ie.RParen = p.pos
	return ie
}

func (p *parser) parseMathExpression(lp token.Pos) *ast.MathExpr {
	me := new(ast.MathExpr)
	me.LParen = lp
	me.Nodes = make([]ast.Node, 0)
	me.OpLit = p.lit
	p.next()
	for p.tok != token.RPAREN && p.tok != token.EOF {
		me.Nodes = append(me.Nodes, p.parseSubExpression())
	}
	if len(me.Nodes) < 2 {
		p.addError("Math expressions must have at least 2 arguments")
		return nil
	}
	//me.ExprList = p.parseExpressionList()
	me.RParen = p.pos
	return me
}

func (p *parser) parseNumber() *ast.Number {
	i, err := strconv.ParseInt(p.lit, 0, 64)
	if err != nil {
		p.addError(err)
	}
	return &ast.Number{p.pos, p.lit, int(i)}
}

func (p *parser) parsePrintExpression(lparen token.Pos) *ast.PrintExpr {
	pe := new(ast.PrintExpr)
	pe.LParen = lparen
	pe.Nodes = make([]ast.Node, 0)
	p.next()
	for p.tok != token.RPAREN {
		pe.Nodes = append(pe.Nodes, p.parseSubExpression2())
	}
	if p.tok != token.RPAREN {
		p.addError("Unknown token:", p.lit, "Expected: ')'")
	}
	pe.RParen = p.pos
	return pe
}

func (p *parser) parseSetExpression(lparen token.Pos) *ast.SetExpr {
	se := new(ast.SetExpr)
	se.LParen = lparen
	// TODO: eventually expand this for multiple assignment?
	p.next()
	if p.tok != token.IDENT {
		p.addError("First argument to set must be an identifier")
		return nil
	}
	se.Name = p.parseIdentifier().Lit
	p.next()
	se.Value = p.parseSubExpression()
	if p.tok != token.RPAREN {
		p.addError("Unknown token:", p.lit, "Expected: ')'")
	}
	se.RParen = p.pos
	p.curScope.Insert(se.Name, se)
	return se
}

func (p *parser) parseString() *ast.String {
	return &ast.String{p.pos, p.lit}
}

func (p *parser) parseSubExpression() ast.Node {
	for p.tok == token.COMMENT {
		p.next()
	}
	var n ast.Node
	switch p.tok {
	case token.IDENT:
		i := p.parseIdentifier()
		if p.curScope.Lookup(i.Lit) == nil {
			p.addError("Undeclared identifier: ", i.Lit)
			p.next()
			return nil
		}
		n = i
	case token.LPAREN:
		n = p.parseExpression()
	case token.NUMBER:
		n = p.parseNumber()
	case token.STRING:
		p.addError("Expected Number or Expression, got String:",
			p.lit)
	default:
		p.addError("Unexpected token: ", p.lit)
	}
	p.next()
	return n
}

func (p *parser) parseSubExpression2() ast.Node {
	for p.tok == token.COMMENT {
		p.next()
	}
	if p.tok == token.STRING {
		n := p.parseString()
		p.next()
		return n
	}
	return p.parseSubExpression()
}

func (p *parser) parseUserExpression(lp token.Pos) *ast.UserExpr {
	ident := p.curScope.Lookup(p.lit)
	if ident == nil {
		p.addError("Undeclared identifier: ", p.lit)
		return nil
	}
	de, ok := ident.(*ast.DefineExpr)
	if !ok {
		p.addError("Undeclared function: ", p.lit)
		return nil
	}
	ue := new(ast.UserExpr)
	ue.Name = p.lit
	p.next()
	for p.tok != token.RPAREN {
		e := p.parseSubExpression2()
		if e != nil {
			ue.Nodes = append(ue.Nodes, e)
		}
	}
	if len(ue.Nodes) != len(de.Args) {
		p.addError("Parameter count mismatch. Function takes ",
			len(de.Args), " parameters, got:", len(ue.Nodes))
		return nil
	}
	ue.RParen = p.pos
	return ue
}
