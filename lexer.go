package lex

import (
	"io"

	"github.com/evovetech/lex/token"
)

type Lexer struct {
	scan token.Scanner
}

func NewLexer(rd io.Reader) *Lexer {
	return &Lexer{
		scan: token.NewScanner(rd),
	}
}

func (l *Lexer) NextToken() (tok Token) {
	tok.kind = l.scan.NextToken(&tok.val)
	return
}
