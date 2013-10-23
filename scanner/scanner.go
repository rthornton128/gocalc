package scanner

import (
	//"fmt"
	"misc/calc/token"
)

type Scanner struct {
	ch        byte
	str       string
	off, roff int
	file      *token.File
}

func (s *Scanner) Init(file *token.File, expr string) {
	s.str = expr
	s.file = file
	s.off = 0
	s.roff = 0
	s.next()
}

func (s *Scanner) Scan() (tok token.Token, pos token.Pos, lit string) {
	s.skipWhitespace()

	if isAlpha(s.ch) {
		pos, lit = s.scanIdentifier()
		tok = token.Lookup(lit)
		return
	}
	if isDigit(s.ch) {
		pos, lit = s.scanNumber()
		return token.NUMBER, pos, lit
	}
	ch := s.ch
	s.next()
	switch ch {
	case '+':
		tok = token.ADD
	case '-':
		if isDigit(s.ch) { // is '-' a unary operator for a number?
			pos, lit = s.scanNumber()
			pos, lit, tok = pos-1, string('-')+lit, token.NUMBER
			return
		}
		tok = token.SUB
	case '*':
		tok = token.MUL
	case '/':
		tok = token.DIV
	case '%':
		tok = token.MOD
	case '(':
		tok = token.LPAREN
	case ')':
		tok = token.RPAREN
	case ';':
		pos, lit = s.scanComment()
		tok = token.COMMENT
		return
	case '=':
		tok = token.EQ
	case '<':
		switch s.ch {
		case '=':
			lit = string(ch) + string(s.ch)
			pos, tok = token.Pos(s.off-1), token.LTE
			s.next()
			return
		case '>':
			lit = string(ch) + string(s.ch)
			pos, tok = token.Pos(s.off-1), token.NEQ
			s.next()
			return
		}
		tok = token.LT
	case '>':
		if s.ch == '=' {
			lit = string(ch) + string(s.ch)
			pos, tok = token.Pos(s.off-1), token.GTE
			s.next()
			return
		}
		tok = token.GT
	default:
		if s.off >= len(s.str) {
			tok = token.EOF
			pos = token.Pos(s.off)
			return
		}
	}
	pos = token.Pos(s.off - 1)
	lit = string(ch)
	return
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (s *Scanner) next() {
	s.off = s.roff
	if s.roff < len(s.str) {
		s.ch = s.str[s.off]
	} else {
		s.ch = 0
	}
	if s.ch == '\n' {
		s.file.AddLine(s.off)
	}
	s.roff++
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func (s *Scanner) scanComment() (token.Pos, string) {
	start := s.off - 1
	for s.ch != '\n' && s.off < len(s.str) {
		s.next()
	}
	return token.Pos(start), s.str[start:s.off]
}

func (s *Scanner) scanIdentifier() (token.Pos, string) {
	start := s.off
	for isDigit(s.ch) || isAlpha(s.ch) || s.ch == '_' || s.ch == '-' {
		s.next()
	}
	return token.Pos(start), s.str[start:s.off]
}

func (s *Scanner) scanNumber() (token.Pos, string) {
	start := s.off
	for isDigit(s.ch) {
		s.next()
	}
	return token.Pos(start), s.str[start:s.off]
}
