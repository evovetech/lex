package lex

import (
	"fmt"
	"strconv"

	"github.com/evovetech/lex/ast"
	"github.com/evovetech/lex/op"
	"github.com/evovetech/lex/token"
)

type Parser struct {
	lex *Lexer
	cur *Token
}

func NewParser(lex *Lexer) *Parser {
	p := &Parser{
		lex: lex,
	}
	return p
}

func (p *Parser) init() {
	if p.cur == nil {
		tok := p.lex.NextToken()
		p.cur = &tok
	}
}

func (p *Parser) CurToken() Token {
	p.init()
	return *p.cur
}

func (p *Parser) NextToken() {
	p.init()
	*p.cur = p.lex.NextToken()
}

func (p *Parser) ParseNumberExpr() (num *ast.NumberExpr, err error) {
	tok := p.CurToken()
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

	if p.CurToken().kind == token.RPAREN {
		// eat ')'
		p.NextToken()
		return
	}

	return nil, fmt.Errorf("expected %s", token.RPAREN)
}

func (p *Parser) ParseIdentifierExpr() (ast.Expression, error) {
	var name = p.CurToken().val.RawString()
	var args []ast.Expression

	p.NextToken()
	if p.CurToken().kind != token.LPAREN {
		return &ast.VariableExpr{Name: name}, nil
	}

	// eat '('
	p.NextToken()
	if p.CurToken().kind == token.RPAREN {
		goto done
	}

	for {
		if arg, err := p.ParseExpression(); err != nil {
			return nil, err
		} else {
			args = append(args, arg)
		}
		cur := p.CurToken()
		if cur.kind == token.RPAREN {
			goto done
		}
		if cur.kind != token.COMMA {
			return nil, fmt.Errorf("error, expected comma but got %s", cur)
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
	switch tok := p.CurToken(); tok.kind {
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

func (p *Parser) ParseBinOpRhs(exprPrec op.Precedence, lhs ast.Expression) (ast.Expression, error) {
	expr := lhs
	for {
		// If this is a binop that binds at least as tightly as the current binop,
		// consume it, otherwise we are done.
		var tokPrec op.Precedence
		if tokPrec = p.CurToken().Precedence(); tokPrec < exprPrec {
			return lhs, nil
		}

		// Okay, we know this is a binop.
		binOp := p.CurToken().Value()
		p.NextToken() // eat binop

		// Parse the primary expression after the binary operator.
		var rhs ast.Expression
		var err error
		if rhs, err = p.ParsePrimary(); err != nil {
			return nil, err
		}

		// If BinOp binds less tightly with RHS than the operator after RHS, let
		// the pending operator take RHS as its LHS.
		if nextPrec := p.CurToken().Precedence(); tokPrec < nextPrec {
			if rhs, err = p.ParseBinOpRhs(tokPrec+1, rhs); err != nil {
				return nil, err
			}
		}

		// merge lhs/rhs
		lhs = &ast.BinaryExpr{
			Op:    binOp,
			Left:  lhs,
			Right: rhs,
		}
	}
	return nil, fmt.Errorf("error for %v %s", exprPrec, expr)
}

func (p *Parser) ParsePrototype() (*ast.PrototypeExpr, error) {
	tok := p.CurToken()
	if tok.kind != token.IDENTIFIER {
		return nil, fmt.Errorf("expected function name in prototype: got %s", tok)
	}

	name := tok.val.RawString()
	p.NextToken()

	if tok = p.CurToken(); tok.kind != token.LPAREN {
		return nil, fmt.Errorf("expected '(' in prototype: got %s", tok)
	}

	var argNames []string
	for {
		p.NextToken()
		if tok = p.CurToken(); tok.kind == token.IDENTIFIER {
			argNames = append(argNames, tok.val.RawString())
			continue
		}
		break
	}

	if tok = p.CurToken(); tok.kind != token.RPAREN {
		return nil, fmt.Errorf("expected ')' in prototype: got %s", tok)
	}
	p.NextToken() // eat '('

	proto := &ast.PrototypeExpr{
		Name: name,
		Args: argNames,
	}
	return proto, nil
}

func (p *Parser) ParseDefinition() (f *ast.FunctionExpr, err error) {
	p.NextToken() // eat def

	var proto *ast.PrototypeExpr
	if proto, err = p.ParsePrototype(); err != nil {
		return
	}

	var body ast.Expression
	if body, err = p.ParseExpression(); err != nil {
		return
	}

	f = &ast.FunctionExpr{
		Proto: proto,
		Body:  body,
	}
	return
}

func (p *Parser) ParseExtern() (*ast.PrototypeExpr, error) {
	p.NextToken() // eat extern
	return p.ParsePrototype()
}

func (p *Parser) ParseTopLevelExpression() (f *ast.FunctionExpr, err error) {
	var expr ast.Expression
	if expr, err = p.ParseExpression(); err == nil {
		f = &ast.FunctionExpr{
			// Make anonymous proto
			Proto: &ast.PrototypeExpr{
				Name: "__anon_expr",
				Args: []string{},
			},
			Body: expr,
		}
	}
	return
}

func (p *Parser) ParseExpression() (ast.Expression, error) {
	lhs, err := p.ParsePrimary()
	if err != nil {
		return nil, err
	}
	return p.ParseBinOpRhs(op.NOOP, lhs)
}
