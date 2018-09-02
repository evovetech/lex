package token

type Token uint8

const (
	noneToken Token = iota

	ILLEGAL
	EOF

	// TODO:
	UNKNOWN

	// Identifier / Literal
	IDENTIFIER
	NUMBER

	// Operators
	ASSIGN
	PLUS
	MINUS
	STAR

	LT
	GT

	COMMA
	SEMICOLON

	LPAREN
	RPAREN

	// others
	COMMENT

	// keywords
	DEF
	EXTERN
	IF
	THEN
	ELSE
	FOR
	IN

	maxToken
)

func (tok Token) IsValid() bool {
	return noneToken < tok && tok < maxToken
}

func (tok Token) String() string {
	return tokenNames[tok]
}

func (tok Token) IsDone() bool {
	return tok <= EOF || tok >= maxToken
}

var tokenNames = [...]string{
	ILLEGAL: "illegal token",
	EOF:     "end of file",

	UNKNOWN: "unknown",

	IDENTIFIER: "identifier",
	NUMBER:     "number",

	ASSIGN: "assign",
	PLUS:   "plus",
	MINUS:  "minus",
	STAR:   "star",

	LT: "less than",
	GT: "greater than",

	COMMA:     "comma",
	SEMICOLON: "semi-colon",

	LPAREN: "left parenthesis",
	RPAREN: "right parenthesis",

	COMMENT: "comment",

	DEF:    "def",
	EXTERN: "extern",
	IF:     "if",
	THEN:   "then",
	ELSE:   "else",
	FOR:    "for",
	IN:     "in",
}

// keywordToken records the special tokens for
// strings that should not be treated as ordinary identifiers.
var keywordToken = map[string]Token{
	"def":    DEF,
	"extern": EXTERN,
	"if":     IF,
	"then":   THEN,
	"else":   ELSE,
	"for":    FOR,
	"in":     IN,
}

type Tok interface {
	Kind() Token
	Val() Value
}

type Reader interface {
	Peek(buf []Char) error
	Consume(n int) error
}

type Char interface {
	Val() rune
	Size() int
}

type Tokenizer interface {
	Start(ch Char) bool
	Continue(i int, ch Char) bool
	End() (Token, error)
}


/*

Lexer:    in(chars  []rune)       => out(token  Token)
Parser:   in(tokens []Token)      => out(expr   Expression)
Eval:     in(exprs  []Expression) => out(result Result)
Compiler: in(exprs  []Expression) => out(prog   Program)

*/
