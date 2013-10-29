package parser

import (
	//	"errors"
	"fmt"
	"misc/calc/ast"
	"misc/calc/scanner"
	"misc/calc/token"
	"strconv"
)

func ParseExpr(expr string) ast.Node {
	f := token.NewFile("", expr)
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
	for n := p.parse(); n != nil; n = p.parse() {
		root.Nodes = append(root.Nodes, n)
		p.next()
	}
	//if p.file.NumErrors() > 0 {
	//	p.file.PrintErrors()
	//	return nil
	//}
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

/*
var closeError = errors.New("Unexpected ')'")
var eofError = errors.New("Reached end of file")
var openError = errors.New("Opening '(' with no closing bracket.")
*/
func (p *parser) next() {
	p.tok, p.pos, p.lit = p.scan.Scan()
	p.pos += p.file.Base()
	//fmt.Println("tok:", p.tok)
	//fmt.Println("pos:", p.pos)
	//fmt.Println("lit:", p.lit)
}

func (p *parser) parse() ast.Node {
	var n ast.Node = nil
	switch p.tok {
	case token.IDENT:
		n = p.parseIdentifier()
	case token.NUMBER:
		n = p.parseNumber()
	case token.LPAREN:
		n = p.parseExpression()
	case token.COMMENT:
		// consume comment and move on
		p.next()
		return p.parse()
	case token.EOF:
		return nil
	default:
		p.file.AddError(p.pos, "Unexpected token outside of expression: ", p.lit)
		return nil
	}
	return n
}

func (p *parser) parseComparisonExpression(lp token.Pos) *ast.CompExpr {
	ce := new(ast.CompExpr)
	ce.LParen = lp
	ce.CompLit = p.lit
	p.next()
	ce.A = p.parseSubExpression()
	ce.B = p.parseSubExpression()
	if ce.A == nil || ce.B == nil { // doesn't seem right...
		p.file.AddError(p.pos, "Some kind of conditional error")
	}
	//p.expect(token.RPAREN)
	if p.tok != token.RPAREN {
		p.file.AddError(p.pos, "Expected ')', got:", p.lit)
		return nil
	}
	ce.RParen = p.pos
	return ce
}

func (p *parser) parseDefineExpression(lparen token.Pos) *ast.DefineExpr {
	d := new(ast.DefineExpr)
	d.LParen = lparen
	d.Args = make([]string, 0) // TODO: remove?
	//fmt.Println("up scope!")
	tmp := p.curScope
	d.Scope = ast.NewScope(p.curScope)
	p.curScope = d.Scope
	d.Impl = make([]ast.Node, 0)
	p.next()
	switch p.tok {
	case token.LPAREN:
		e := p.parseIdentifierList()
		l := e.Nodes
		d.Name = l[0].(*ast.Identifier).Lit
		l = l[1:]
		for _, v := range l {
			d.Args = append(d.Args, v.(*ast.Identifier).Lit) //TODO: remove?
			d.Scope.Insert(v.(*ast.Identifier).Lit, nil)
		}
	case token.IDENT:
		d.Name = p.parseIdentifier().Lit
		p.next()
	default:
		p.file.AddError(p.pos, "Expected identifier(s) but got: ", p.lit)
		return nil
	}
	//fmt.Println("parseDefine:", d.Name)
	for p.tok != token.RPAREN {
		if p.tok != token.LPAREN {
			p.file.AddError(p.pos, "Expected expression but got: ", p.lit)
			return nil
		}
		d.Impl = append(d.Impl, p.parseExpression())
		p.next()
	}
	if len(d.Impl) < 1 {
		p.file.AddError(p.pos, "Expected list of expressions but got: ", p.lit)
		return nil
	}
	if p.tok != token.RPAREN {
		p.file.AddError(p.pos, "Expected closing paren but got: ", p.lit)
		return nil
	}
	//fmt.Println("down scope!")
	//fmt.Println(d.Name, "had", len(d.Impl), "expressions as arguments")
	p.curScope = tmp
	return d
}

func (p *parser) parseExpression() ast.Node {
	lparen := p.pos
	p.next()
	switch p.tok {
	case token.LPAREN:
		p.file.AddError(p.pos, "Parse: First element of an expression may not "+
			"be another expression!")
		return nil
	case token.RPAREN:
		p.file.AddError(p.pos, "Parse: Empty expression not allowed.")
		return nil
	case token.LT, token.LTE, token.GT, token.GTE, token.EQ, token.NEQ:
		return p.parseComparisonExpression(lparen)
	case token.ADD, token.SUB, token.MUL, token.DIV, token.MOD:
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
		p.file.AddError(p.pos, "Expected identifier or rparen, got: ", p.lit)
		return nil
	}
	e.RParen = p.pos
	p.next()
	return e
}

func (p *parser) parseIfExpression(lparen token.Pos) *ast.IfExpr {
	ie := new(ast.IfExpr)
	ie.LParen, ie.Else = lparen, nil
	p.next()
	ie.Comp = p.parseSubExpression()
	ie.Then = p.parseSubExpression()
	if p.tok == token.RPAREN {
		ie.Else = nil
		ie.RParen = p.pos
		return ie
	}
	ie.Else = p.parseSubExpression()
	if p.tok != token.RPAREN {
		p.file.AddError(p.pos, "Expected closing paren, got: ", p.lit)
		return nil
	}
	ie.RParen = p.pos
	if ie.Comp == nil || ie.Then == nil {
		return nil
	}
	return ie
}

func (p *parser) parseMathExpression(lp token.Pos) *ast.MathExpr {
	me := new(ast.MathExpr)
	me.OpLit = p.lit
	p.next()
	for p.tok != token.RPAREN && p.tok != token.EOF {
		me.ExprList = append(me.ExprList, p.parseSubExpression())
	}
	if len(me.ExprList) < 2 {
		p.file.AddError(p.pos, "Math expressions must have at least 2 arguments")
		return nil
	}
	//me.ExprList = p.parseExpressionList()
	me.RParen = p.pos
	return me
}

func (p *parser) parseNumber() *ast.Number {
	i, err := strconv.ParseInt(p.lit, 0, 64)
	if err != nil {
		p.file.AddError(p.pos, "Parse:", err)
	}
	return &ast.Number{p.pos, p.lit, int(i)}
}

func (p *parser) parsePrintExpression(lparen token.Pos) *ast.PrintExpr {
	pe := new(ast.PrintExpr)
	pe.LParen = lparen
	pe.Nodes = make([]ast.Node, 0)
	p.next()
	for p.tok != token.RPAREN {
		pe.Nodes = append(pe.Nodes, p.parseSubExpression())
	}
	if p.tok != token.RPAREN {
		p.file.AddError(p.pos, "Unknown token:", p.lit, "Expected: ')'")
	}
	pe.RParen = p.pos
	return pe
}

func (p *parser) parseSetExpression(lparen token.Pos) *ast.SetExpr {
	se := new(ast.SetExpr)
	se.LParen = lparen
	// eventually expand this for multiple assignment
	p.next()
	if p.tok != token.IDENT {
		p.file.AddError(p.pos, "First argument to set must be an identifier")
		return nil
	}
	se.Name = p.parseIdentifier().Lit
	p.next()
	se.Value = p.parseSubExpression()
	if p.tok != token.RPAREN {
		p.file.AddError(p.pos, "Unknown token:", p.lit, "Expected: ')'")
	}
	se.RParen = p.pos
	return se
}

func (p *parser) parseSubExpression() ast.Node {
	for p.tok == token.COMMENT {
		p.next()
	}
	var n ast.Node
	switch p.tok {
	case token.IDENT:
		n = p.parseIdentifier()
	case token.LPAREN:
		n = p.parseExpression()
	case token.NUMBER:
		n = p.parseNumber()
	default:
		//fmt.Println("subexpr, bad lit:", p.pos, "-", p.lit)
		p.file.AddError(p.pos, "Unexpected token: ", p.lit)
	}
	p.next()
	return n
}

func (p *parser) parseUserExpression(lp token.Pos) *ast.UserExpr {
	ue := new(ast.UserExpr)
	ue.Name = p.lit
	p.next()
	for p.tok != token.RPAREN {
		e := p.parseSubExpression()
		if e != nil {
			ue.Nodes = append(ue.Nodes, e)
		}
	}
	ue.RParen = p.pos
	return ue
}
