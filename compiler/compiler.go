package compiler

import (
	"fmt"

	"github.com/evovetech/lex/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

type Compiler interface {
	GetContext() llvm.Context
	GetModule() llvm.Module
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

func (c *compiler) GetModule() llvm.Module {
	return c.module
}

func (c *compiler) Compile(node ast.Node) (val llvm.Value, err error) {
	//fmt.Printf("compiling %T -> %s >>>\n", node, node)
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
	case *ast.IfExpr:
		val, err = c.compileIfExpression(e)
	default:
		err = fmt.Errorf("error compiling. node type not handled: %s", node)
	}
	//fmt.Printf("<<< compiled %T -> (val=%v, err=%v)\n", node, val, err)
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
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				fmt.Printf("error: %s\n", e.Error())
				err = e
			default:
				err = fmt.Errorf("%v", e)
				fmt.Printf("error: %s\n", err.Error())
			}
		}
	}()

	// Look up the name in the global module table.
	var fn llvm.Value
	if fn = c.module.NamedFunction(e.Callee); fn.IsNil() {
		err = fmt.Errorf("unknown function referenced: %s", e)
		return
	}

	var size int
	if size = fn.ParamsCount(); size != len(e.Args) {
		err = fmt.Errorf("incorrect # of arguments passed: %d (expected %d)",
			len(e.Args), size)
		return
	}

	var args []llvm.Value
	for _, arg := range e.Args {
		var argVal llvm.Value
		if argVal, err = c.Compile(arg); err == nil {
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
	if fn = c.module.NamedFunction(e.Proto.Name); fn.IsNil() {
		if fn, err = c.Compile(e.Proto); err != nil {
			return
		}
		if fn.IsNil() {
			err = fmt.Errorf("function is nil")
			return
		}
	}
	if fn.BasicBlocksCount() > 0 {
		if fn.Name() == "__anon_expr" {
			fn.EraseFromParentAsFunction()
			return c.compileFunction(e)
		}
		err = fmt.Errorf("function cannot be redefined")
		return
	}

	// Create a new basic block to start insertion into.
	block := c.AddBasicBlock(fn, "entry")
	c.builder.SetInsertPointAtEnd(block)

	// Record the function arguments in the NamedValues map.
	m := c.namedValues
	for k := range m {
		delete(m, k)
	}
	for _, arg := range fn.Params() {
		m[arg.Name()] = arg
	}

	var body llvm.Value
	if body, err = c.Compile(e.Body); err != nil || body.IsNil() {
		if err == nil {
			err = fmt.Errorf("body is nil: %s", e.Body)
		}
		fn.EraseFromParentAsFunction()
		return
	}

	// finish off the function
	c.builder.CreateRet(body)

	// validate the generated code, checking for consistency.
	llvm.VerifyFunction(fn, llvm.PrintMessageAction)

	// optimize the function
	//c.fpm.RunFunc(fn)
	return
}

func (c *compiler) compileIfExpression(e *ast.IfExpr) (ret llvm.Value, err error) {
	var ifVal llvm.Value
	if ifVal, err = c.Compile(e.Cond); err != nil {
		return
	}

	// Convert condition to a bool by comparing non-equal to 0.0.
	ifVal = c.builder.CreateFCmp(llvm.FloatONE, ifVal, c.float64(0), "ifcond")

	// begin
	startBB := c.builder.GetInsertBlock()
	function := startBB.Parent()

	// Emit 'then' value
	thenBB := c.AddBasicBlock(function, "then")
	c.builder.SetInsertPointAtEnd(thenBB)
	var thenVal llvm.Value
	if thenVal, err = c.Compile(e.Then); err != nil {
		return
	}
	/*
	Codegen of 'then' can change the current block, update then_bb for the
       * phi. We create a new name because one is used for the phi node, and the
       * other is used for the conditional branch.
	 */
	newThenBB := c.builder.GetInsertBlock()

	// Emit else value
	elseBB := c.AddBasicBlock(function, "else")
	c.builder.SetInsertPointAtEnd(elseBB)
	var elseVal llvm.Value
	if elseVal, err = c.Compile(e.Else); err != nil {
		return
	}
	newElseBB := c.builder.GetInsertBlock()

	// Emit merge block
	mergeBB := c.AddBasicBlock(function, "ifcont")
	c.builder.SetInsertPointAtEnd(mergeBB)
	phi := c.builder.CreatePHI(c.DoubleType(), "iftmp")
	phi.AddIncoming([]llvm.Value{
		thenVal,
		elseVal,
	}, []llvm.BasicBlock{
		thenBB,
		elseBB,
	})

	// return to start block to add conditional branch
	c.builder.SetInsertPointAtEnd(startBB)
	c.builder.CreateCondBr(ifVal, thenBB, elseBB)

	// set an unconditionaal branch at the end of
	// the then block and the else block
	// to the merge block
	c.builder.SetInsertPointAtEnd(newThenBB)
	c.builder.CreateBr(mergeBB)
	c.builder.SetInsertPointAtEnd(newElseBB)
	c.builder.CreateBr(mergeBB)

	// finally set the builder to the end of the merge block
	c.builder.SetInsertPointAtEnd(mergeBB)
	ret = phi
	return
}

func (c *compiler) float64(num float64) llvm.Value {
	t := c.DoubleType()
	return llvm.ConstFloat(t, num)
}
