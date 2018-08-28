package main

import (
	"fmt"
	"os"

	. "github.com/evovetech/lex"
	"github.com/spf13/cobra"
	"llvm.org/llvm/bindings/go/llvm"
)

var lex = &cobra.Command{
	Use:   "lex",
	Short: "lex input",
	Run:   run,
}

func main() {
	// execute
	if err := lex.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	lex.AddCommand(
		&cobra.Command{
			Use:   "version",
			Short: "Print the version",
			Run: func(cmd *cobra.Command, args []string) {
				version := llvm.Version
				fmt.Printf("llvm %s", version)
				fmt.Println()
			},
		},
	)
}

func run(cmd *cobra.Command, args []string) {
	in, out, err := os.Stdin, os.Stdout, os.Stderr
	eval := NewEval(in, out, err)
	eval.Loop()
}
