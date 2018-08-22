package lex

import (
	"fmt"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	code := `# Compute the x'th fibonacci number.
def fib(x)
  if x < 3 then
    1
  else
    fib(x-1)+fib(x-2)

# This expression will compute the 40th number.
fib(40)
`
	rd := strings.NewReader(code)
	lex := NewLexer(rd)

	fmt.Println("begin >>>")

	i, max := 0, len(code)
	for tok := lex.NextToken(); !tok.IsDone(); tok = lex.NextToken() {
		fmt.Printf("  %s\n", tok)
		if i++; i >= max {
			break
		}
	}

	fmt.Println("<<< end")
}
