package op

import "github.com/evovetech/lex/token"

type Precedence int32

const (
	INVALID Precedence = 10 * (iota - 1)
	NOOP
	LT
	PLUS
	MINUS
	MULT
)

var ops = map[token.Token]Precedence{
	token.LT:    LT,
	token.PLUS:  PLUS,
	token.MINUS: MINUS,
	token.STAR:  MULT,
}

func GetPrecedence(token token.Token) Precedence {
	if op, ok := ops[token]; ok {
		return op
	}
	return INVALID
}
