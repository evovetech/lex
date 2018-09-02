package lex

import (
	"fmt"
	"io"

	"github.com/evovetech/lex/compiler"
	"github.com/evovetech/lex/token"
	"llvm.org/llvm/bindings/go/llvm"
)

type Repl interface {
	Loop()
}

type repl struct {
	*Parser
	compiler    compiler.Compiler
	ee          llvm.ExecutionEngine
	out, errOut io.Writer
}

func NewRepl(c compiler.Compiler, in io.Reader, out, errOut io.Writer) Repl {
	lex := NewLexer(in)
	parser := NewParser(lex)
	ee, err := llvm.NewExecutionEngine(c.GetModule())
	if err != nil {
		panic(err)
	}
	return &repl{
		Parser:   parser,
		compiler: c,
		ee:       ee,
		out:      out,
		errOut:   errOut,
	}
}

func (r *repl) Loop() {
	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Printf("recovered in Loop() %v\n", r)
	//	}
	//}()

	for {
		fmt.Fprintf(r.out, "ready> ")
		switch tok := r.CurToken(); tok.Kind() {
		case token.EOF:
			return
		case token.SEMICOLON,
			token.COMMENT:
			// ignore top-level semicolons. & comments
			r.NextToken()
		case token.DEF:
			r.handleDefinition()
		case token.EXTERN:
			r.handleExtern()
		default:
			r.handleTopLevelExpression()
		}
	}
}

func (r *repl) handleDefinition() {
	if def, err := r.ParseDefinition(); err == nil {
		fmt.Fprintf(r.out, "parsed a def: %s\n", def)
		if code, err := r.compiler.Compile(def); err == nil {
			fmt.Fprintf(r.out, "Read function definition:")
			code.Dump()
			fmt.Fprint(r.out, "\n")
		} else {
			fmt.Fprintf(r.errOut, "error compiling def: %s\n", err.Error())
		}
	} else {
		fmt.Fprintf(r.errOut, "error handling def: %s\n", err.Error())
		// skip next token for error recovery
		r.NextToken()
	}
}

func (r *repl) handleExtern() {
	if extern, err := r.ParseExtern(); err == nil {
		fmt.Fprintf(r.out, "parsed an extern: %s\n", extern)
		if code, err := r.compiler.Compile(extern); err == nil {
			fmt.Fprintf(r.out, "Read extern:")
			code.Dump()
			fmt.Fprint(r.out, "\n")
		} else {
			fmt.Fprintf(r.errOut, "error compiling extern: %s\n", err.Error())
		}
	} else {
		fmt.Fprintf(r.errOut, "error handling extern: %s\n", err.Error())
		// skip next token for error recovery
		r.NextToken()
	}
}

func (r *repl) handleTopLevelExpression() {
	if topLevel, err := r.ParseTopLevelExpression(); err == nil {
		fmt.Fprintf(r.out, "parsed a top-level: %s\n", topLevel)
		if code, err := r.compiler.Compile(topLevel); err == nil {
			fmt.Fprintf(r.out, "Read top-level epression:")
			code.Dump()
			fmt.Fprint(r.out, "\n")

			result := r.ee.RunFunction(code, []llvm.GenericValue{})
			rfloat := result.Float(r.compiler.GetContext().DoubleType())
			fmt.Fprintf(r.out, "Evaluated to a %v\n", rfloat)
		} else {
			fmt.Fprintf(r.errOut, "error compiling top-level expression: %s\n", err.Error())
		}
	} else {
		fmt.Fprintf(r.errOut, "error handling top level: %s\n", err.Error())
		// skip next token for error recovery
		r.NextToken()
	}
}
