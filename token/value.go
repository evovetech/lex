package token

type Value struct {
	raw      []rune
	beg, end Position
	err      error
}

func (v *Value) Raw() []rune {
	return v.raw
}

func (v *Value) Pos() (beg, end Position) {
	beg, end = v.beg, v.end
	return
}

func (v *Value) Error() error {
	return v.err
}

func (v Value) String() string {
	if err := v.err; err != nil {
		return err.Error()
	}
	return string(v.raw)
}
