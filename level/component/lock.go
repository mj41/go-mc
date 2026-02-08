package component

import (
	"io"

	"github.com/Tnze/go-mc/nbt/dynbt"
	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*Lock)(nil)

type Lock struct {
	Data dynbt.Value
}

// ID implements DataComponent.
func (Lock) ID() string {
	return "minecraft:lock"
}

// ReadFrom implements DataComponent.
func (l *Lock) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.NBTField{V: &l.Data, AllowUnknownFields: true}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (l *Lock) WriteTo(w io.Writer) (n int64, err error) {
	return pk.NBTField{V: &l.Data, AllowUnknownFields: true}.WriteTo(w)
}
