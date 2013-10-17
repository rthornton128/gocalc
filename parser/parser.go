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

func ParseFile(f *token.File, str string) ast.Node {
	if f.Size() != len(str) {
		fmt.Println("File size does not match string length.")
		return nil
	}
	p := new(Parser)
	p.Init(f, str)
	// This section is not ideal and should be handled better. Basically,
	// anything after the root node is an error. After
	// root is retrieved, I need to determine if anything follows and what kind
	// of error to generate.
	root, err := p.next()
	for {
		n, err := p.next()
		if err != nil {
			p.file.AddError(root.EndPos(), "Parse: ", err)
		}
		if n == nil {
			break
		}
	}
	if err != nil {
		if root == nil {
			p.file.AddError(token.Pos(1), "Parse: ", err)
		} else {
			p.file.AddError(root.EndPos(), "Parse: ", err)
		}
	}
	if _, ok := root.(*ast.Operator); ok {
		if root == nil {
			p.file.AddError(token.Pos(1), "Parse: Invalid expression")
		} else {
			p.file.AddError(root.EndPos(), "Parse: Invalid expression")
		}
	}

	if p.file.NumErrors() > 0 {
		p.file.PrintErrors()
		return nil
	}
	return root
}

type Parser struct {
	file *token.File
	scan *scanner.Scanner
}

func (p *Parser) Init(file *token.File, expr string) {
	p.file = file
	p.scan = new(scanner.Scanner)
	p.scan.Init(file, expr)
}

var closeError = errors.New("Unexpected ')'")
var eofError = errors.New("Reached end of file")
var openError = errors.New("Opening '(' with no closing bracket.")

func (p *Parser) next() (ast.Node, error) {
	tok, pos, lit := p.scan.Scan()
	open := false
	switch tok {
	case token.ADD, token.SUB, token.MUL, token.DIV, token.MOD:
		return &ast.Operator{pos, lit[0]}, nil
	case token.IDENT:
		return &ast.Identifier{pos, lit}, nil
	case token.NUMBER:
		i, err := strconv.ParseInt(lit, 0, 64)
		if err != nil {
			p.file.AddError(pos, "Parse:", err)
		}
		return &ast.Number{pos, lit, int(i)}, nil
	case token.LPAREN:
		open = true
		e := new(ast.Expression)
		e.LParen = pos
		offset := token.Pos(1)
		for {
			n, err := p.next()
			if err == closeError {
				open = false
				break
			}
			if n == nil {
				break
			}
			e.Nodes = append(e.Nodes, n)
			offset = n.EndPos() // - n.BegPos()
		}
		if open == true {
			return nil, openError
		}
		e.RParen = pos + offset
		return e, nil
	case token.RPAREN:
		if open == false {
			return nil, closeError
		}
	}
	return nil, nil // eofError
}
