package lex

import (
	"fmt"
	"io"

	"github.com/evovetech/lex/token"
)

type Repl interface {
	Loop()
}

type repl struct {
	*Parser
	out, err io.Writer
}

func NewRepl(in io.Reader, out, err io.Writer) Repl {
	lex := NewLexer(in)
	parser := NewParser(lex)
	return &repl{
		Parser: parser,
		out:    out,
		err:    err,
	}
}

func (r *repl) Loop() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered in Loop()", r)
		}
	}()
	for {
		fmt.Fprintf(r.out, "ready> ")
		switch tok := r.CurToken(); tok.Kind() {
		case token.EOF:
			return
		case token.SEMICOLON:
			// ignore top-level semicolons.
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
	} else {
		fmt.Fprintf(r.err, "error handling def: %s\n", err.Error())
		// skip next token for error recovery
		r.NextToken()
	}
}

func (r *repl) handleExtern() {
	if extern, err := r.ParseExtern(); err == nil {
		fmt.Fprintf(r.out, "parsed an extern: %s\n", extern)
	} else {
		fmt.Fprintf(r.err, "error handling extern: %s\n", err.Error())
		// skip next token for error recovery
		r.NextToken()
	}
}

func (r *repl) handleTopLevelExpression() {
	if topLevel, err := r.ParseTopLevelExpression(); err == nil {
		fmt.Fprintf(r.out, "parsed a top level: %s\n", topLevel)
	} else {
		fmt.Fprintf(r.err, "error handling top level: %s\n", err.Error())
		// skip next token for error recovery
		r.NextToken()
	}
}
