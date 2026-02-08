package component

import (
	"io"
)

var _ DataComponent = (*Repairable)(nil)

type Repairable struct {
	Items IDSet
}

// ID implements DataComponent.
func (Repairable) ID() string {
	return "minecraft:repairable"
}

// ReadFrom implements DataComponent.
func (rep *Repairable) ReadFrom(r io.Reader) (n int64, err error) {
	return rep.Items.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (rep *Repairable) WriteTo(w io.Writer) (n int64, err error) {
	return rep.Items.WriteTo(w)
}
