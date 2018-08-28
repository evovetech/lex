package lex

import (
	"fmt"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	code := `# Compute the x'th fibonacci number.
def fib(x)
  if x < 3 then
    1
  else
    fib(x-1)+fib(x-2)

# This expression will compute the 40th number.
fib(40)

extern sin(arg);
extern cos(arg);
extern atan2(arg1 arg2);

atan2(sin(.4), cos(42))
`
	rd := strings.NewReader(code)
	lex := NewLexer(rd)
	p := NewParser(lex)

	fmt.Println("begin >>>")

	i, max := 0, len(code)
	for {
		if expr, err := p.ParsePrimary(); err != nil {
			fmt.Printf("  %s\n", err.Error())
		} else {
			fmt.Printf("  %s\n", expr)
		}
		if i++; i >= max {
			break
		}
		if p.CurToken().IsDone() {
			break
		}
	}

	fmt.Println("<<< end")
}
