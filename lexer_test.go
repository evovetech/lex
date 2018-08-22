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

	tok := lex.NextToken()
	for !tok.IsDone() {
		fmt.Printf("  %s\n", tok)
		tok = lex.NextToken()
	}

	fmt.Println("<<< end")
}
