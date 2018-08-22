package lex

import (
	"fmt"
	"io"

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

type Lexer struct {
	scan token.Scanner
	tok  TokenBuf
}

func NewLexer(rd io.Reader) *Lexer {
	return &Lexer{
		scan: token.NewScanner(rd),
		tok: TokenBuf{
			tok: make([]Token, 1),
		},
	}
}

func (l *Lexer) NextToken() Token {
	tok := l.tok.tok[0]
	tok.kind = l.scan.NextToken(&tok.val)
	return tok
}
