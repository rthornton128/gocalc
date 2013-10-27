package parser

import (
	"errors"
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
	//p.topScope = root.Scope
	//p.curScope = root.Scope
	for n, err := p.parse(); ; n, err = p.parse() {
		if err != nil {
			p.file.AddError(err.pos, "Parse: ", err.msg)
		}
		if _, ok := n.(*ast.Operator); ok {
			p.file.AddError(n.Pos(), "Parse: Invalid expression, operator "+
				"may not be outside of expression")
		}
		if n == nil {
			break
		}
		root.Nodes = append(root.Nodes, n)
		//p.topScope.Nodes = append(p.topScope.Nodes, n)
		p.next()
	}
	//if p.file.NumErrors() > 0 {
	//	p.file.PrintErrors()
	//	return nil
	//}
	return root
}

type parser struct {
	file *token.File
	scan *scanner.Scanner
	//topScope *ast.Scope
	//curScope *ast.Scope
	tok token.Token
	pos token.Pos
	lit string
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

var closeError = errors.New("Unexpected ')'")
var eofError = errors.New("Reached end of file")
var openError = errors.New("Opening '(' with no closing bracket.")

func (p *parser) next() {
	p.tok, p.pos, p.lit = p.scan.Scan()
	p.pos += p.file.Base()
	//fmt.Println("tok:", p.tok)
	//fmt.Println("pos:", p.pos)
	//fmt.Println("lit:", p.lit)
}

func (p *parser) parse() (ast.Node, *perror) {
	var n ast.Node = nil
	switch p.tok {
	case token.ADD, token.SUB, token.MUL, token.DIV, token.MOD, token.LT,
		token.LTE, token.GT, token.GTE, token.EQ, token.NEQ:
		n = &ast.Operator{p.pos, p.lit}
	case token.IDENT:
		n = p.parseIdentifier()
	case token.NUMBER:
		n = p.parseNumber()
	case token.LPAREN:
		return p.parseExpression()
	case token.RPAREN:
		// this is gradually getting phased out, never fear!
		// I don't really like this solution...it feels clunky/messy. It's like
		// using exception handling. Receiving an RPAREN is not an error unless
		// it's out of place, so treating it like one no matter the situation
		// seems stupid.
		return nil, &perror{p.pos, closeError}
	case token.COMMENT:
		// consume comment and move on
		p.next()
		return p.parse()
	}
	return n, nil // eofError
}

func (p *parser) parseExpression() (ast.Node, *perror) {
	// an lparen was found. scan until an rparen found. determine expression
	// type: define (arg list, body), if, print and set expressions or
	// a typical math expression (operator).
	open := true
	e := new(ast.Expression)
	//e.Scope = ast.NewScope(p.curScope)
	//p.curScope = e.Scope
	e.LParen = p.pos
	//offset := token.Pos(1)

	p.next()
	switch p.tok {
	case token.LPAREN:
		p.file.AddError(p.pos, "Parse: First element of an expression may not "+
			"be another expression!")
		return nil, nil
	case token.RPAREN:
		p.file.AddError(p.pos, "Parse: Empty expression not allowed.")
		return nil, nil
	case token.DEFINE:
		return p.parseDefineExpression(e.LParen), nil
	case token.IF:
		return p.parseIfExpression(e.LParen), nil
	case token.PRINT:
		return p.parsePrintExpression(e.LParen), nil
	case token.SET:
		return p.parseSetExpression(e.LParen), nil
	default:
		// actually, I want to remove this section entirely...
		for {
			// here is where I could attack scope, rather than during the evaluation
			// stage. I will need to track the current (inner) scope in p (parser).
			// I could also make more intelligent decisions about what type of
			// expression I have rather than trying to resolve it runtime, too. This
			// could allow me to make better errors and optimizations later on

			n, err := p.parse()
			if err != nil {
				if err.msg == closeError {
					open = false
					break
				}
			}

			if n == nil {
				break
			}
			e.Nodes = append(e.Nodes, n)
			//offset = n.End() // - n.Pos()
			p.next()
		}
	}
	if open == true {
		return nil, &perror{p.pos, openError}
	}
	e.RParen = p.pos
	return e, nil
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

func (p *parser) parseNumber() *ast.Number {
	i, err := strconv.ParseInt(p.lit, 0, 64)
	if err != nil {
		p.file.AddError(p.pos, "Parse:", err)
	}
	return &ast.Number{p.pos, p.lit, int(i)}
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
		n, _ = p.parseExpression()
	case token.NUMBER:
		n = p.parseNumber()
	default:
		p.file.AddError(p.pos, "Unexpected token: ", p.lit)
		n = nil
	}
	p.next()
	return n
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

func (p *parser) parseDefineExpression(lparen token.Pos) *ast.DefineExpr {
	d := new(ast.DefineExpr)
	d.LParen = lparen
	d.Args = make([]string, 0)
	d.Impl = new(ast.Expression)
	p.next()
	var n ast.Node
	switch p.tok {
	case token.LPAREN:
		e := p.parseIdentifierList()
		l := e.Nodes
		d.Name = l[0].(*ast.Identifier).Lit
		l = l[1:]
		for _, v := range l {
			d.Args = append(d.Args, v.(*ast.Identifier).Lit)
		}
	case token.IDENT:
		d.Name = p.parseIdentifier().Lit
		p.next()
	default:
		p.file.AddError(p.pos, "Expected identifier(s) but got: ", p.lit)
		return nil
	}
	if p.tok != token.LPAREN {
		p.file.AddError(p.pos, "Expected expression but got: ", p.lit)
		return nil
	}
	var err *perror
	n, err = p.parseExpression()
	if err != nil {
		p.file.AddError(err.pos, err.msg)
		return nil
	}
	d.Impl = n
	p.next()
	if p.tok != token.RPAREN {
		p.file.AddError(p.pos, "Expected closing paren but got: ", p.lit)
		return nil
	}
	return d
}
