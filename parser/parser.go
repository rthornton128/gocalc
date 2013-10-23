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

	p := new(Parser)
	p.init(f, str)
	//p.parse() // will become parseFile()
	root := ast.NewFile(token.Pos(1), token.Pos(len(str)+1))
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
		p.next()
	}
	//if p.file.NumErrors() > 0 {
	//	p.file.PrintErrors()
	//	return nil
	//}
	return root
}

type Parser struct {
	file *token.File
	scan *scanner.Scanner
	tok  token.Token
	pos  token.Pos
	lit  string
}

func (p *Parser) init(file *token.File, expr string) {
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

func (p *Parser) next() {
	p.tok, p.pos, p.lit = p.scan.Scan()
	p.pos += p.file.Base()
	//fmt.Println("tok:", p.tok)
	//fmt.Println("pos:", p.pos)
	//fmt.Println("lit:", p.lit)
}

func (p *Parser) parse() (ast.Node, *perror) {
	var n ast.Node = nil
	switch p.tok {
	case token.ADD, token.SUB, token.MUL, token.DIV, token.MOD, token.LT,
		token.LTE, token.GT, token.GTE, token.EQ, token.NEQ:
		n = &ast.Operator{p.pos, p.lit}
	case token.IDENT, token.IF, token.PRINT, token.SET: // temporary
		n = p.parseIdentifier()
	case token.NUMBER:
		n = p.parseNumber()
	case token.LPAREN:
		return p.parseExpression()
	case token.RPAREN:
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

func (p *Parser) parseExpression() (*ast.Expression, *perror) {
	// an lparen was found. scan until an rparen found. determine expression
	// type: define (arg list, body), if, print and set expressions or
	// a typical math expression (operator).
	open := true
	e := new(ast.Expression)
	e.LParen = p.pos
	offset := token.Pos(1)
	for {
		// here is where I could attack scope, rather than during the evaluation
		// stage. I will need to track the current (inner) scope in p (parser).
		// I could also make more intelligent decisions about what type of
		// expression I have rather than trying to resolve it runtime, too. This
		// could allow me to make better errors and optimizations later on
		p.next()
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
	}
	if open == true {
		return nil, &perror{p.pos, openError}
	}
	e.RParen = p.pos + offset
	return e, nil
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{p.pos, p.lit}
}

func (p *Parser) parseNumber() *ast.Number {
	i, err := strconv.ParseInt(p.lit, 0, 64)
	if err != nil {
		p.file.AddError(p.pos, "Parse:", err)
	}
	return &ast.Number{p.pos, p.lit, int(i)}
}
