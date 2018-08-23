package lex

import (
	"fmt"
	"strconv"

	"github.com/evovetech/lex/ast"
	"github.com/evovetech/lex/token"
)

type Parser struct {
	lex *Lexer
	cur Token
}

func NewParser(lex *Lexer) *Parser {
	p := &Parser{
		lex: lex,
		cur: lex.NextToken(),
	}
	return p
}

func (p *Parser) NextToken() {
	p.cur = p.lex.NextToken()
}

func (p *Parser) ParsePrimary() (ast.Expression, error) {
	switch tok := p.cur; tok.kind {
	case token.IDENT:
		return p.ParseIdentExpr()
	case token.NUMBER:
		return p.ParseNumberExpr()
	default:
		p.NextToken()
		return nil, fmt.Errorf("error for token: %s", tok)
	}
}

func (p *Parser) ParseIdentExpr() (ast.Expression, error) {
	var name = p.cur.val.RawString()
	var args []ast.Expression

	p.NextToken()
	tok := p.cur
	val := tok.val
	if val.Raw()[0] != '(' {
		return &ast.VariableExpr{Name: name}, nil
	}

	// eat '('
	p.NextToken()

	// TODO: p.ParseExpr()
	for p.cur.val.Raw()[0] != ')' {
		p.NextToken()
	}

	// eat ')'
	p.NextToken()

	// return
	call := &ast.CallExpr{
		Callee: name,
		Args:   args,
	}
	return call, nil
}

func (p *Parser) ParseNumberExpr() (num *ast.NumberExpr, err error) {
	tok := p.cur
	var val float64
	if val, err = strconv.ParseFloat(tok.Value().RawString(), 64); err == nil {
		num = &ast.NumberExpr{Val: val}
		p.NextToken()
	}
	return
}
