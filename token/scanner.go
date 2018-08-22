package token

import (
	"bufio"
	"io"
	"unicode"
)

type Scanner interface {
	NextToken(val *Value) Token
}

func NewScanner(rd io.Reader) Scanner {
	name := "src"
	sc := &scanner{
		rd:  bufio.NewReader(rd),
		pk:  new(ch),
		pos: MakePosition(&name, 1, 0),
	}
	// init
	sc.fill()
	return sc
}

type scanner struct {
	rd  *bufio.Reader
	pk  *ch
	pos Position
}

func (sc *scanner) NextToken(val *Value) Token {
	var ch = sc.peek()
	sc.eatSpace()

	var pos = sc.pos

	// start/end token
	sc.startToken(val)
	defer sc.endToken(val, &pos)

	if err, tok := ch.error(); err {
		val.err = ch.err
		return tok
	}

_:
	ch = sc.peek()
	if err, tok := ch.error(); err {
		val.err = ch.err
		return tok
	}

	r := ch.val
	switch {
	case isIdentStart(r):
		return sc.scanIdent(val, &pos)
	}

	// TODO:
	val.raw = append(val.raw, r)
	sc.read()
	return ILLEGAL
}

func (sc *scanner) scanIdent(val *Value, pos *Position) Token {
	for {
		ch := sc.peek()
		if err, _ := ch.error(); err {
			break
		}

		if r := ch.val; isIdent(r) {
			*pos = sc.pos
			val.raw = append(val.raw, r)
			sc.read()
			continue
		}
		break
	}

	if k, ok := keywordToken[val.RawString()]; ok {
		return k
	}

	return IDENT
}

func (sc *scanner) eatSpace() {
	for ch := sc.peek(); unicode.IsSpace(ch.val); ch = sc.peek() {
		sc.read()
	}
}

func (sc *scanner) startToken(val *Value) {
	val.raw = val.raw[0:0]
	val.beg, val.end = sc.pos, sc.pos
	val.err = nil
}

func (sc *scanner) endToken(val *Value, pos *Position) {
	val.end = *pos
}

func (sc *scanner) peek() ch {
	return *sc.pk
}

func (sc *scanner) read() ch {
	next := *sc.pk
	sc.fill()
	if next.isNewline() {
		if next.val == '\r' && sc.pk.val == '\n' {
			// combine into '\n'
			next = *sc.pk
			sc.fill()
		}
		sc.pos.Line++
		sc.pos.Col = 1
	}
	return next
}

func (sc *scanner) fill() {
	sc.pk.fill(sc.rd)
	sc.pos.Col++
}

type ch struct {
	val  rune
	size int
	err  error
}

func (c *ch) error() (bool, Token) {
	if c.err == nil {
		return false, noneToken
	}
	switch c.err {
	case io.EOF:
		return true, EOF
	default:
		return true, ERROR
	}
}

func (c *ch) fill(rd *bufio.Reader) {
	c.val, c.size, c.err = rd.ReadRune()
}

func (c *ch) isNewline() bool {
	return isNewline(c.val)
}

func isNewline(r rune) bool {
	// This property isn't the same as Z; special-case it.
	if uint32(r) <= unicode.MaxLatin1 {
		switch r {
		case '\r', '\n', '\u0085':
			return true
		default:
			return false
		}
	}
	return unicode.Is(unicode.Zl, r)
}

// isIdent reports whether c is an identifier rune.
func isIdent(c rune) bool {
	return isdigit(c) || isIdentStart(c)
}

func isIdentStart(c rune) bool {
	return 'a' <= c && c <= 'z' ||
		'A' <= c && c <= 'Z' ||
		c == '_' ||
		unicode.IsLetter(c)
}

func isdigit(c rune) bool  { return '0' <= c && c <= '9' }
func isodigit(c rune) bool { return '0' <= c && c <= '7' }
func isxdigit(c rune) bool { return isdigit(c) || 'A' <= c && c <= 'F' || 'a' <= c && c <= 'f' }
func isbdigit(c rune) bool { return '0' == c || c == '1' }
