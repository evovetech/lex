package lex

import (
	"fmt"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	code := `def fib(x)`
	rd := strings.NewReader(code)
	lex := NewLexer(rd)

	fmt.Println("begin >>>")

	for tok := lex.NextToken(); !tok.IsDone(); tok = lex.NextToken() {
		fmt.Printf("  %s\n", tok)
	}

	fmt.Println("<<< end")
}
