// Copyright (c) 2013, Rob Thornton
// All rights reserved.
// This software is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package token

type Token int

const (
	EOF Token = iota
	COMMENT

	op_start
	ADD
	SUB
	MUL
	DIV
	MOD
	LPAREN
	RPAREN
	EQ
	GT
	GTE
	LT
	LTE
	NEQ
	op_end

	lit_start
	IDENT
	NUMBER
	STRING
	lit_end

	key_start
	AND
	CASE
	DEFINE
	IF
	IMPORT
	OR
	PRINT
	SET
	SWITCH
	key_end
)

var tokens = map[string]Token{
	"and":    AND,
	"case":   CASE,
	"define": DEFINE,
	"if":     IF,
	"import": IMPORT,
	"or":     OR,
	"print":  PRINT,
	"set":    SET,
	"switch": SWITCH,
}

func Lookup(ident string) Token {
	if t, ok := tokens[ident]; ok {
		return t
	}
	return IDENT
}
