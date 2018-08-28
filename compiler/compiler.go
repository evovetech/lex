package compiler

import (
	"fmt"

	"github.com/evovetech/lex/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

type Compiler interface {
	GetContext() llvm.Context
	Compile(node ast.Node) (llvm.Value, error)
}

func NewCompiler() Compiler {
	ctx := llvm.GlobalContext()
	return &compiler{
		Context:     ctx,
		builder:     ctx.NewBuilder(),
		namedValues: make(map[string]llvm.Value),
	}
}

type compiler struct {
	llvm.Context
	builder     llvm.Builder
	module      llvm.Module
	namedValues map[string]llvm.Value
}

func (c *compiler) GetContext() llvm.Context {
	return c.Context
}

func (c *compiler) Compile(node ast.Node) (val llvm.Value, err error) {
	switch e := node.(type) {
	case *ast.NumberExpr:
		val = c.float64(e.Val)
	case *ast.VariableExpr:
		val, err = c.compileVariableExpr(e)
	case *ast.BinaryExpr:
		val, err = c.compileBinaryExpr(e)
	default:
		err = fmt.Errorf("error compiling. node type not handled: %s", node)
	}
	return
}

func (c *compiler) compileVariableExpr(e *ast.VariableExpr) (val llvm.Value, err error) {
	// Look this variable up in the function
	var ok bool
	if val, ok = c.namedValues[e.Name]; !ok {
		err = fmt.Errorf("unknown variable name: %s", e.Name)
	}
	return
}

func (c *compiler) compileBinaryExpr(e *ast.BinaryExpr) (val llvm.Value, err error) {
	var lhs, rhs llvm.Value
	if lhs, err = c.Compile(e.Left); err != nil {
		return
	}
	if rhs, err = c.Compile(e.Right); err != nil {
		return
	}

	// TODO: better switch
	switch e.Op.Raw()[0] {
	case '+':
		val = c.builder.CreateFAdd(lhs, rhs, "addtmp")
	case '-':
		val = c.builder.CreateFSub(lhs, rhs, "subtmp")
	case '*':
		val = c.builder.CreateFMul(lhs, rhs, "multmp")
	case '<':
		boolVal := c.builder.CreateFCmp(llvm.FloatULT, lhs, rhs, "cmptmp")
		// Convert bool 0/1 to double 0.0 or 1.0
		val = c.builder.CreateUIToFP(boolVal, c.DoubleType(), "booltmp")
	default:
		err = fmt.Errorf("invalid binary operator: %s", t.Op)
	}
	return
}

func (c *compiler) float64(num float64) llvm.Value {
	t := c.DoubleType()
	return llvm.ConstFloat(t, num)
}
