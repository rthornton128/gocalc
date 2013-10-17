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
	op_end

	lit_start
	IDENT
	NUMBER
	lit_end
)

type Pos int
