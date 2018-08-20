package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"llvm.org/llvm/bindings/go/llvm"
)

var lex = &cobra.Command{
	Use:   "lex",
	Short: "lex input",
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
