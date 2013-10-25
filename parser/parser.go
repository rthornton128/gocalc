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
	case token.IDENT, token.IF: // temporary
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
	offset := token.Pos(1)

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
		//e.Nodes = p.parseDefineExpression()
		//open = false
	case token.PRINT:
		e.Nodes = p.parsePrintExpression()
		open = false
	case token.SET:
		e.Nodes = p.parseSetExpression()
		open = false
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
			offset = n.End() // - n.Pos()
			p.next()
		}
	}
	if open == true {
		return nil, &perror{p.pos, openError}
	}
	e.RParen = p.pos + offset
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

func (p *parser) parsePrintExpression() []ast.Node {
	nodes := make([]ast.Node, 0)
	nodes = append(nodes, p.parseIdentifier()) // blah
	p.next()
	var n ast.Node
	for p.tok != token.RPAREN {
		switch p.tok {
		case token.LPAREN:
			var err *perror
			n, err = p.parseExpression()
			if err != nil {
				p.file.AddError(err.pos, err.msg)
			}
		case token.NUMBER:
			n = p.parseNumber()
		case token.IDENT:
			//if n = p.curScope.Lookup(p.lit); n == nil {
			//p.file.AddError(p.pos, "Undeclared identifier:", p.lit)
			//}
			n = p.parseIdentifier()
		default:
			p.file.AddError(p.pos, "Invalid argument to print: ", p.lit)
		}
		//p.curScope.Nodes = append(p.curScope.Nodes, n)
		nodes = append(nodes, n)
		p.next()
	}
	//if p.tok != token.RPAREN {
	//p.file.AddError(p.pos, "Unknown token:", p.lit, "Expected: ')'")
	//}
	return nodes
}

func (p *parser) parseSetExpression() []ast.Node {
	// eventually expand this for multiple assignment
	nodes := make([]ast.Node, 0)
	nodes = append(nodes, p.parseIdentifier()) // blah
	p.next()
	if p.tok != token.IDENT {
		p.file.AddError(p.pos, "First argument to set must be an identifier")
		return nil
	}
	i := p.parseIdentifier()
	nodes = append(nodes, i)
	p.next()
	var n ast.Node
	switch p.tok {
	case token.LPAREN:
		var err *perror
		n, err = p.parseExpression() // should really be parseMath?Expr()
		if err != nil {
			p.file.AddError(err.pos, err.msg)
		}
	case token.NUMBER:
		n = p.parseNumber()
	case token.IDENT:
		/*if p.curScope.Lookup(p.lit) == nil {
			p.file.AddError(p.pos, "Undeclared identifier:", p.lit)
		}*/
		n = p.parseIdentifier()
	}
	p.next()
	if p.tok != token.RPAREN {
		p.file.AddError(p.pos, "Unknown token:", p.lit, "Expected: ')'")
	}
	nodes = append(nodes, n)
	return nodes
}

func (p *parser) parseDefineExpression(lparen token.Pos) *ast.DefineExpr {
	d := new(ast.DefineExpr)
	d.LParen = lparen
	//d.Decl = make([]ast.Node, 1)
	d.Args = make([]string, 0)
	d.Impl = new(ast.Expression)
	//d.Decl = append(d.Decl, p.parseIdentifier()) // blah
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
	//d.Decl = append(d.Decl, n)
	//p.next()
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
