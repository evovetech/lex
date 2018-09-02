package main

import (
	"fmt"
	"os"

	. "github.com/evovetech/lex"
	"github.com/evovetech/lex/compiler"
	"github.com/spf13/cobra"
	"llvm.org/llvm/bindings/go/llvm"
)

var lex = &cobra.Command{
	Use:   "lex",
	Short: "lex input",
	Run:   repl,
}

func main() {
	// execute
	if err := lex.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	lex.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run:   version,
	})
	lex.AddCommand(&cobra.Command{
		Use:   "compile",
		Short: "Compiles",
		RunE:  compile,
	})
}

func version(_ *cobra.Command, _ []string) {
	version := llvm.Version
	fmt.Printf("llvm %s", version)
	fmt.Println()
}

func repl(_ *cobra.Command, _ []string) {
	in, out, err := os.Stdin, os.Stdout, os.Stderr
	c := compiler.NewCompiler("lex")
	repl := NewRepl(c, in, out, err)
	repl.Loop()
}

func compile(_ *cobra.Command, _ []string) (e error) {
	var m *compiler.Machine
	if m, e = compiler.NewDefaultMachine(); e != nil {
		return
	}


	c := m.NewCompiler("lex")
	in, out, err := os.Stdin, os.Stdout, os.Stderr
	repl := NewRepl(c.GetCompiler(), in, out, err)
	repl.Loop()

	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()

	// TODO:
	fType := llvm.AssemblyFile
	e = c.Write("sample/output.bc", fType)
	return
}
