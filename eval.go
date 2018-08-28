package lex

import (
	"fmt"
	"io"

	"github.com/evovetech/lex/token"
)

type Eval interface {
	Loop()
}

type eval struct {
	*Parser
	out, err io.Writer
}

func NewEval(in io.Reader, out, err io.Writer) Eval {
	lex := NewLexer(in)
	parser := NewParser(lex)
	return &eval{
		Parser: parser,
		out:    out,
		err:    err,
	}
}

func (e *eval) Loop() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered in Loop()", r)
		}
	}()
	for {
		fmt.Fprintf(e.out, "ready> ")
		switch tok := e.CurToken(); tok.Kind() {
		case token.EOF:
			return
		case token.SEMICOLON:
			// ignore top-level semicolons.
			e.NextToken()
		case token.DEF:
			e.handleDefinition()
		case token.EXTERN:
			e.handleExtern()
		default:
			e.handleTopLevelExpression()
		}
	}
}

func (e *eval) handleDefinition() {
	if def, err := e.ParseDefinition(); err == nil {
		fmt.Fprintf(e.out, "parsed a def: %s\n", def)
	} else {
		fmt.Fprintf(e.err, "error handling def: %s", err.Error())
		// skip next token for error recovery
		e.NextToken()
	}
}

func (e *eval) handleExtern() {
	if extern, err := e.ParseExtern(); err == nil {
		fmt.Fprintf(e.out, "parsed an extern: %s\n", extern)
	} else {
		fmt.Fprintf(e.err, "error handling extern: %s", err.Error())
		// skip next token for error recovery
		e.NextToken()
	}
}

func (e *eval) handleTopLevelExpression() {
	if topLevel, err := e.ParseTopLevelExpression(); err == nil {
		fmt.Fprintf(e.out, "parsed a top level: %s\n", topLevel)
	} else {
		fmt.Fprintf(e.err, "error handling top level: %s", err.Error())
		// skip next token for error recovery
		e.NextToken()
	}
}
