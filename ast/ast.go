package ast

import "github.com/evovetech/lex"

// The base Node interface
type Node interface {
	Token() lex.Token
	String() string
}

type node struct {
	tok lex.Token
}

func (n *node) Token() lex.Token {
	return n.tok
}

func (n *node) String() string {
	return n.tok.String()
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

type stmt struct {
	node
}

func (s *stmt) statementNode() {}

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode()
}

type expr struct {
	node
}

func (e *expr) expressionNode() {}

type NumberExpr struct {
	expr
	Val float64
}

type VariableExpr struct {
	expr
	Name string
}

type BinaryExpr struct {
	expr
	Op          uint8
	Left, Right Expression
}

type CallExpr struct {
	expr
	Callee string
	Args   []Expression
}
