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

func NewCompiler(name string) Compiler {
	ctx := llvm.GlobalContext()
	mod := ctx.NewModule(name)
	return &compiler{
		Context:     ctx,
		builder:     ctx.NewBuilder(),
		module:      mod,
		fpm:         newFunctionPassManager(mod),
		namedValues: make(map[string]llvm.Value),
	}
}

func newFunctionPassManager(mod llvm.Module) llvm.PassManager {
	// Create a new pass manager attached to module.
	fpm := llvm.NewFunctionPassManagerForModule(mod)

	// Do simple "peephole" optimizations and bit-twiddling optzns.
	fpm.AddInstructionCombiningPass()

	// Reassociate expressions.
	fpm.AddReassociatePass()

	// Eliminate Common SubExpressions.
	fpm.AddGVNPass()

	// Simplify the control flow graph (deleting unreachable blocks, etc).
	fpm.AddCFGSimplificationPass()

	// init
	fpm.InitializeFunc()

	return fpm
}

type compiler struct {
	llvm.Context
	builder     llvm.Builder
	module      llvm.Module
	fpm         llvm.PassManager
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
	case *ast.CallExpr:
		val, err = c.compileCallExpr(e)
	case *ast.PrototypeExpr:
		val, err = c.compilePrototype(e)
	case *ast.FunctionExpr:
		val, err = c.compileFunction(e)
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
		err = fmt.Errorf("invalid binary operator: %s", e.Op)
	}
	return
}

func (c *compiler) compileCallExpr(e *ast.CallExpr) (val llvm.Value, err error) {
	// Look up the name in the global module table.
	var fn llvm.Value
	if fn = c.module.NamedFunction(e.Callee); fn.IsNil() {
		err = fmt.Errorf("unknown function referenced: %s", e)
		return
	}

	var size int
	if size = fn.ParamsCount(); size != len(e.Args) {
		err = fmt.Errorf("incorrect # of arguments passed: %d (expected %d)",
			size, len(e.Args))
		return
	}

	var args []llvm.Value
	for _, arg := range e.Args {
		var argVal llvm.Value
		if argVal, err = c.Compile(arg); err != nil {
			args = append(args, argVal)
			continue
		}
		return
	}

	val = c.builder.CreateCall(fn, args, "calltmp")
	return
}

func (c *compiler) compilePrototype(e *ast.PrototypeExpr) (fn llvm.Value, err error) {
	var size = len(e.Args)
	var types = make([]llvm.Type, size)
	for i := 0; i < size; i++ {
		types[i] = c.DoubleType()
	}

	fnType := llvm.FunctionType(c.DoubleType(), types, false)
	fn = llvm.AddFunction(c.module, e.Name, fnType)
	fn.SetLinkage(llvm.ExternalLinkage)

	// Set names for all arguments.
	for i, arg := range fn.Params() {
		arg.SetName(e.Args[i])
	}
	return
}

func (c *compiler) compileFunction(e *ast.FunctionExpr) (fn llvm.Value, err error) {
	fn = c.module.NamedFunction(e.Proto.Name)
	if fn.IsNil() {
		if fn, err = c.Compile(e.Proto); err != nil {
			return
		}
		if fn.IsNil() {
			err = fmt.Errorf("function is nil")
			return
		}
	}
	if fn.BasicBlocksCount() > 0 {
		err = fmt.Errorf("function cannot be redefined")
		return
	}

	// Create a new basic block to start insertion into.
	block := c.AddBasicBlock(fn, "entry")
	c.builder.SetInsertPoint(block, llvm.Value{})

	// Record the function arguments in the NamedValues map.
	m := c.namedValues
	for k := range m {
		delete(m, k)
	}
	for _, arg := range fn.Params() {
		m[arg.Name()] = arg
	}

	var body llvm.Value
	if body, err = c.Compile(e.Body); err == nil {
		if !body.IsNil() {
			// finish off the function
			c.builder.CreateRet(body)

			// validate the generated code, checking for consistency.
			llvm.VerifyFunction(fn, llvm.PrintMessageAction)

			// optimize the function
			c.fpm.RunFunc(fn)
			return
		}
		err = fmt.Errorf("body is nil: %s", e.Body)
	}

	fn.EraseFromParentAsFunction()
	return
}

func (c *compiler) float64(num float64) llvm.Value {
	t := c.DoubleType()
	return llvm.ConstFloat(t, num)
}
