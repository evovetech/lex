package ast

import "fmt"

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
	Op          uint8
	Left, Right Expression
}

type CallExpr struct {
	expr
	Callee string
	Args   []Expression
}

func (c *CallExpr) String() string {
	first := true
	var args string
	for _, arg := range c.Args {
		if !first {
			args += ", "
		}
		first = false
		args += arg.String()
	}
	return fmt.Sprintf("%s(%s)", c.Callee, args)
}
