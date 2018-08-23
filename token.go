package lex

import (
	"fmt"

	"github.com/evovetech/lex/token"
)

type Token struct {
	kind token.Token
	val  token.Value
}

func (t *Token) Kind() token.Token {
	return t.kind
}

func (t *Token) Value() token.Value {
	return t.val
}

func (t *Token) IsDone() bool {
	return t.val.Error() != nil || t.kind.IsDone()
}

func (t Token) String() string {
	return fmt.Sprintf("<%s>%s", t.kind, t.val)
}

type TokenBuf struct {
	tok  []Token
	r, w int32
}
