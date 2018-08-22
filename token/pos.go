package token

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// A Position describes the location of a rune of input.
type Position struct {
	file *string // filename (indirect for compactness)
	Line int32   // 1-based line number
	Col  int32   // 1-based column number (strictly: rune)
}

// IsValid reports whether the position is valid.
func (p Position) IsValid() bool {
	return p.Line >= 1
}

// Filename returns the name of the file containing this position.
func (p Position) Filename() string {
	if p.file != nil {
		return *p.file
	}
	return "<unknown>"
}

// MakePosition returns position with the specified components.
func MakePosition(file *string, line, col int32) Position { return Position{file, line, col} }

// add returns the position at the end of s, assuming it starts at p.
func (p Position) add(s string) Position {
	if n := strings.Count(s, "\n"); n > 0 {
		p.Line += int32(n)
		s = s[strings.LastIndex(s, "\n")+1:]
		p.Col = 1
	}
	p.Col += int32(utf8.RuneCountInString(s))
	return p
}

func (p Position) String() string {
	if p.Col > 0 {
		return fmt.Sprintf("%s:%d:%d", p.Filename(), p.Line, p.Col)
	}
	return fmt.Sprintf("%s:%d", p.Filename(), p.Line)
}

func (p Position) isBefore(q Position) bool {
	if p.Line != q.Line {
		return p.Line < q.Line
	}
	return p.Col < q.Col
}
