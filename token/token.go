package token

type Token uint8

const (
	noneToken Token = iota

	ILLEGAL
	EOF

	// TODO:
	UNKNOWN

	// commands
	DEF
	EXTERN

	// primary
	IDENT
	NUMBER

	// others
	COMMENT

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
	DEF:     "def",
	EXTERN:  "extern",
	IDENT:   "identifier",
	NUMBER:  "number",
	COMMENT: "comment",
}

// keywordToken records the special tokens for
// strings that should not be treated as ordinary identifiers.
var keywordToken = map[string]Token{
	"def":    DEF,
	"extern": EXTERN,
}
