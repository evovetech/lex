package token

import "fmt"

type Value struct {
	raw      []rune
	beg, end Position
	err      error
}

func (v Value) Raw() []rune {
	return v.raw
}

func (v Value) RawString() string {
	return string(v.raw)
}

func (v Value) Pos() (beg, end Position) {
	beg, end = v.beg, v.end
	return
}

func (v Value) Error() error {
	return v.err
}

func (v Value) String() string {
	if err := v.err; err != nil {
		return fmt.Sprintf("{err=%s}", err.Error())
	}
	return fmt.Sprintf("{raw=%s, pos=(%s -- %s)}", v.RawString(), v.beg, v.end)
}
