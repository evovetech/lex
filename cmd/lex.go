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
var vers = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run:   version,
}
var comp = &cobra.Command{
	Use:   "compile",
	Short: "Compiles",
	RunE:  compile,
}
var optimize bool

func main() {
	// execute
	if err := lex.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	comp.Flags().BoolVar(&optimize, "optimize", false, "optimize llvm output -- runs function pass manager")
	lex.AddCommand(
		vers,
		comp,
	)
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
	c.GetCompiler().GetOptions().Optimize = optimize
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
	var name string
	if optimize {
		name = "output"
	} else {
		name = "output-unoptimized"
	}
	file := fmt.Sprintf("sample/%s.bc", name)
	e = c.Write(file, fType)
	return
}
