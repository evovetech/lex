package op

import "github.com/evovetech/lex/token"

type Precedence int32

const (
	Invalid Precedence = 10 * (iota - 1)
	NoOp
	LtGt
	AddSub
	_
	MultDiv
)

var ops = map[token.Token]Precedence{
	token.LT:    LtGt,
	token.PLUS:  AddSub,
	token.MINUS: AddSub,
	token.STAR:  MultDiv,
}

func GetPrecedence(token token.Token) Precedence {
	if op, ok := ops[token]; ok {
		return op
	}
	return Invalid
}
