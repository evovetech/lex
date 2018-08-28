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

func (p *Parser) ParseNumberExpr() (num *ast.NumberExpr, err error) {
	tok := p.cur
	var val float64
	if val, err = strconv.ParseFloat(tok.Value().RawString(), 64); err == nil {
		num = &ast.NumberExpr{Val: val}
		p.NextToken()
	}
	return
}

func (p *Parser) ParseParenExpr() (expr ast.Expression, err error) {
	// eat '('
	p.NextToken()

	if expr, err = p.ParseExpression(); err != nil {
		return
	}

	if p.cur.kind == token.RPAREN {
		// eat ')'
		p.NextToken()
		return
	}

	return nil, fmt.Errorf("expected %s", token.RPAREN)
}

func (p *Parser) ParseIdentifierExpr() (ast.Expression, error) {
	var name = p.cur.val.RawString()
	var args []ast.Expression

	p.NextToken()
	if p.cur.kind != token.LPAREN {
		return &ast.VariableExpr{Name: name}, nil
	}

	// eat '('
	p.NextToken()
	if p.cur.kind == token.RPAREN {
		goto done
	}

	for {
		if arg, err := p.ParseExpression(); err != nil {
			return nil, err
		} else {
			args = append(args, arg)
		}
		if p.cur.kind == token.RPAREN {
			goto done
		}
		if p.cur.kind != token.COMMA {
			return nil, fmt.Errorf("error, expected comma but got %s", p.cur)
		}
		p.NextToken()
	}

done:

// eat ')'
	p.NextToken()

	// return
	call := &ast.CallExpr{
		Callee: name,
		Args:   args,
	}
	return call, nil
}

func (p *Parser) ParsePrimary() (ast.Expression, error) {
	switch tok := p.cur; tok.kind {
	case token.IDENTIFIER:
		return p.ParseIdentifierExpr()
	case token.NUMBER:
		return p.ParseNumberExpr()
	case token.LPAREN:
		return p.ParseParenExpr()
	default:
		p.NextToken()
		return nil, fmt.Errorf("error for token: %s", tok)
	}
}

func (p *Parser) ParseBinOpRhs(exprPrec int, lhs ast.Expression) (ast.Expression, error) {
	return nil, nil
}

func (p *Parser) ParseExpression() (ast.Expression, error) {
	_, err := p.ParsePrimary()
	if err != nil {
		return nil, err
	}

	// TODO:
	err = fmt.Errorf("TODO: parse expression '%s'", p.cur)
	p.NextToken()
	return &ast.ErrorExpr{Err: err}, nil
}
