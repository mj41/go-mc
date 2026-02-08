package component

import (
	"io"
)

var _ DataComponent = (*UseRemainder)(nil)

type UseRemainder struct {
	Remainder SlotData
}

// ID implements DataComponent.
func (UseRemainder) ID() string {
	return "minecraft:use_remainder"
}

// ReadFrom implements DataComponent.
func (u *UseRemainder) ReadFrom(r io.Reader) (n int64, err error) {
	return u.Remainder.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (u *UseRemainder) WriteTo(w io.Writer) (n int64, err error) {
	return u.Remainder.WriteTo(w)
}
