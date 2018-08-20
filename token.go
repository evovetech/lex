package lex

type Token uint8

const (
	ILLEGAL Token = iota
	EOF

	NEWLINE
	SPACE

	LETTER
	MARK
	NUMBER

	CONTROL
	PUNCT
	SYMBOL

	maxToken
)
