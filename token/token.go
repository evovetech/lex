package token

type Token uint8

const (
	noneToken Token = iota

	ILLEGAL
	EOF

	// TODO:
	UNKNOWN

	// Identifier / Literal
	IDENT
	NUMBER

	// Operators
	ASSIGN
	PLUS
	MINUS

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

	IDENT:  "identifier",
	NUMBER: "number",

	ASSIGN: "assign",
	PLUS:   "plus",
	MINUS:  "minus",

	LT: "less than",
	GT: "greater than",

	COMMA:     "comma",
	SEMICOLON: "semi-colon",

	LPAREN: "left parenthesis",
	RPAREN: "right parenthesis",

	COMMENT: "comment",

	DEF:    "def",
	EXTERN: "extern",
}

// keywordToken records the special tokens for
// strings that should not be treated as ordinary identifiers.
var keywordToken = map[string]Token{
	"def":    DEF,
	"extern": EXTERN,
}
