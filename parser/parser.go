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
	p.Init(f, str)
	root := ast.NewFile(token.Pos(1), token.Pos(len(str)+1))
	for n, err := p.next(); ; n, err = p.next() {
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
}

func (p *Parser) Init(file *token.File, expr string) {
	p.file = file
	p.scan = new(scanner.Scanner)
	p.scan.Init(file, expr)
}

type perror struct {
	pos token.Pos
	msg error
}

var closeError = errors.New("Unexpected ')'")
var eofError = errors.New("Reached end of file")
var openError = errors.New("Opening '(' with no closing bracket.")

func (p *Parser) next() (ast.Node, *perror) {
	tok, off, lit := p.scan.Scan()
	pos := p.file.Base() + off
	//fmt.Println("tok:", tok)
	//fmt.Println("pos:", pos)
	//fmt.Println("lit:", lit)
	open := false
	switch tok {
	case token.ADD, token.SUB, token.MUL, token.DIV, token.MOD:
		return &ast.Operator{pos, lit[0]}, nil
	case token.IDENT:
		//fmt.Println("Found identifier", lit, "at pos:", pos)
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
			return nil, &perror{pos, openError}
		}
		e.RParen = pos + offset
		return e, nil
	case token.RPAREN:
		if open == false {
			return nil, &perror{pos, closeError}
		}
	case token.COMMENT:
		//fmt.Println("Found comment:", lit)
		return p.next()
	}
	return nil, nil // eofError
}
