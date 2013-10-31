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
	lit_end

	key_start
	AND
	DEFINE
	IF
	OR
	PRINT
	SET
	key_end
)

var tokens = map[string]Token{
	"and":    AND,
	"define": DEFINE,
	"if":     IF,
	"or":     OR,
	"print":  PRINT,
	"set":    SET,
}

func Lookup(ident string) Token {
	if t, ok := tokens[ident]; ok {
		return t
	}
	return IDENT
}
