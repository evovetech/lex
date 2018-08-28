package ast

import (
	"fmt"
	"strings"

	"github.com/evovetech/lex/token"
)

// The base Node interface
type Node interface {
	String() string
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

type stmt struct{}

func (s *stmt) statementNode() {}

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode()
}

type expr struct{}

func (e *expr) expressionNode() {}

// TEMP
type ErrorExpr struct {
	expr
	Err error
}

func (e *ErrorExpr) String() string {
	return e.Err.Error()
}

type NumberExpr struct {
	expr
	Val float64
}

func (n *NumberExpr) String() string {
	return fmt.Sprintf("%.3f", n.Val)
}

type VariableExpr struct {
	expr
	Name string
}

func (v *VariableExpr) String() string {
	return v.Name
}

type BinaryExpr struct {
	expr
	Op          token.Token
	Left, Right Expression
}

func (b *BinaryExpr) String() string {
	op := string(append([]rune{}, rune(b.Op)))
	return fmt.Sprintf("(%s %s %s)", b.Left, op, b.Right)
}

type CallExpr struct {
	expr
	Callee string
	Args   []Expression
}

func (c *CallExpr) String() string {
	var exprs []string
	for _, arg := range c.Args {
		exprs = append(exprs, arg.String())
	}
	args := strings.Join(exprs, ", ")
	return fmt.Sprintf("%s(%s)", c.Callee, args)
}

type PrototypeExpr struct {
	expr
	Name string
	Args []string
}

func (p *PrototypeExpr) String() string {
	args := strings.Join(p.Args, ", ")
	return fmt.Sprintf("%s(%s)", p.Name, args)
}

type FunctionExpr struct {
	expr
	Prototype *PrototypeExpr
	Body      Expression
}

func (f *FunctionExpr) String() string {
	return fmt.Sprintf("%s %s", f.Prototype, f.Body)
}
