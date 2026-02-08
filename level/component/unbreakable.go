package component

import "io"

var _ DataComponent = (*Unbreakable)(nil)

type Unbreakable struct{}

// ID implements DataComponent.
func (Unbreakable) ID() string {
	return "minecraft:unbreakable"
}

// ReadFrom implements DataComponent.
func (u *Unbreakable) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

// WriteTo implements DataComponent.
func (u *Unbreakable) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}
